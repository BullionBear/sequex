package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	// Example of REST API usage (existing functionality)
	fmt.Println("=== Binance REST API Example ===")
	restAPIExample()

	fmt.Println("\n=== Binance WebSocket Example ===")
	websocketExample()

	fmt.Println("\n=== Binance User Data Stream Example ===")
	userDataStreamExample()
}

func restAPIExample() {
	// Configure your API credentials here
	cfg := &binance.Config{
		BaseURL:   "https://api.binance.com/api",
		APIKey:    "your-api-key",
		APISecret: "your-api-secret",
	}

	client := binance.NewClient(cfg)

	// Example: Get recent trades (public endpoint)
	// This is a placeholder - implement based on existing client methods
	fmt.Println("REST API client configured successfully")
	_ = client // Use client for actual API calls
}

func websocketExample() {
	// Create WebSocket client (using port 9443 for optimal performance)
	wsClient := binance.NewWSClient(binance.WSConfig{
		BaseURL: binance.MainnetWSBaseUrl9443,
	})

	// Configure Kline subscription options
	klineOptions := binance.KlineSubscriptionOptions{
		OnConnect: func() {
			fmt.Println("🟢 Connected to Binance Kline WebSocket")
		},
		OnReconnect: func() {
			fmt.Println("🔄 Reconnected to Binance Kline WebSocket")
		},
		OnError: func(err error) {
			fmt.Printf("❌ Kline WebSocket Error: %v\n", err)
		},
		OnKline: func(kline binance.WSKline) {
			fmt.Printf("📊 Kline Update: %s %s | Open: %s | High: %s | Low: %s | Close: %s | Volume: %s | Closed: %t\n",
				kline.Symbol, kline.Interval, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, kline.IsClosed)
		},
		OnDisconnect: func() {
			fmt.Println("🔴 Disconnected from Binance Kline WebSocket")
		},
	}

	// Configure Aggregate Trade subscription options
	aggTradeOptions := binance.AggTradeSubscriptionOptions{
		OnConnect: func() {
			fmt.Println("🟢 Connected to Binance AggTrade WebSocket")
		},
		OnReconnect: func() {
			fmt.Println("🔄 Reconnected to Binance AggTrade WebSocket")
		},
		OnError: func(err error) {
			fmt.Printf("❌ AggTrade WebSocket Error: %v\n", err)
		},
		OnAggTrade: func(aggTrade binance.WSAggTrade) {
			buyerType := "Seller"
			if aggTrade.IsBuyerMaker {
				buyerType = "Buyer"
			}
			fmt.Printf("💰 AggTrade: %s | Price: %s | Qty: %s | %s is Maker | AggID: %d\n",
				aggTrade.Symbol, aggTrade.Price, aggTrade.Quantity, buyerType, aggTrade.AggTradeId)
		},
		OnDisconnect: func() {
			fmt.Println("🔴 Disconnected from Binance AggTrade WebSocket")
		},
	}

	// Configure Raw Trade subscription options
	tradeOptions := binance.TradeSubscriptionOptions{
		OnConnect: func() {
			fmt.Println("🟢 Connected to Binance Trade WebSocket")
		},
		OnReconnect: func() {
			fmt.Println("🔄 Reconnected to Binance Trade WebSocket")
		},
		OnError: func(err error) {
			fmt.Printf("❌ Trade WebSocket Error: %v\n", err)
		},
		OnTrade: func(trade binance.WSTrade) {
			buyerType := "Seller"
			if trade.IsBuyerMaker {
				buyerType = "Buyer"
			}
			fmt.Printf("🔥 Trade: %s | Price: %s | Qty: %s | %s is Maker | TradeID: %d\n",
				trade.Symbol, trade.Price, trade.Quantity, buyerType, trade.TradeId)
		},
		OnDisconnect: func() {
			fmt.Println("🔴 Disconnected from Binance Trade WebSocket")
		},
	}

	// Configure Depth subscription options
	depthOptions := binance.DepthSubscriptionOptions{
		OnConnect: func() {
			fmt.Println("🟢 Connected to Binance Depth WebSocket")
		},
		OnReconnect: func() {
			fmt.Println("🔄 Reconnected to Binance Depth WebSocket")
		},
		OnError: func(err error) {
			fmt.Printf("❌ Depth WebSocket Error: %v\n", err)
		},
		OnDepth: func(depth binance.WSDepth) {
			// Show best bid/ask for concise output
			bestBid := "N/A"
			bestAsk := "N/A"
			if len(depth.Bids) > 0 {
				bestBid = depth.Bids[0][0] // Best bid price
			}
			if len(depth.Asks) > 0 {
				bestAsk = depth.Asks[0][0] // Best ask price
			}
			fmt.Printf("📖 Depth Update: UpdateID=%d | Best Bid: %s | Best Ask: %s | Levels: Bids=%d Asks=%d\n",
				depth.LastUpdateId, bestBid, bestAsk, len(depth.Bids), len(depth.Asks))
		},
		OnDisconnect: func() {
			fmt.Println("🔴 Disconnected from Binance Depth WebSocket")
		},
	}

	// Configure Depth Update subscription options
	depthUpdateOptions := binance.DepthUpdateSubscriptionOptions{
		OnConnect: func() {
			fmt.Println("🟢 Connected to Binance DepthUpdate WebSocket")
		},
		OnReconnect: func() {
			fmt.Println("🔄 Reconnected to Binance DepthUpdate WebSocket")
		},
		OnError: func(err error) {
			fmt.Printf("❌ DepthUpdate WebSocket Error: %v\n", err)
		},
		OnDepthUpdate: func(update binance.WSDepthUpdate) {
			fmt.Printf("🔄 Depth Update: %s | UpdateIDs: %d-%d | Bid Changes: %d | Ask Changes: %d\n",
				update.Symbol, update.FirstUpdateId, update.FinalUpdateId, len(update.BidUpdates), len(update.AskUpdates))
		},
		OnDisconnect: func() {
			fmt.Println("🔴 Disconnected from Binance DepthUpdate WebSocket")
		},
	}

	// Subscribe to BTCUSDT 1-minute klines (using lowercase symbol)
	fmt.Println("Subscribing to BTCUSDT 1m kline stream...")
	unsubscribeKline, err := wsClient.SubscribeKline("btcusdt", "1m", klineOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to kline stream: %v", err)
	}

	// Subscribe to BTCUSDT aggregate trades
	fmt.Println("Subscribing to BTCUSDT aggregate trade stream...")
	unsubscribeAggTrade, err := wsClient.SubscribeAggTrade("btcusdt", aggTradeOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to aggregate trade stream: %v", err)
	}

	// Subscribe to BTCUSDT raw trades
	fmt.Println("Subscribing to BTCUSDT raw trade stream...")
	unsubscribeTrade, err := wsClient.SubscribeTrade("btcusdt", tradeOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to trade stream: %v", err)
	}

	// Subscribe to BTCUSDT partial book depth (5 levels, 1000ms updates)
	fmt.Println("Subscribing to BTCUSDT depth stream (5 levels)...")
	unsubscribeDepth, err := wsClient.SubscribeDepth("btcusdt", 5, "", depthOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to depth stream: %v", err)
	}

	// Subscribe to BTCUSDT differential depth updates (1000ms updates)
	fmt.Println("Subscribing to BTCUSDT differential depth update stream...")
	unsubscribeDepthUpdate, err := wsClient.SubscribeDepthUpdate("btcusdt", "", depthUpdateOptions)
	if err != nil {
		log.Fatalf("Failed to subscribe to depth update stream: %v", err)
	}

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("WebSocket clients are running. Press Ctrl+C to stop...")
	fmt.Println("You should see:")
	fmt.Println("  📊 Kline updates every ~1 minute (OHLCV data)")
	fmt.Println("  💰 Aggregate trade updates (aggregated for taker orders)")
	fmt.Println("  🔥 Raw trade updates (individual buyer/seller trades)")
	fmt.Println("  📖 Depth updates every second (order book snapshots)")
	fmt.Println("  🔄 Depth update changes (incremental order book changes)")

	// Wait for shutdown signal
	<-sigChan

	fmt.Println("\nShutting down...")
	unsubscribeKline()
	unsubscribeAggTrade()
	unsubscribeTrade()
	unsubscribeDepth()
	unsubscribeDepthUpdate()
	wsClient.Close()

	// Give some time for cleanup
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Example completed.")
}

