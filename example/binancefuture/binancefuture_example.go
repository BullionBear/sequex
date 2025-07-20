package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binancefuture"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Binance Futures connectivity test...")

	// Create configuration
	// For testing, you can use testnet
	config := binancefuture.DefaultConfig()

	// Set your API credentials (you'll need to set these as environment variables)
	config.APIKey = os.Getenv("BINANCE_API_KEY")
	config.APISecret = os.Getenv("BINANCE_API_SECRET")

	if config.APIKey == "" || config.APISecret == "" {
		log.Fatal("Please set BINANCE_API_KEY and BINANCE_API_SECRET environment variables")
	}

	// Create client
	client := binancefuture.NewClient(config)
	wsClient := binancefuture.NewWSStreamClient(config)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test 1: Get account balance
	log.Println("\n=== Test 1: Get Account Balance ===")
	if err := testGetAccountBalance(ctx, client); err != nil {
		log.Printf("Failed to get account balance: %v", err)
		return
	}

	// Test 2: Get ADAUSDT price
	log.Println("\n=== Test 2: Get ADAUSDT Price ===")
	adaPrice, err := testGetADAUSDTPrice(ctx, client)
	if err != nil {
		log.Printf("Failed to get ADAUSDT price: %v", err)
		return
	}

	// Test 3: Subscribe to user data stream
	log.Println("\n=== Test 3: Subscribe to User Data Stream ===")
	listenKey, unsubscribe, err := testSubscribeUserDataStream(ctx, client, wsClient)
	if err != nil {
		log.Printf("Failed to subscribe to user data stream: %v", err)
		return
	}
	defer unsubscribe()

	// Test 4: Send a limit buy order of ADAUSDT at 0.9 * current price
	log.Println("\n=== Test 4: Send Limit Buy Order ===")
	limitOrderID, err := testLimitBuyOrder(ctx, client, adaPrice)
	if err != nil {
		log.Printf("Failed to place limit buy order: %v", err)
		return
	}

	// Test 5: Get open orders
	log.Println("\n=== Test 5: Get Open Orders ===")
	if err := testGetOpenOrders(ctx, client); err != nil {
		log.Printf("Failed to get open orders: %v", err)
		return
	}

	// Test 6: Cancel the limit order
	log.Println("\n=== Test 6: Cancel Limit Order ===")
	if err := testCancelOrder(ctx, client, limitOrderID); err != nil {
		log.Printf("Failed to cancel order: %v", err)
		return
	}

	// Test 7: Send a market buy order of ADAUSDT
	log.Println("\n=== Test 7: Send Market Buy Order ===")
	if err := testMarketBuyOrder(ctx, client); err != nil {
		log.Printf("Failed to place market buy order: %v", err)
		return
	}

	// Test 8: Get open positions
	log.Println("\n=== Test 8: Get Open Positions ===")
	if err := testGetOpenPositions(ctx, client); err != nil {
		log.Printf("Failed to get open positions: %v", err)
		return
	}

	// Test 9: Send a market sell order of ADAUSDT
	log.Println("\n=== Test 9: Send Market Sell Order ===")
	if err := testMarketSellOrder(ctx, client); err != nil {
		log.Printf("Failed to place market sell order: %v", err)
		return
	}

	// Test 10: Unsubscribe from user data stream
	log.Println("\n=== Test 10: Unsubscribe from User Data Stream ===")
	if err := testUnsubscribeUserDataStream(ctx, client, listenKey); err != nil {
		log.Printf("Failed to unsubscribe from user data stream: %v", err)
		return
	}

	log.Println("\n=== All tests completed successfully! ===")
}

func testGetAccountBalance(ctx context.Context, client *binancefuture.Client) error {
	account, err := client.GetAccount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	log.Printf("Account Type: %s", account.AccountType)
	log.Printf("Can Trade: %t", account.CanTrade)
	log.Printf("Can Withdraw: %t", account.CanWithdraw)
	log.Printf("Can Deposit: %t", account.CanDeposit)
	log.Printf("Update Time: %d", account.UpdateTime)
	log.Printf("Total Wallet Balance: %s", account.TotalWalletBalance)
	log.Printf("Available Balance: %s", account.AvailableBalance)

	log.Println("Assets:")
	for _, asset := range account.Assets {
		if asset.Asset == "USDT" || asset.Asset == "ADA" {
			log.Printf("  %s: Wallet Balance=%s, Available Balance=%s",
				asset.Asset, asset.WalletBalance, asset.AvailableBalance)
		}
	}

	log.Println("Positions:")
	for _, position := range account.Positions {
		if position.Symbol == "ADAUSDT" {
			log.Printf("  %s: Position Amount=%s, Entry Price=%s, Mark Price=%s, Unrealized PnL=%s",
				position.Symbol, position.PositionAmt, position.EntryPrice,
				position.MarkPrice, position.UnrealizedProfit)
		}
	}

	return nil
}

