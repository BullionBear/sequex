package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Binance connectivity test...")

	// Create configuration
	// For testing, you can use testnet
	config := binance.DefaultConfig()

	// Set your API credentials (you'll need to set these as environment variables)
	config.APIKey = os.Getenv("BINANCE_API_KEY")
	config.APISecret = os.Getenv("BINANCE_API_SECRET")

	if config.APIKey == "" || config.APISecret == "" {
		log.Fatal("Please set BINANCE_API_KEY and BINANCE_API_SECRET environment variables")
	}

	// Create client
	client := binance.NewClient(config)
	wsClient := binance.NewWSStreamClient(config)

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

	// Test 5: Cancel the limit order
	log.Println("\n=== Test 5: Cancel Limit Order ===")
	if err := testCancelOrder(ctx, client, limitOrderID); err != nil {
		log.Printf("Failed to cancel order: %v", err)
		return
	}

	// Test 6: Send a market buy order of ADAUSDT
	log.Println("\n=== Test 6: Send Market Buy Order ===")
	if err := testMarketBuyOrder(ctx, client); err != nil {
		log.Printf("Failed to place market buy order: %v", err)
		return
	}

	// Test 7: Send a market sell order of ADAUSDT
	log.Println("\n=== Test 7: Send Market Sell Order ===")
	if err := testMarketSellOrder(ctx, client); err != nil {
		log.Printf("Failed to place market sell order: %v", err)
		return
	}

	// Test 8: Unsubscribe from user data stream
	log.Println("\n=== Test 8: Unsubscribe from User Data Stream ===")
	if err := testUnsubscribeUserDataStream(ctx, client, listenKey); err != nil {
		log.Printf("Failed to unsubscribe from user data stream: %v", err)
		return
	}

	log.Println("\n=== All tests completed successfully! ===")
}

func testGetAccountBalance(ctx context.Context, client *binance.Client) error {
	account, err := client.GetAccount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	log.Printf("Account Type: %s", account.AccountType)
	log.Printf("Can Trade: %t", account.CanTrade)
	log.Printf("Can Withdraw: %t", account.CanWithdraw)
	log.Printf("Can Deposit: %t", account.CanDeposit)
	log.Printf("Update Time: %d", account.UpdateTime)

	log.Println("Balances:")
	for _, balance := range account.Balances {
		if balance.Asset == "USDT" || balance.Asset == "ADA" {
			log.Printf("  %s: Free=%s, Locked=%s", balance.Asset, balance.Free, balance.Locked)
		}
	}

	return nil
}

func testGetADAUSDTPrice(ctx context.Context, client *binance.Client) (float64, error) {
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

func testSubscribeUserDataStream(ctx context.Context, client *binance.Client, wsClient *binance.WSStreamClient) (string, func() error, error) {
	// Create user data stream
	streamResp, err := client.CreateUserDataStream(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create user data stream: %w", err)
	}

	listenKey := streamResp.ListenKey
	log.Printf("Created user data stream with listen key: %s", listenKey)

	// Subscribe to user data stream
	userDataOptions := &binance.UserDataSubscriptionOptions{}
	userDataOptions.WithConnect(func() {
		log.Println("Connected to user data stream")
	}).WithDisconnect(func() {
		log.Println("Disconnected from user data stream")
	}).WithError(func(err error) {
		log.Printf("User data stream error: %v", err)
	}).WithAccountUpdate(func(data *binance.WSOutboundAccountPosition) error {
		log.Printf("Account update received: %+v", data)
		return nil
	}).WithBalanceUpdate(func(data *binance.WSBalanceUpdate) error {
		log.Printf("Balance update received: %+v", data)
		return nil
	}).WithExecutionReport(func(data *binance.WSExecutionReport) error {
		log.Printf("Execution report received: %+v", data)
		return nil
	})

	unsubscribe, err := wsClient.SubscribeToUserDataStream(listenKey, userDataOptions)
	if err != nil {
		return "", nil, fmt.Errorf("failed to subscribe to user data stream: %w", err)
	}

	// Wait a bit for connection to establish
	time.Sleep(2 * time.Second)

	return listenKey, unsubscribe, nil
}

func testLimitBuyOrder(ctx context.Context, client *binance.Client, currentPrice float64) (int64, error) {
	// Calculate limit price (0.9 * current price)
	limitPrice := currentPrice * 0.9

	// Calculate quantity (ensure minimum notional value but within available balance)
	// Binance requires minimum notional value, so we'll use a reasonable quantity
	quantity := 10.0 // 3 ADA to ensure minimum notional and fit available balance

	// Create limit buy order
	orderReq := &binance.NewOrderRequest{
		Symbol:      "ADAUSDT",
		Side:        binance.SideBuy,
		Type:        binance.OrderTypeLimit,
		TimeInForce: binance.TimeInForceGTC,
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

func testCancelOrder(ctx context.Context, client *binance.Client, orderID int64) error {
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

func testMarketBuyOrder(ctx context.Context, client *binance.Client) error {
	// Calculate quantity (ensure minimum notional value but within available balance)
	quantity := 10.0 // 3 ADA to ensure minimum notional and fit available balance

	// Create market buy order
	orderReq := &binance.NewOrderRequest{
		Symbol:   "ADAUSDT",
		Side:     binance.SideBuy,
		Type:     binance.OrderTypeMarket,
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
	log.Printf("  Executed Quantity: %s", order.ExecutedQty)
	log.Printf("  Cumulative Quote Quantity: %s", order.CummulativeQuoteQty)

	// Print fill information if available
	if len(order.Fills) > 0 {
		log.Printf("  Fills:")
		for i, fill := range order.Fills {
			log.Printf("    Fill %d: Price=%s, Qty=%s, Commission=%s %s",
				i+1, fill.Price, fill.Qty, fill.Commission, fill.CommissionAsset)
		}
	}

	return nil
}

func testMarketSellOrder(ctx context.Context, client *binance.Client) error {
	// Calculate quantity (ensure minimum notional value but within available balance)
	quantity := 10.0 // 3 ADA to ensure minimum notional and fit available balance

	// Create market sell order
	orderReq := &binance.NewOrderRequest{
		Symbol:   "ADAUSDT",
		Side:     binance.SideSell,
		Type:     binance.OrderTypeMarket,
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
	log.Printf("  Executed Quantity: %s", order.ExecutedQty)
	log.Printf("  Cumulative Quote Quantity: %s", order.CummulativeQuoteQty)

	// Print fill information if available
	if len(order.Fills) > 0 {
		log.Printf("  Fills:")
		for i, fill := range order.Fills {
			log.Printf("    Fill %d: Price=%s, Qty=%s, Commission=%s %s",
				i+1, fill.Price, fill.Qty, fill.Commission, fill.CommissionAsset)
		}
	}

	return nil
}

func testUnsubscribeUserDataStream(ctx context.Context, client *binance.Client, listenKey string) error {
	log.Printf("Closing user data stream with listen key: %s", listenKey)

	err := client.CloseUserDataStream(ctx, listenKey)
	if err != nil {
		return fmt.Errorf("failed to close user data stream: %w", err)
	}

	log.Println("User data stream closed successfully")
	return nil
}
