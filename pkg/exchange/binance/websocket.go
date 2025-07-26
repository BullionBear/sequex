package binance

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSConnection interface for all WebSocket connection types
type WSConnection interface {
	Connect() error
	Disconnect()
	SetOnMessage(func([]byte))
}

const (
	pingInterval      = 20 * time.Second
	reconnectDelay    = 5 * time.Second
	keepaliveInterval = 30 * time.Minute // Keepalive interval for user data streams
)

type BinanceWSConn struct {
	conn      *websocket.Conn
	url       string
	mu        sync.Mutex
	connected bool
	ctx       context.Context
	cancel    context.CancelFunc
	reconnect bool
	OnMessage func([]byte) // Callback for handling messages
}

func NewBinanceWSConn(baseURL, streamPath string) *BinanceWSConn {
	ctx, cancel := context.WithCancel(context.Background())
	return &BinanceWSConn{
		url:       baseURL + streamPath,
		ctx:       ctx,
		cancel:    cancel,
		reconnect: true,
	}
}

func (w *BinanceWSConn) Connect() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	dialer := websocket.DefaultDialer
	c, _, err := dialer.Dial(w.url, nil)
	if err != nil {
		return err
	}
	w.conn = c
	w.connected = true
	go w.readLoop()
	go w.pingLoop()
	return nil
}

func (w *BinanceWSConn) SetOnMessage(handler func([]byte)) {
	w.OnMessage = handler
}

func (w *BinanceWSConn) readLoop() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}

		w.mu.Lock()
		conn := w.conn
		w.mu.Unlock()

		if conn == nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			// Check if this is a graceful shutdown
			if w.ctx.Err() != nil {
				return
			}
			log.Printf("[WS] Read error: %v", err)
			w.handleDisconnect()
			return
		}

		// Call the message handler if set
		if w.OnMessage != nil {
			w.OnMessage(message)
		}
	}
}

func (w *BinanceWSConn) pingLoop() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.mu.Lock()
			if w.conn != nil {
				if err := w.conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					log.Printf("[WS] Pong error: %v", err)
				}
			}
			w.mu.Unlock()
		}
	}
}

func (w *BinanceWSConn) handleDisconnect() {
	w.mu.Lock()
	w.connected = false
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}
	shouldReconnect := w.reconnect && w.ctx.Err() == nil
	w.mu.Unlock()

	if shouldReconnect {
		log.Printf("[WS] Reconnecting in %v...", reconnectDelay)
		time.Sleep(reconnectDelay)
		if err := w.Connect(); err != nil {
			log.Printf("[WS] Reconnect failed: %v", err)
		}
	}
}

func (w *BinanceWSConn) Disconnect() {
	w.mu.Lock()
	w.reconnect = false
	if w.conn != nil {
		// Set read deadline to unblock ReadMessage immediately
		w.conn.SetReadDeadline(time.Now())
		// Send close frame before closing
		w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		w.conn.Close()
		w.conn = nil
	}
	w.connected = false
	w.mu.Unlock()
	w.cancel()
}

func (w *BinanceWSConn) IsConnected() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.connected
}

// UserDataWSConn handles WebSocket connections for user data streams with listen key management
type UserDataWSConn struct {
	conn               *websocket.Conn
	baseURL            string
	listenKey          string
	client             *Client
	mu                 sync.Mutex
	connected          bool
	ctx                context.Context
	cancel             context.CancelFunc
	reconnect          bool
	OnMessage          func([]byte) // Callback for handling messages
	options            UserDataSubscriptionOptions
	keepaliveTimer     *time.Timer
	reconnectRequested bool
}

func NewUserDataWSConn(baseURL, listenKey string, client *Client, options UserDataSubscriptionOptions) *UserDataWSConn {
	ctx, cancel := context.WithCancel(context.Background())
	return &UserDataWSConn{
		baseURL:   baseURL,
		listenKey: listenKey,
		client:    client,
		ctx:       ctx,
		cancel:    cancel,
		reconnect: true,
		options:   options,
	}
}

func (w *UserDataWSConn) Connect() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Construct URL for user data stream
	url := w.baseURL + "/ws/" + w.listenKey

	dialer := websocket.DefaultDialer
	c, _, err := dialer.Dial(url, nil)
	if err != nil {
		return err
	}
	w.conn = c
	w.connected = true
	w.reconnectRequested = false

	go w.readLoop()
	go w.pingLoop()
	w.startKeepaliveTimer()

	return nil
}

func (w *UserDataWSConn) SetOnMessage(handler func([]byte)) {
	w.OnMessage = handler
}