func testGetADAUSDTPrice(ctx context.Context, client *binancefuture.Client) (float64, error) {
	ticker, err := client.GetTickerPrice(ctx, "ADAUSDT")
	if err != nil {
		return 0, fmt.Errorf("failed to get ADAUSDT ticker: %w", err)
	}

	if !ticker.IsSingle() {
		return 0, fmt.Errorf("expected single ticker, got array")
	}

	price, err := strconv.ParseFloat(ticker.GetSingle().Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	log.Printf("ADAUSDT Current Price: $%.4f", price)
	return price, nil
}

func testSubscribeUserDataStream(ctx context.Context, client *binancefuture.Client, wsClient *binancefuture.WSStreamClient) (string, func() error, error) {
	// Create user data stream
	streamResp, err := client.CreateUserDataStream(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create user data stream: %w", err)
	}

	listenKey := streamResp.ListenKey
	log.Printf("Created user data stream with listen key: %s", listenKey)

	// Subscribe to user data stream
	userDataOptions := binancefuture.NewUserDataSubscriptionOptions()
	userDataOptions.WithConnect(func() {
		log.Println("Connected to user data stream")
	}).WithDisconnect(func() {
		log.Println("Disconnected from user data stream")
	}).WithError(func(err error) {
		log.Printf("User data stream error: %v", err)
	}).WithAccountUpdateEvent(func(data *binancefuture.WSAccountUpdateEvent) error {
		log.Printf("Account update event received: %+v", data)
		return nil
	}).WithListenKeyExpired(func(data *binancefuture.WSListenKeyExpiredEvent) error {
		log.Printf("Listen key expired: %+v", data)
		return nil
	}).WithMarginCall(func(data *binancefuture.WSMarginCallEvent) error {
		log.Printf("Margin call received: %+v", data)
		return nil
	}).WithOrderTradeUpdate(func(data *binancefuture.WSOrderTradeUpdateEvent) error {
		log.Printf("Order trade update received: %+v", data)
		return nil
	}).WithAccountConfigUpdate(func(data *binancefuture.WSAccountConfigUpdateEvent) error {
		log.Printf("Account config update received: %+v", data)
		return nil
	})

	unsubscribe, err := wsClient.SubscribeToUserDataStream(userDataOptions)
	if err != nil {
		return "", nil, fmt.Errorf("failed to subscribe to user data stream: %w", err)
	}

	// Wait a bit for connection to establish
	time.Sleep(2 * time.Second)

	return listenKey, unsubscribe, nil
}

func testLimitBuyOrder(ctx context.Context, client *binancefuture.Client, currentPrice float64) (int64, error) {
	// Calculate limit price (0.9 * current price)
	limitPrice := currentPrice * 0.9

	// Calculate quantity (ensure minimum notional value but within available balance)
	// Binance Futures requires minimum notional value, so we'll use a reasonable quantity
	quantity := 10.0 // 10 ADA to ensure minimum notional and fit available balance

	// Create limit buy order
	orderReq := &binancefuture.NewOrderRequest{
		Symbol:      "ADAUSDT",
		Side:        binancefuture.SideBuy,
		Type:        binancefuture.OrderTypeLimit,
		TimeInForce: binancefuture.TimeInForceGTC,
		Quantity:    fmt.Sprintf("%.1f", quantity),
		Price:       fmt.Sprintf("%.4f", limitPrice),
	}

	log.Printf("Placing limit buy order: %s %s ADA at $%.4f", orderReq.Side, orderReq.Quantity, limitPrice)

	order, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		return 0, fmt.Errorf("failed to place limit buy order: %w", err)
	}

	log.Printf("Limit buy order placed successfully:")
	log.Printf("  Order ID: %d", order.OrderId)
	log.Printf("  Status: %s", order.Status)
	log.Printf("  Price: %s", order.Price)
	log.Printf("  Quantity: %s", order.OrigQty)

	return order.OrderId, nil
}

func testGetOpenOrders(ctx context.Context, client *binancefuture.Client) error {
	orders, err := client.GetOpenOrders(ctx, "ADAUSDT")
	if err != nil {
		return fmt.Errorf("failed to get open orders: %w", err)
	}

	log.Printf("Found %d open orders for ADAUSDT:", len(orders))
	for i, order := range orders {
		log.Printf("  Order %d:", i+1)
		log.Printf("    Order ID: %d", order.OrderId)
		log.Printf("    Symbol: %s", order.Symbol)
		log.Printf("    Side: %s", order.Side)
		log.Printf("    Type: %s", order.Type)
		log.Printf("    Status: %s", order.Status)
		log.Printf("    Price: %s", order.Price)
		log.Printf("    Quantity: %s", order.OrigQty)
		log.Printf("    Executed Quantity: %s", order.ExecutedQty)
	}

	return nil
}

