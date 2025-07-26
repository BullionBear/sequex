package binance

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var BinanceWssBase = "wss://stream.binance.com:9443"

const (
	pingInterval   = 20 * time.Second
	reconnectDelay = 5 * time.Second
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

func (w *BinanceWSConn) readLoop() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}

		_, message, err := w.conn.ReadMessage()
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
