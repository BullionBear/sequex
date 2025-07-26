package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
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

	fmt.Println("🚀 Starting User Data Stream Test with Active Trading...")

	// Create REST API client for listen key management and trading
	cfg := &binance.Config{
		BaseURL:   binance.MainnetBaseUrl,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
	client := binance.NewClient(cfg)

	// Step 1: Check USDT balance
	fmt.Println("📊 Checking account balance...")
	ctx := context.Background()
	accountResp, err := client.GetAccountInfo(ctx, binance.GetAccountInfoRequest{})
	if err != nil {
		panic(fmt.Sprintf("Failed to get account info: %v", err))
	}

	var usdtBalance float64
	var hasUSDT bool
	for _, balance := range accountResp.Data.Balances {
		if balance.Asset == "USDT" {
			// Parse free balance
			if parsed, err := strconv.ParseFloat(balance.Free, 64); err == nil {
				usdtBalance = parsed
				hasUSDT = true
			}
			break
		}
	}

	if !hasUSDT || usdtBalance < 10.0 {
		panic(fmt.Sprintf("Insufficient USDT balance: %.8f (required: >= 10.0)", usdtBalance))
	}

	fmt.Printf("✅ USDT Balance: %.8f (sufficient for testing)\n", usdtBalance)

	// Create WebSocket client with REST client for user data streams
	wsClient := binance.NewWSClientWithRESTClient(binance.WSConfig{
		BaseURL: binance.MainnetWSBaseUrl9443,
	}, client)

	// Step 2: Set up User Data Stream monitoring
	var orderCount int
	var executionCount int

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
			fmt.Printf("💰 Account Position Update: %d balances updated at %d\n",
				len(event.BalanceArray), event.EventTime)
			for _, balance := range event.BalanceArray {
				if balance.Asset == "USDT" || balance.Asset == "USDC" {
					fmt.Printf("  %s: Free=%s, Locked=%s\n", balance.Asset, balance.Free, balance.Locked)
				}
			}
		},
		OnBalanceUpdate: func(event binance.WSBalanceUpdateEvent) {
			if event.Asset == "USDT" || event.Asset == "USDC" {
				fmt.Printf("💸 Balance Update: %s delta=%s at %d\n", event.Asset, event.BalanceDelta, event.EventTime)
			}
		},
		OnExecutionReport: func(event binance.WSExecutionReportEvent) {
			executionCount++
			fmt.Printf("📋 Order #%d: %s %s %s %s@%s (Status: %s -> %s) at %d\n",
				executionCount, event.Symbol, event.Side, event.CurrentExecutionType,
				event.OrderQuantity, event.OrderPrice, event.CurrentOrderStatus,
				event.CurrentExecutionType, event.EventTime)

			if event.CurrentOrderStatus == binance.OrderStatusCanceled {
				fmt.Printf("    ✅ Order %d successfully canceled\n", event.OrderId)
			} else if event.CurrentOrderStatus == binance.OrderStatusNew {
				fmt.Printf("    📝 Order %d placed successfully\n", event.OrderId)
			}
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
		panic(fmt.Sprintf("Failed to subscribe to user data stream: %v", err))
	}

	// Wait for connection to establish
	time.Sleep(2 * time.Second)

	// Step 3: Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Step 4: Start order placement/cancellation loop
	fmt.Println("\n🎯 Starting continuous USDC/USDT order placement test...")
	fmt.Println("   Symbol: USDCUSDT")
	fmt.Println("   Side: BUY")
	fmt.Println("   Price: 0.9000")
	fmt.Println("   Quantity: 11 USDC (small test amount)")
	fmt.Println("   Strategy: Place order -> Wait 1s -> Cancel -> Repeat")
	fmt.Println("\nPress Ctrl+C to stop...\n")

	// Run order loop in a goroutine
	orderLoop := make(chan bool)
	go func() {
		defer close(orderLoop)

		for {
			select {
			case <-sigChan:
				return
			default:
				orderCount++

				// Place a small buy order for USDC/USDT
				orderReq := binance.CreateOrderRequest{
					Symbol:           "USDCUSDT",
					Side:             binance.OrderSideBuy,
					Type:             binance.OrderTypeLimit,
					TimeInForce:      binance.TimeInForceGTC,
					Quantity:         "11.0",   // Small quantity for testing
					Price:            "0.9000", // Price below market to avoid accidental fills
					NewOrderRespType: binance.NewOrderRespTypeFull,
				}

				fmt.Printf("📤 Placing order #%d...\n", orderCount)
				createResp, err := client.CreateOrder(ctx, orderReq)
				if err != nil {
					fmt.Printf("❌ Failed to place order #%d: %v\n", orderCount, err)
					time.Sleep(1 * time.Second)
					continue
				}

				orderId := createResp.Data.OrderId
				fmt.Printf("📝 Order #%d placed: ID=%d, Status=%s\n",
					orderCount, orderId, createResp.Data.Status)

				// Wait 1 second
				time.Sleep(1 * time.Second)

				// Cancel the order
				cancelReq := binance.CancelOrderRequest{
					Symbol:  "USDCUSDT",
					OrderId: orderId,
				}

				fmt.Printf("🗑️  Canceling order #%d (ID: %d)...\n", orderCount, orderId)
				cancelResp, err := client.CancelOrder(ctx, cancelReq)
				if err != nil {
					fmt.Printf("❌ Failed to cancel order #%d: %v\n", orderCount, err)
				} else {
					fmt.Printf("✅ Order #%d canceled: Status=%s\n",
						orderCount, cancelResp.Data.Status)
				}

				// Brief pause before next iteration
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan

	fmt.Println("\n🛑 Shutdown signal received...")
	fmt.Printf("📊 Test Summary:\n")
	fmt.Printf("   Orders Placed: %d\n", orderCount)
	fmt.Printf("   Execution Reports: %d\n", executionCount)

	// Cancel any remaining orders
	fmt.Println("🧹 Cleaning up any remaining orders...")
	if openOrdersResp, err := client.ListOpenOrders(ctx, binance.ListOpenOrdersRequest{
		Symbol: "USDCUSDT",
	}); err == nil && openOrdersResp.Data != nil {
		for _, order := range *openOrdersResp.Data {
			if order.Symbol == "USDCUSDT" {
				fmt.Printf("🗑️  Canceling remaining order: %d\n", order.OrderId)
				client.CancelOrder(ctx, binance.CancelOrderRequest{
					Symbol:  "USDCUSDT",
					OrderId: order.OrderId,
				})
			}
		}
	}

	fmt.Println("🔌 Shutting down user data stream...")
	unsubscribeUserData()
	wsClient.Close()

	time.Sleep(500 * time.Millisecond)
	fmt.Println("✅ User Data Stream test completed successfully!")
}