func (w *UserDataWSConn) Disconnect() {
	w.mu.Lock()
	w.reconnect = false
	if w.keepaliveTimer != nil {
		w.keepaliveTimer.Stop()
	}
	w.cancel()
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}
	w.connected = false
	w.mu.Unlock()

	// Close the listen key
	if w.client != nil {
		ctx := context.Background()
		if _, err := w.client.CloseUserDataStream(ctx, w.listenKey); err != nil {
			log.Printf("[UserDataWS] Failed to close listen key: %v", err)
		}
	}
}

func (w *UserDataWSConn) readLoop() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}

		w.mu.Lock()
		conn := w.conn
		w.mu.Unlock()

		if conn == nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			// Check if this is a graceful shutdown
			if w.ctx.Err() != nil {
				return
			}
			log.Printf("[UserDataWS] Read error: %v", err)
			w.handleDisconnect()
			return
		}

		// Call the message handler if set
		if w.OnMessage != nil {
			w.OnMessage(message)
		}
	}
}

func (w *UserDataWSConn) pingLoop() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.mu.Lock()
			if w.conn != nil {
				if err := w.conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					log.Printf("[UserDataWS] Pong error: %v", err)
				}
			}
			w.mu.Unlock()
		}
	}
}

func (w *UserDataWSConn) startKeepaliveTimer() {
	if w.keepaliveTimer != nil {
		w.keepaliveTimer.Stop()
	}

	w.keepaliveTimer = time.AfterFunc(keepaliveInterval, func() {
		w.sendKeepalive()
	})
}

func (w *UserDataWSConn) sendKeepalive() {
	ctx := context.Background()
	if _, err := w.client.KeepaliveUserDataStream(ctx, w.listenKey); err != nil {
		log.Printf("[UserDataWS] Keepalive failed: %v", err)
		// Don't trigger reconnection for keepalive failures - let listen key expiry handle it
	} else {
		log.Printf("[UserDataWS] Keepalive sent successfully")
		// Schedule next keepalive
		w.startKeepaliveTimer()
	}
}

func (w *UserDataWSConn) handleDisconnect() {
	w.mu.Lock()
	w.connected = false
	if w.keepaliveTimer != nil {
		w.keepaliveTimer.Stop()
	}
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}
	w.mu.Unlock()

	if w.reconnectRequested {
		// This disconnect was requested due to listen key expiry
		w.handleReconnectWithNewListenKey()
	} else if w.reconnect {
		// Normal reconnection with existing listen key
		w.handleReconnect()
	}
}

func (w *UserDataWSConn) handleReconnect() {
	if !w.reconnect {
		return
	}
	log.Printf("[UserDataWS] Attempting to reconnect...")
	time.Sleep(reconnectDelay)
	if err := w.Connect(); err != nil {
		log.Printf("[UserDataWS] Reconnect failed: %v", err)
		go w.handleReconnect()
	} else {
		log.Printf("[UserDataWS] Reconnected successfully")
		if w.options.OnReconnect != nil {
			w.options.OnReconnect()
		}
	}
}

func (w *UserDataWSConn) reconnectWithNewListenKey() {
	w.mu.Lock()
	w.reconnectRequested = true
	if w.conn != nil {
		w.conn.Close()
	}
	w.mu.Unlock()
}

func (w *UserDataWSConn) handleReconnectWithNewListenKey() {
	if !w.reconnect {
		return
	}

	log.Printf("[UserDataWS] Attempting to reconnect with new listen key...")

	// Get new listen key
	ctx := context.Background()
	resp, err := w.client.StartUserDataStream(ctx)
	if err != nil {
		log.Printf("[UserDataWS] Failed to get new listen key: %v", err)
		if w.options.OnError != nil {
			w.options.OnError(err)
		}
		time.Sleep(reconnectDelay)
		go w.handleReconnectWithNewListenKey()
		return
	}

	if resp.Data == nil || resp.Data.ListenKey == "" {
		log.Printf("[UserDataWS] Invalid listen key received")
		time.Sleep(reconnectDelay)
		go w.handleReconnectWithNewListenKey()
		return
	}

	// Update listen key
	w.mu.Lock()
	w.listenKey = resp.Data.ListenKey
	w.mu.Unlock()

	// Try to connect with new listen key
	if err := w.Connect(); err != nil {
		log.Printf("[UserDataWS] Reconnect with new listen key failed: %v", err)
		time.Sleep(reconnectDelay)
		go w.handleReconnectWithNewListenKey()
	} else {
		log.Printf("[UserDataWS] Reconnected successfully with new listen key")
		if w.options.OnReconnect != nil {
			w.options.OnReconnect()
		}
	}
}