func testCancelOrder(ctx context.Context, client *binancefuture.Client, orderID int64) error {
	log.Printf("Canceling order ID: %d", orderID)

	cancelResp, err := client.CancelOrder(ctx, "ADAUSDT", orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	log.Printf("Order canceled successfully:")
	log.Printf("  Order ID: %d", cancelResp.OrderId)
	log.Printf("  Status: %s", cancelResp.Status)
	log.Printf("  Executed Quantity: %s", cancelResp.ExecutedQty)

	return nil
}

func testMarketBuyOrder(ctx context.Context, client *binancefuture.Client) error {
	// Calculate quantity (ensure minimum notional value but within available balance)
	quantity := 10.0 // 10 ADA to ensure minimum notional and fit available balance

	// Create market buy order
	orderReq := &binancefuture.NewOrderRequest{
		Symbol:   "ADAUSDT",
		Side:     binancefuture.SideBuy,
		Type:     binancefuture.OrderTypeMarket,
		Quantity: fmt.Sprintf("%.1f", quantity),
	}

	log.Printf("Placing market buy order: %s %s ADA", orderReq.Side, orderReq.Quantity)

	order, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		return fmt.Errorf("failed to place market buy order: %w", err)
	}

	log.Printf("Market buy order placed successfully:")
	log.Printf("  Order ID: %d", order.OrderId)
	log.Printf("  Status: %s", order.Status)
	log.Printf("  Quantity: %s", order.OrigQty)
	log.Printf("  Executed Quantity: %s", order.ExecutedQty)
	log.Printf("  Average Price: %s", order.AvgPrice)

	return nil
}

func testGetOpenPositions(ctx context.Context, client *binancefuture.Client) error {
	positions, err := client.GetPositionRisk(ctx, "ADAUSDT")
	if err != nil {
		return fmt.Errorf("failed to get position risk: %w", err)
	}

	log.Printf("Found %d positions for ADAUSDT:", len(positions))
	for i, position := range positions {
		log.Printf("  Position %d:", i+1)
		log.Printf("    Symbol: %s", position.Symbol)
		log.Printf("    Position Amount: %s", position.PositionAmt)
		log.Printf("    Entry Price: %s", position.EntryPrice)
		log.Printf("    Mark Price: %s", position.MarkPrice)
		log.Printf("    Unrealized PnL: %s", position.UnrealizedPnl)
		log.Printf("    Leverage: %s", position.Leverage)
		log.Printf("    Margin Type: %s", position.MarginType)
		log.Printf("    Has Position: %t", position.HasPosition())
		if position.HasPosition() {
			log.Printf("    Is Long: %t", position.IsLong())
			log.Printf("    Is Short: %t", position.IsShort())
		}
	}

	return nil
}

func testMarketSellOrder(ctx context.Context, client *binancefuture.Client) error {
	// Calculate quantity (ensure minimum notional value but within available balance)
	quantity := 10.0 // 10 ADA to ensure minimum notional and fit available balance

	// Create market sell order
	orderReq := &binancefuture.NewOrderRequest{
		Symbol:   "ADAUSDT",
		Side:     binancefuture.SideSell,
		Type:     binancefuture.OrderTypeMarket,
		Quantity: fmt.Sprintf("%.1f", quantity),
	}

	log.Printf("Placing market sell order: %s %s ADA", orderReq.Side, orderReq.Quantity)

	order, err := client.PlaceOrder(ctx, orderReq)
	if err != nil {
		return fmt.Errorf("failed to place market sell order: %w", err)
	}

	log.Printf("Market sell order placed successfully:")
	log.Printf("  Order ID: %d", order.OrderId)
	log.Printf("  Status: %s", order.Status)
	log.Printf("  Quantity: %s", order.OrigQty)
	log.Printf("  Executed Quantity: %s", order.ExecutedQty)
	log.Printf("  Average Price: %s", order.AvgPrice)

	return nil
}

func testUnsubscribeUserDataStream(ctx context.Context, client *binancefuture.Client, listenKey string) error {
	log.Printf("Closing user data stream with listen key: %s", listenKey)

	err := client.CloseUserDataStream(ctx, listenKey)
	if err != nil {
		return fmt.Errorf("failed to close user data stream: %w", err)
	}

	log.Println("User data stream closed successfully")
	return nil
}