func userDataStreamExample() {
	// Configure your API credentials here (required for user data stream)
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		fmt.Println("⚠️  User Data Stream Example Skipped")
		fmt.Println("Please set BINANCE_API_KEY and BINANCE_API_SECRET environment variables to test user data streams")
		return
	}

	// Create REST API client for listen key management
	cfg := &binance.Config{
		BaseURL:   binance.MainnetBaseUrl,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
	client := binance.NewClient(cfg)

	// Create WebSocket client with REST client for user data streams
	wsClient := binance.NewWSClientWithRESTClient(binance.WSConfig{
		BaseURL: binance.MainnetWSBaseUrl9443,
	}, client)

	// Configure user data subscription options
	userDataOptions := binance.UserDataSubscriptionOptions{
		OnConnect: func() {
			fmt.Println("🟢 Connected to Binance User Data Stream")
		},
		OnReconnect: func() {
			fmt.Println("🔄 Reconnected to User Data Stream")
		},
		OnError: func(err error) {
			fmt.Printf("❌ User Data Stream Error: %v\n", err)
		},
		OnAccountPosition: func(event binance.WSOutboundAccountPositionEvent) {
			fmt.Printf("💰 Account Position Update: %d balances updated\n", len(event.BalanceArray))
			for _, balance := range event.BalanceArray {
				if balance.Free != "0.00000000" || balance.Locked != "0.00000000" {
					fmt.Printf("  %s: Free=%s, Locked=%s\n", balance.Asset, balance.Free, balance.Locked)
				}
			}
		},
		OnBalanceUpdate: func(event binance.WSBalanceUpdateEvent) {
			fmt.Printf("💸 Balance Update: %s delta=%s\n", event.Asset, event.BalanceDelta)
		},
		OnExecutionReport: func(event binance.WSExecutionReportEvent) {
			fmt.Printf("📋 Order Update: %s %s %s %s@%s (Status: %s)\n",
				event.Symbol, event.Side, event.CurrentExecutionType,
				event.OrderQuantity, event.OrderPrice, event.CurrentOrderStatus)
		},
		OnListenKeyExpired: func(event binance.WSListenKeyExpiredEvent) {
			fmt.Printf("🔑 Listen Key Expired: %s (reconnection will be handled automatically)\n", event.ListenKey)
		},
		OnDisconnect: func() {
			fmt.Println("🔴 Disconnected from User Data Stream")
		},
	}

	// Subscribe to user data stream
	unsubscribeUserData, err := wsClient.SubscribeUserData(userDataOptions)
	if err != nil {
		fmt.Printf("❌ Failed to subscribe to user data stream: %v\n", err)
		return
	}

	// Set up shutdown signal handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("User Data Stream is active. Press Ctrl+C to stop...")
	fmt.Println("You should see:")
	fmt.Println("  💰 Account position updates when balances change")
	fmt.Println("  💸 Balance updates when funds are deposited/withdrawn/transferred")
	fmt.Println("  📋 Order execution reports when orders are placed/filled/cancelled")
	fmt.Println("  🔑 Listen key expiry notifications (handled automatically)")
	fmt.Println("\nTo see user data events, try:")
	fmt.Println("  - Place a test order on Binance")
	fmt.Println("  - Transfer funds between Spot and other accounts")
	fmt.Println("  - Deposit or withdraw funds")

	// Wait for shutdown signal
	<-sigChan

	fmt.Println("\nShutting down user data stream...")
	unsubscribeUserData()
	wsClient.Close()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("User Data Stream example completed.")
}
