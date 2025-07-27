package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/BullionBear/sequex/pkg/exchange/binanceperp"
)

func main() {
	fmt.Println("=== Binance Perpetual Futures REST API Example ===")
	restAPIExample()
}

func restAPIExample() {
	// Configure your API credentials here
	cfg := &binanceperp.Config{
		BaseURL: binanceperp.MainnetBaseUrl,
		// For signed requests, you would also need:
		// APIKey:    "your_api_key",
		// APISecret: "your_api_secret",
	}

	client := binanceperp.NewClient(cfg)

	// Example 1: Get Server Time (unsigned request)
	fmt.Println("\n--- Get Server Time ---")
	timeResp, err := client.GetServerTime(context.Background())
	if err != nil {
		log.Printf("GetServerTime error: %v", err)
		return
	}

	if timeResp.Code != 0 {
		log.Printf("GetServerTime failed with code %d: %s", timeResp.Code, timeResp.Message)
		return
	}

	fmt.Printf("Server Time: %d\n", timeResp.Data.ServerTime)
	fmt.Printf("Response Code: %d\n", timeResp.Code)
	fmt.Printf("Response Message: %s\n", timeResp.Message)

	// Example 2: Get Order Book Depth (unsigned request with parameters)
	fmt.Println("\n--- Get Order Book Depth ---")
	depthReq := binanceperp.GetDepthRequest{
		Symbol: "BTCUSDT",
		Limit:  5,
	}
	depthResp, err := client.GetDepth(context.Background(), depthReq)
	if err != nil {
		log.Printf("GetDepth error: %v", err)
		return
	}

	if depthResp.Code != 0 {
		log.Printf("GetDepth failed with code %d: %s", depthResp.Code, depthResp.Message)
		return
	}

	fmt.Printf("Last Update ID: %d\n", depthResp.Data.LastUpdateId)
	fmt.Printf("Message Output Time: %d\n", depthResp.Data.E)
	fmt.Printf("Transaction Time: %d\n", depthResp.Data.T)
	fmt.Printf("Number of Bids: %d\n", len(depthResp.Data.Bids))
	fmt.Printf("Number of Asks: %d\n", len(depthResp.Data.Asks))

	if len(depthResp.Data.Bids) > 0 {
		fmt.Printf("Best Bid: %s @ %s\n", depthResp.Data.Bids[0][1], depthResp.Data.Bids[0][0])
	}
	if len(depthResp.Data.Asks) > 0 {
		fmt.Printf("Best Ask: %s @ %s\n", depthResp.Data.Asks[0][1], depthResp.Data.Asks[0][0])
	}

	// Example 3: Get Recent Trades (unsigned request with parameters)
	fmt.Println("\n--- Get Recent Trades ---")
	tradesReq := binanceperp.GetRecentTradesRequest{
		Symbol: "BTCUSDT",
		Limit:  3,
	}
	tradesResp, err := client.GetRecentTrades(context.Background(), tradesReq)
	if err != nil {
		log.Printf("GetRecentTrades error: %v", err)
		return
	}

	if tradesResp.Code != 0 {
		log.Printf("GetRecentTrades failed with code %d: %s", tradesResp.Code, tradesResp.Message)
		return
	}

	fmt.Printf("Number of Recent Trades: %d\n", len(*tradesResp.Data))
	for i, trade := range *tradesResp.Data {
		fmt.Printf("Trade %d:\n", i+1)
		fmt.Printf("  ID: %d\n", trade.Id)
		fmt.Printf("  Price: %s\n", trade.Price)
		fmt.Printf("  Quantity: %s\n", trade.Qty)
		fmt.Printf("  Quote Quantity: %s\n", trade.QuoteQty)
		fmt.Printf("  Time: %d\n", trade.Time)
		fmt.Printf("  Is Buyer Maker: %t\n", trade.IsBuyerMaker)
	}

	// Example 4: Get Aggregate Trades (unsigned request with parameters)
	fmt.Println("\n--- Get Aggregate Trades ---")
	aggTradesReq := binanceperp.GetAggTradesRequest{
		Symbol: "BTCUSDT",
		Limit:  3,
	}
	aggTradesResp, err := client.GetAggTrades(context.Background(), aggTradesReq)
	if err != nil {
		log.Printf("GetAggTrades error: %v", err)
		return
	}

	if aggTradesResp.Code != 0 {
		log.Printf("GetAggTrades failed with code %d: %s", aggTradesResp.Code, aggTradesResp.Message)
		return
	}

	fmt.Printf("Number of Aggregate Trades: %d\n", len(*aggTradesResp.Data))
	for i, trade := range *aggTradesResp.Data {
		fmt.Printf("Aggregate Trade %d:\n", i+1)
		fmt.Printf("  Agg Trade ID: %d\n", trade.AggTradeId)
		fmt.Printf("  Price: %s\n", trade.Price)
		fmt.Printf("  Quantity: %s\n", trade.Quantity)
		fmt.Printf("  First Trade ID: %d\n", trade.FirstTradeId)
		fmt.Printf("  Last Trade ID: %d\n", trade.LastTradeId)
		fmt.Printf("  Timestamp: %d\n", trade.Timestamp)
		fmt.Printf("  Is Buyer Maker: %t\n", trade.IsBuyerMaker)
	}

	// Example 5: Get Klines/Candlestick Data (unsigned request with parameters)
	fmt.Println("\n--- Get Klines/Candlestick Data ---")
	klinesReq := binanceperp.GetKlinesRequest{
		Symbol:   "BTCUSDT",
		Interval: "1m",
		Limit:    3,
	}
	klinesResp, err := client.GetKlines(context.Background(), klinesReq)
	if err != nil {
		log.Printf("GetKlines error: %v", err)
		return
	}

	if klinesResp.Code != 0 {
		log.Printf("GetKlines failed with code %d: %s", klinesResp.Code, klinesResp.Message)
		return
	}

	fmt.Printf("Number of Klines: %d\n", len(*klinesResp.Data))
	for i, kline := range *klinesResp.Data {
		fmt.Printf("Kline %d:\n", i+1)
		fmt.Printf("  Open Time: %d\n", kline.OpenTime)
		fmt.Printf("  Open: %s\n", kline.Open)
		fmt.Printf("  High: %s\n", kline.High)
		fmt.Printf("  Low: %s\n", kline.Low)
		fmt.Printf("  Close: %s\n", kline.Close)
		fmt.Printf("  Volume: %s\n", kline.Volume)
		fmt.Printf("  Close Time: %d\n", kline.CloseTime)
		fmt.Printf("  Quote Asset Volume: %s\n", kline.QuoteAssetVolume)
		fmt.Printf("  Number of Trades: %d\n", kline.NumberOfTrades)
		fmt.Printf("  Taker Buy Base Asset Volume: %s\n", kline.TakerBuyBaseAssetVolume)
		fmt.Printf("  Taker Buy Quote Asset Volume: %s\n", kline.TakerBuyQuoteAssetVolume)
	}

	// Example 6: Get Mark Price and Funding Rate (unsigned request with optional parameters)
	fmt.Println("\n--- Get Mark Price (Single Symbol) ---")
	markPriceReq := binanceperp.GetMarkPriceRequest{
		Symbol: "BTCUSDT",
	}
	markPriceResp, err := client.GetMarkPrice(context.Background(), markPriceReq)
	if err != nil {
		log.Printf("GetMarkPrice error: %v", err)
		return
	}

	if markPriceResp.Code != 0 {
		log.Printf("GetMarkPrice failed with code %d: %s", markPriceResp.Code, markPriceResp.Message)
		return
	}

	fmt.Printf("Number of Mark Price Entries: %d\n", len(*markPriceResp.Data))
	for i, mp := range *markPriceResp.Data {
		fmt.Printf("Mark Price %d:\n", i+1)
		fmt.Printf("  Symbol: %s\n", mp.Symbol)
		fmt.Printf("  Mark Price: %s\n", mp.MarkPrice)
		fmt.Printf("  Index Price: %s\n", mp.IndexPrice)
		fmt.Printf("  Estimated Settle Price: %s\n", mp.EstimatedSettlePrice)
		fmt.Printf("  Last Funding Rate: %s\n", mp.LastFundingRate)
		fmt.Printf("  Interest Rate: %s\n", mp.InterestRate)
		fmt.Printf("  Next Funding Time: %d\n", mp.NextFundingTime)
		fmt.Printf("  Time: %d\n", mp.Time)
	}

	// Example 7: Get Price Ticker (unsigned request with optional parameters)
	fmt.Println("\n--- Get Price Ticker (Single Symbol) ---")
	priceTickerReq := binanceperp.GetPriceTickerRequest{
		Symbol: "BTCUSDT",
	}
	priceTickerResp, err := client.GetPriceTicker(context.Background(), priceTickerReq)
	if err != nil {
		log.Printf("GetPriceTicker error: %v", err)
		return
	}

	if priceTickerResp.Code != 0 {
		log.Printf("GetPriceTicker failed with code %d: %s", priceTickerResp.Code, priceTickerResp.Message)
		return
	}

	fmt.Printf("Number of Price Ticker Entries: %d\n", len(*priceTickerResp.Data))
	for i, ticker := range *priceTickerResp.Data {
		fmt.Printf("Price Ticker %d:\n", i+1)
		fmt.Printf("  Symbol: %s\n", ticker.Symbol)
		fmt.Printf("  Price: %s\n", ticker.Price)
		fmt.Printf("  Time: %d\n", ticker.Time)
	}

	// Example 8: Get Book Ticker (unsigned request with optional parameters)
	fmt.Println("\n--- Get Book Ticker (Single Symbol) ---")
	bookTickerReq := binanceperp.GetBookTickerRequest{
		Symbol: "BTCUSDT",
	}
	bookTickerResp, err := client.GetBookTicker(context.Background(), bookTickerReq)
	if err != nil {
		log.Printf("GetBookTicker error: %v", err)
		return
	}

	if bookTickerResp.Code != 0 {
		log.Printf("GetBookTicker failed with code %d: %s", bookTickerResp.Code, bookTickerResp.Message)
		return
	}

	fmt.Printf("Number of Book Ticker Entries: %d\n", len(*bookTickerResp.Data))
	for i, ticker := range *bookTickerResp.Data {
		fmt.Printf("Book Ticker %d:\n", i+1)
		fmt.Printf("  Symbol: %s\n", ticker.Symbol)
		fmt.Printf("  Best Bid: %s @ %s\n", ticker.BidQty, ticker.BidPrice)
		fmt.Printf("  Best Ask: %s @ %s\n", ticker.AskQty, ticker.AskPrice)
		fmt.Printf("  Time: %d\n", ticker.Time)
	}

	// Example 9: Get Account Balance (signed request - requires API credentials)
	fmt.Println("\n--- Get Account Balance (Signed Request) ---")
	// This example requires API credentials to be set in the config
	// Uncomment and set your credentials to test this

	// Update config with credentials for signed requests
	cfg.APIKey = os.Getenv("BINANCEPERP_API_KEY")
	cfg.APISecret = os.Getenv("BINANCEPERP_API_SECRET")

	balanceReq := binanceperp.GetAccountBalanceRequest{
		RecvWindow: 5000,
	}
	balanceResp, err := client.GetAccountBalance(context.Background(), balanceReq)
	if err != nil {
		log.Printf("GetAccountBalance error: %v", err)
		fmt.Println("Skipping account balance - requires valid API credentials")
	} else {

		if balanceResp.Code != 0 {
			log.Printf("GetAccountBalance failed with code %d: %s", balanceResp.Code, balanceResp.Message)
		} else {

			fmt.Printf("Number of Account Balance Entries: %d\n", len(*balanceResp.Data))
			for i, balance := range *balanceResp.Data {
				fmt.Printf("Balance %d:\n", i+1)
				fmt.Printf("  Account Alias: %s\n", balance.AccountAlias)
				fmt.Printf("  Asset: %s\n", balance.Asset)
				fmt.Printf("  Balance: %s\n", balance.Balance)
				fmt.Printf("  Cross Wallet Balance: %s\n", balance.CrossWalletBalance)
				fmt.Printf("  Cross UnPnL: %s\n", balance.CrossUnPnl)
				fmt.Printf("  Available Balance: %s\n", balance.AvailableBalance)
				fmt.Printf("  Max Withdraw Amount: %s\n", balance.MaxWithdrawAmount)
				fmt.Printf("  Margin Available: %t\n", balance.MarginAvailable)
				fmt.Printf("  Update Time: %d\n", balance.UpdateTime)
				fmt.Println()
			}
		}
	}

	// Example 10: Create Order (TRADE - signed request - requires API credentials)
	fmt.Println("\n--- Create Order (TRADE - Signed Request) ---")
	// WARNING: This creates a real order! Uncomment only for testing on testnet

	// Example: Create a LIMIT BUY order
	createOrderReq := binanceperp.CreateOrderRequest{
		Symbol:      "BTCUSDT",
		Side:        binanceperp.OrderSideBuy,
		Type:        binanceperp.OrderTypeLimit,
		TimeInForce: binanceperp.TimeInForceGTC,
		Quantity:    "0.001",
		Price:       "100000.00",
		RecvWindow:  5000,
	}

	createResp, err := client.CreateOrder(context.Background(), createOrderReq)
	if err != nil {
		log.Printf("CreateOrder error: %v", err)
		return
	}

	if createResp.Code != 0 {
		log.Printf("CreateOrder failed with code %d: %s", createResp.Code, createResp.Message)
		return
	}

	fmt.Printf("Order Created Successfully:\n")
	fmt.Printf("  Order ID: %d\n", createResp.Data.OrderId)
	fmt.Printf("  Client Order ID: %s\n", createResp.Data.ClientOrderId)
	fmt.Printf("  Symbol: %s\n", createResp.Data.Symbol)
	fmt.Printf("  Side: %s\n", createResp.Data.Side)
	fmt.Printf("  Type: %s\n", createResp.Data.Type)
	fmt.Printf("  Status: %s\n", createResp.Data.Status)
	fmt.Printf("  Quantity: %s\n", createResp.Data.OrigQty)
	fmt.Printf("  Price: %s\n", createResp.Data.Price)

	// Example 11: Cancel Order (TRADE - signed request - requires API credentials)
	fmt.Println("\n--- Cancel Order (TRADE - Signed Request) ---")

	// Cancel the order created above using order ID
	cancelOrderReq := binanceperp.CancelOrderRequest{
		Symbol:  "BTCUSDT",
		OrderId: createResp.Data.OrderId,
	}

	cancelResp, err := client.CancelOrder(context.Background(), cancelOrderReq)
	if err != nil {
		log.Printf("CancelOrder error: %v", err)
		return
	}

	if cancelResp.Code != 0 {
		log.Printf("CancelOrder failed with code %d: %s", cancelResp.Code, cancelResp.Message)
		return
	}

	fmt.Printf("Order Canceled Successfully:\n")
	fmt.Printf("  Order ID: %d\n", cancelResp.Data.OrderId)
	fmt.Printf("  Client Order ID: %s\n", cancelResp.Data.ClientOrderId)
	fmt.Printf("  Symbol: %s\n", cancelResp.Data.Symbol)
	fmt.Printf("  Status: %s\n", cancelResp.Data.Status)

	// Example 12: Cancel All Orders (TRADE - signed request - requires API credentials)
	fmt.Println("\n--- Cancel All Orders (TRADE - Signed Request) ---")

	// Cancel all open orders for a symbol
	cancelAllReq := binanceperp.CancelAllOrdersRequest{
		Symbol: "BTCUSDT",
	}

	cancelAllResp, err := client.CancelAllOrders(context.Background(), cancelAllReq)
	if err != nil {
		log.Printf("CancelAllOrders error: %v", err)
		return
	}

	if cancelAllResp.Code != 0 {
		log.Printf("CancelAllOrders failed with code %d: %s", cancelAllResp.Code, cancelAllResp.Message)
		return
	}

	fmt.Printf("All Orders Canceled Successfully:\n")
	fmt.Printf("  Code: %d\n", cancelAllResp.Data.Code)
	fmt.Printf("  Message: %s\n", cancelAllResp.Data.Msg)

	// Example 13: Query Order (USER_DATA - signed request - requires API credentials)
	fmt.Println("\n--- Query Order (USER_DATA - Signed Request) ---")

	// Query an order by order ID
	queryOrderReq := binanceperp.QueryOrderRequest{
		Symbol:  "BTCUSDT",
		OrderId: createResp.Data.OrderId, // Use the order ID from create order
	}

	queryResp, err := client.QueryOrder(context.Background(), queryOrderReq)
	if err != nil {
		log.Printf("QueryOrder error: %v", err)
		return
	}

	if queryResp.Code != 0 {
		log.Printf("QueryOrder failed with code %d: %s", queryResp.Code, queryResp.Message)
		return
	}

	fmt.Printf("Order Status Query Successful:\n")
	fmt.Printf("  Order ID: %d\n", queryResp.Data.OrderId)
	fmt.Printf("  Client Order ID: %s\n", queryResp.Data.ClientOrderId)
	fmt.Printf("  Symbol: %s\n", queryResp.Data.Symbol)
	fmt.Printf("  Side: %s\n", queryResp.Data.Side)
	fmt.Printf("  Type: %s\n", queryResp.Data.Type)
	fmt.Printf("  Status: %s\n", queryResp.Data.Status)
	fmt.Printf("  Original Quantity: %s\n", queryResp.Data.OrigQty)
	fmt.Printf("  Executed Quantity: %s\n", queryResp.Data.ExecutedQty)
	fmt.Printf("  Price: %s\n", queryResp.Data.Price)
	fmt.Printf("  Average Price: %s\n", queryResp.Data.AvgPrice)
	fmt.Printf("  Order Time: %d\n", queryResp.Data.Time)
	fmt.Printf("  Update Time: %d\n", queryResp.Data.UpdateTime)

	// Example 14: Query Current Open Order (USER_DATA - signed request - requires API credentials)
	fmt.Println("\n--- Query Current Open Order (USER_DATA - Signed Request) ---")
	/*
		// Query a current open order by order ID
		// This will only work if the order is still open (not filled or cancelled)
		queryOpenOrderReq := binanceperp.QueryCurrentOpenOrderRequest{
			Symbol:  "BTCUSDT",
			OrderId: createResp.Data.OrderId,  // Use the order ID from create order
		}

		queryOpenResp, err := client.QueryCurrentOpenOrder(context.Background(), queryOpenOrderReq)
		if err != nil {
			log.Printf("QueryCurrentOpenOrder error: %v", err)
			return
		}

		if queryOpenResp.Code != 0 {
			log.Printf("QueryCurrentOpenOrder failed with code %d: %s", queryOpenResp.Code, queryOpenResp.Message)
			return
		}

		fmt.Printf("Open Order Query Successful:\n")
		fmt.Printf("  Order ID: %d\n", queryOpenResp.Data.OrderId)
		fmt.Printf("  Client Order ID: %s\n", queryOpenResp.Data.ClientOrderId)
		fmt.Printf("  Symbol: %s\n", queryOpenResp.Data.Symbol)
		fmt.Printf("  Side: %s\n", queryOpenResp.Data.Side)
		fmt.Printf("  Type: %s\n", queryOpenResp.Data.Type)
		fmt.Printf("  Status: %s\n", queryOpenResp.Data.Status)
		fmt.Printf("  Original Quantity: %s\n", queryOpenResp.Data.OrigQty)
		fmt.Printf("  Executed Quantity: %s\n", queryOpenResp.Data.ExecutedQty)
		fmt.Printf("  Price: %s\n", queryOpenResp.Data.Price)
		fmt.Printf("  Average Price: %s\n", queryOpenResp.Data.AvgPrice)
		fmt.Printf("  Order Time: %d\n", queryOpenResp.Data.Time)
		fmt.Printf("  Update Time: %d\n", queryOpenResp.Data.UpdateTime)
	*/
	fmt.Println("Query Current Open Order example is commented out - requires existing OPEN order!")

	// Example 15: Get My Trades (USER_DATA - signed request - requires API credentials)
	fmt.Println("\n--- Get My Trades (USER_DATA - Signed Request) ---")
	/*
		// Get recent trades for a symbol
		myTradesReq := binanceperp.GetMyTradesRequest{
			Symbol: "BTCUSDT",
			Limit:  10, // Get last 10 trades
		}

		myTradesResp, err := client.GetMyTrades(context.Background(), myTradesReq)
		if err != nil {
			log.Printf("GetMyTrades error: %v", err)
			return
		}

		if myTradesResp.Code != 0 {
			log.Printf("GetMyTrades failed with code %d: %s", myTradesResp.Code, myTradesResp.Message)
			return
		}

		fmt.Printf("My Trades Retrieved Successfully:\n")
		fmt.Printf("Number of Trades: %d\n", len(*myTradesResp.Data))
		for i, trade := range *myTradesResp.Data {
			fmt.Printf("\nTrade %d:\n", i+1)
			fmt.Printf("  Trade ID: %d\n", trade.Id)
			fmt.Printf("  Order ID: %d\n", trade.OrderId)
			fmt.Printf("  Symbol: %s\n", trade.Symbol)
			fmt.Printf("  Side: %s\n", trade.Side)
			fmt.Printf("  Position Side: %s\n", trade.PositionSide)
			fmt.Printf("  Price: %s\n", trade.Price)
			fmt.Printf("  Quantity: %s\n", trade.Qty)
			fmt.Printf("  Quote Quantity: %s\n", trade.QuoteQty)
			fmt.Printf("  Commission: %s %s\n", trade.Commission, trade.CommissionAsset)
			fmt.Printf("  Realized PnL: %s\n", trade.RealizedPnl)
			fmt.Printf("  Maker: %t\n", trade.Maker)
			fmt.Printf("  Buyer: %t\n", trade.Buyer)
			fmt.Printf("  Time: %d\n", trade.Time)
		}
	*/
	fmt.Println("Get My Trades example is commented out - requires API credentials!")

	// Example 16: Get Positions (USER_DATA - signed request - requires API credentials)
	fmt.Println("\n--- Get Positions (USER_DATA - Signed Request) ---")
	/*
		// Get all current positions
		positionsReq := binanceperp.GetPositionsRequest{
			// Leave empty to get all positions
			// Symbol: "BTCUSDT", // Uncomment to get specific symbol position
		}

		positionsResp, err := client.GetPositions(context.Background(), positionsReq)
		if err != nil {
			log.Printf("GetPositions error: %v", err)
			return
		}

		if positionsResp.Code != 0 {
			log.Printf("GetPositions failed with code %d: %s", positionsResp.Code, positionsResp.Message)
			return
		}

		fmt.Printf("Positions Retrieved Successfully:\n")
		fmt.Printf("Number of Positions: %d\n", len(*positionsResp.Data))

		// Filter and display non-zero positions
		activePositions := 0
		for i, position := range *positionsResp.Data {
			// Only show positions with non-zero amount
			if position.PositionAmt != "0" {
				activePositions++
				fmt.Printf("\nActive Position %d:\n", activePositions)
				fmt.Printf("  Symbol: %s\n", position.Symbol)
				fmt.Printf("  Position Side: %s\n", position.PositionSide)
				fmt.Printf("  Position Amount: %s\n", position.PositionAmt)
				fmt.Printf("  Entry Price: %s\n", position.EntryPrice)
				fmt.Printf("  Break Even Price: %s\n", position.BreakEvenPrice)
				fmt.Printf("  Mark Price: %s\n", position.MarkPrice)
				fmt.Printf("  Unrealized PnL: %s\n", position.UnRealizedProfit)
				fmt.Printf("  Liquidation Price: %s\n", position.LiquidationPrice)
				fmt.Printf("  Leverage: %s\n", position.Leverage)
				fmt.Printf("  Margin Type: %s\n", position.MarginType)
				fmt.Printf("  Update Time: %d\n", position.UpdateTime)
			} else if i < 3 { // Show first few zero positions as examples
				fmt.Printf("\nPosition %d (Zero):\n", i+1)
				fmt.Printf("  Symbol: %s\n", position.Symbol)
				fmt.Printf("  Position Side: %s\n", position.PositionSide)
				fmt.Printf("  Position Amount: %s\n", position.PositionAmt)
			}
		}

		if activePositions == 0 {
			fmt.Println("No active positions found.")
		}
	*/
	fmt.Println("Get Positions example is commented out - requires API credentials!")

	fmt.Println("\nTo test trading functions:")
	fmt.Println("1. Set your API credentials in environment variables")
	fmt.Println("2. Uncomment the trading example code")
	fmt.Println("3. Use testnet for safe testing: https://testnet.binancefuture.com")
	fmt.Println("4. Be careful with real trading - orders can lose money!")
	fmt.Println("5. Cancel All Orders will close ALL your open orders for the symbol!")
	fmt.Println("6. Query Order requires an existing order ID to check status!")
	fmt.Println("7. Query Current Open Order only works for orders that are still open!")
	fmt.Println("8. Get My Trades shows your historical trading activity!")
	fmt.Println("9. Get Positions shows your current futures positions and P&L!")
}

func aggTradeExample() {
	fmt.Println("=== AggTrade Stream Example ===")

	client := binanceperp.NewWSClient(nil)

	var aggTradeCount int64

	options := &binanceperp.AggTradeSubscriptionOptions{}
	options.
		WithConnect(func() {
			fmt.Println("✓ Connected to AggTrade stream")
		}).
		WithAggTrade(func(aggTrade binanceperp.WSAggTrade) {
			count := atomic.AddInt64(&aggTradeCount, 1)
			fmt.Printf("AggTrade #%d: %s @ %s, Qty: %s, Buyer Maker: %t\n",
				count, aggTrade.Symbol, aggTrade.Price, aggTrade.Quantity, aggTrade.IsBuyerMaker)
		}).
		WithError(func(err error) {
			fmt.Printf("AggTrade Error: %v\n", err)
		}).
		WithDisconnect(func() {
			fmt.Println("✓ AggTrade stream disconnected")
		})

	unsubscribe, err := client.SubscribeAggTrade("btcusdt", options)
	if err != nil {
		fmt.Printf("Failed to subscribe to AggTrade: %v\n", err)
		return
	}

	fmt.Println("Listening to AggTrade stream for 10 seconds...")
	time.Sleep(10 * time.Second)

	unsubscribe()
	fmt.Printf("Received %d aggregate trades\n", atomic.LoadInt64(&aggTradeCount))
}

func mixedStreamsExample() {
	fmt.Println("=== Mixed Streams Example ===")

	client := binanceperp.NewWSClient(nil)

	var klineCount int64
	var aggTradeCount int64

	// Subscribe to Kline stream
	klineOptions := &binanceperp.KlineSubscriptionOptions{}
	klineOptions.
		WithConnect(func() {
			fmt.Println("✓ Connected to Kline stream")
		}).
		WithKline(func(kline binanceperp.WSKline) {
			count := atomic.AddInt64(&klineCount, 1)
			fmt.Printf("Kline #%d: %s %s OHLC: %s/%s/%s/%s, Closed: %t\n",
				count, kline.Symbol, kline.Interval, kline.Open, kline.High, kline.Low, kline.Close, kline.IsClosed)
		}).
		WithDisconnect(func() {
			fmt.Println("✓ Kline stream disconnected")
		})

	unsubscribeKline, err := client.SubscribeKline("btcusdt", "1m", klineOptions)
	if err != nil {
		fmt.Printf("Failed to subscribe to Kline: %v\n", err)
		return
	}

	// Subscribe to AggTrade stream
	aggTradeOptions := &binanceperp.AggTradeSubscriptionOptions{}
	aggTradeOptions.
		WithConnect(func() {
			fmt.Println("✓ Connected to AggTrade stream")
		}).
		WithAggTrade(func(aggTrade binanceperp.WSAggTrade) {
			count := atomic.AddInt64(&aggTradeCount, 1)
			if count <= 5 { // Only show first 5 to avoid spam
				fmt.Printf("AggTrade #%d: %s @ %s, Qty: %s\n",
					count, aggTrade.Symbol, aggTrade.Price, aggTrade.Quantity)
			}
		}).
		WithDisconnect(func() {
			fmt.Println("✓ AggTrade stream disconnected")
		})

	unsubscribeAggTrade, err := client.SubscribeAggTrade("btcusdt", aggTradeOptions)
	if err != nil {
		fmt.Printf("Failed to subscribe to AggTrade: %v\n", err)
		unsubscribeKline()
		return
	}

	fmt.Printf("Active subscriptions: %d\n", client.GetSubscriptionCount())
	fmt.Println("Listening to both streams for 8 seconds...")
	time.Sleep(8 * time.Second)

	unsubscribeKline()
	unsubscribeAggTrade()

	fmt.Printf("Final counts - Klines: %d, AggTrades: %d\n",
		atomic.LoadInt64(&klineCount), atomic.LoadInt64(&aggTradeCount))
}

func depthExample() {
	fmt.Println("\n=== Binance Perpetual Futures Depth Stream Example ===")

	// Create WebSocket client
	client := binanceperp.NewWSClient(nil) // Use default config

	// Counter for depth updates
	var depthCount int64

	// Create depth subscription options with chain method
	options := &binanceperp.DepthSubscriptionOptions{}
	options.
		WithConnect(func() {
			fmt.Println("✓ Connected to depth stream")
		}).
		WithDepth(func(depth binanceperp.WSDepth) {
			count := atomic.AddInt64(&depthCount, 1)

			fmt.Printf("Depth Update #%d:\n", count)
			fmt.Printf("  Symbol: %s\n", depth.Symbol)
			fmt.Printf("  Event Time: %d\n", depth.EventTime)
			fmt.Printf("  Transaction Time: %d\n", depth.TransactionTime)
			fmt.Printf("  Update IDs: %d -> %d (prev: %d)\n",
				depth.FirstUpdateID, depth.FinalUpdateID, depth.PrevUpdateID)

			// Show top 3 bids and asks
			fmt.Printf("  Top Bids (%d total):\n", len(depth.Bids))
			for i, bid := range depth.Bids {
				if i >= 3 {
					break
				}
				fmt.Printf("    [%d] Price: %s, Quantity: %s\n", i+1, bid[0], bid[1])
			}

			fmt.Printf("  Top Asks (%d total):\n", len(depth.Asks))
			for i, ask := range depth.Asks {
				if i >= 3 {
					break
				}
				fmt.Printf("    [%d] Price: %s, Quantity: %s\n", i+1, ask[0], ask[1])
			}
			fmt.Println("  ---")
		}).
		WithError(func(err error) {
			fmt.Printf("❌ WebSocket Error: %v\n", err)
		}).
		WithDisconnect(func() {
			fmt.Println("✓ Disconnected from depth stream")
		})

	// Subscribe to BTC/USDT depth with Level 5 and 250ms updates (default)
	unsubscribe, err := client.SubscribeDepth("btcusdt", binanceperp.DepthLevel5, binanceperp.DepthUpdate250ms, options)
	if err != nil {
		fmt.Printf("Failed to subscribe to depth stream: %v\n", err)
		return
	}

	fmt.Println("Listening to depth5 stream for 10 seconds...")
	time.Sleep(10 * time.Second)

	unsubscribe()

	fmt.Printf("Final depth update count: %d\n", atomic.LoadInt64(&depthCount))
}

func multiDepthExample() {
	fmt.Println("\n=== Multiple Depth Levels Example ===")

	client := binanceperp.NewWSClient(nil)

	var depth5Count int64
	var depth20Count int64

	// Subscribe to depth5 with 250ms updates
	options5 := &binanceperp.DepthSubscriptionOptions{}
	options5.
		WithConnect(func() {
			fmt.Println("✓ Connected to depth5 stream")
		}).
		WithDepth(func(depth binanceperp.WSDepth) {
			count := atomic.AddInt64(&depth5Count, 1)
			fmt.Printf("Depth5 #%d: %d bids, %d asks\n", count, len(depth.Bids), len(depth.Asks))
		}).
		WithDisconnect(func() {
			fmt.Println("✓ Disconnected from depth5 stream")
		})

	// Subscribe to depth20 with 100ms updates (faster)
	options20 := &binanceperp.DepthSubscriptionOptions{}
	options20.
		WithConnect(func() {
			fmt.Println("✓ Connected to depth20 stream (100ms)")
		}).
		WithDepth(func(depth binanceperp.WSDepth) {
			count := atomic.AddInt64(&depth20Count, 1)
			fmt.Printf("Depth20 #%d: %d bids, %d asks\n", count, len(depth.Bids), len(depth.Asks))
		}).
		WithDisconnect(func() {
			fmt.Println("✓ Disconnected from depth20 stream")
		})

	// Subscribe to both streams
	unsubscribe5, err := client.SubscribeDepth("ethusdt", binanceperp.DepthLevel5, binanceperp.DepthUpdate250ms, options5)
	if err != nil {
		fmt.Printf("Failed to subscribe to depth5: %v\n", err)
		return
	}

	unsubscribe20, err := client.SubscribeDepth("ethusdt", binanceperp.DepthLevel20, binanceperp.DepthUpdate100ms, options20)
	if err != nil {
		fmt.Printf("Failed to subscribe to depth20: %v\n", err)
		return
	}

	fmt.Println("Listening to both depth streams for 8 seconds...")
	time.Sleep(8 * time.Second)

	unsubscribe5()
	unsubscribe20()

	fmt.Printf("Final counts - Depth5: %d, Depth20: %d\n",
		atomic.LoadInt64(&depth5Count), atomic.LoadInt64(&depth20Count))
}

func diffDepthExample() {
	fmt.Println("\n=== Binance Perpetual Futures Differential Depth Stream Example ===")

	// Create WebSocket client
	client := binanceperp.NewWSClient(nil) // Use default config

	// Counter for differential depth updates
	var diffDepthCount int64

	// Create differential depth subscription options with chain method
	options := &binanceperp.DiffDepthSubscriptionOptions{}
	options.
		WithConnect(func() {
			fmt.Println("✓ Connected to differential depth stream")
		}).
		WithDiffDepth(func(diffDepth binanceperp.WSDepth) {
			count := atomic.AddInt64(&diffDepthCount, 1)

			fmt.Printf("Differential Depth Update #%d:\n", count)
			fmt.Printf("  Symbol: %s\n", diffDepth.Symbol)
			fmt.Printf("  Event Time: %d\n", diffDepth.EventTime)
			fmt.Printf("  Transaction Time: %d\n", diffDepth.TransactionTime)
			fmt.Printf("  Update IDs: %d -> %d (prev: %d)\n",
				diffDepth.FirstUpdateID, diffDepth.FinalUpdateID, diffDepth.PrevUpdateID)

			// Show bid/ask changes (can be many changes, not limited to top N)
			fmt.Printf("  Bid Changes: %d entries\n", len(diffDepth.Bids))
			if len(diffDepth.Bids) > 0 {
				fmt.Printf("    Sample: Price %s -> Quantity %s\n",
					diffDepth.Bids[0][0], diffDepth.Bids[0][1])
				if diffDepth.Bids[0][1] == "0" {
					fmt.Printf("    (Quantity 0 means this price level was removed)\n")
				}
			}

			fmt.Printf("  Ask Changes: %d entries\n", len(diffDepth.Asks))
			if len(diffDepth.Asks) > 0 {
				fmt.Printf("    Sample: Price %s -> Quantity %s\n",
					diffDepth.Asks[0][0], diffDepth.Asks[0][1])
				if diffDepth.Asks[0][1] == "0" {
					fmt.Printf("    (Quantity 0 means this price level was removed)\n")
				}
			}
			fmt.Println("  ---")
		}).
		WithError(func(err error) {
			fmt.Printf("❌ WebSocket Error: %v\n", err)
		}).
		WithDisconnect(func() {
			fmt.Println("✓ Disconnected from differential depth stream")
		})

	// Subscribe to BTC/USDT differential depth with 250ms updates (default)
	unsubscribe, err := client.SubscribeDiffDepth("btcusdt", binanceperp.DepthUpdate250ms, options)
	if err != nil {
		fmt.Printf("Failed to subscribe to differential depth stream: %v\n", err)
		return
	}

	fmt.Println("Listening to differential depth stream for 8 seconds...")
	time.Sleep(8 * time.Second)

	unsubscribe()

	fmt.Printf("Final differential depth update count: %d\n", atomic.LoadInt64(&diffDepthCount))
}

func depthComparisonExample() {
	fmt.Println("\n=== Partial vs Differential Depth Comparison ===")

	client := binanceperp.NewWSClient(nil)

	var partialDepthCount int64
	var diffDepthCount int64

	// Subscribe to partial depth (top 5 levels snapshot)
	partialOptions := &binanceperp.DepthSubscriptionOptions{}
	partialOptions.
		WithConnect(func() {
			fmt.Println("✓ Connected to partial depth stream (top 5 levels)")
		}).
		WithDepth(func(depth binanceperp.WSDepth) {
			count := atomic.AddInt64(&partialDepthCount, 1)
			fmt.Printf("PartialDepth #%d: %d bids, %d asks (always ≤ 5)\n",
				count, len(depth.Bids), len(depth.Asks))
		}).
		WithDisconnect(func() {
			fmt.Println("✓ Disconnected from partial depth stream")
		})

	// Subscribe to differential depth (order book changes)
	diffOptions := &binanceperp.DiffDepthSubscriptionOptions{}
	diffOptions.
		WithConnect(func() {
			fmt.Println("✓ Connected to differential depth stream (order book changes)")
		}).
		WithDiffDepth(func(diffDepth binanceperp.WSDepth) {
			count := atomic.AddInt64(&diffDepthCount, 1)
			fmt.Printf("DiffDepth #%d: %d bid changes, %d ask changes (varies by market activity)\n",
				count, len(diffDepth.Bids), len(diffDepth.Asks))
		}).
		WithDisconnect(func() {
			fmt.Println("✓ Disconnected from differential depth stream")
		})

	// Subscribe to both streams for the same symbol
	unsubscribePartial, err := client.SubscribeDepth("ethusdt", binanceperp.DepthLevel5, binanceperp.DepthUpdate250ms, partialOptions)
	if err != nil {
		fmt.Printf("Failed to subscribe to partial depth: %v\n", err)
		return
	}

	unsubscribeDiff, err := client.SubscribeDiffDepth("ethusdt", binanceperp.DepthUpdate250ms, diffOptions)
	if err != nil {
		fmt.Printf("Failed to subscribe to differential depth: %v\n", err)
		return
	}

	fmt.Println("Comparing both stream types for 6 seconds...")
	time.Sleep(6 * time.Second)

	unsubscribePartial()
	unsubscribeDiff()

	fmt.Printf("Final counts - PartialDepth (snapshots): %d, DiffDepth (changes): %d\n",
		atomic.LoadInt64(&partialDepthCount), atomic.LoadInt64(&diffDepthCount))

	fmt.Println("\n📝 Key Differences:")
	fmt.Println("  • Partial Depth: Fixed number of top levels (5/10/20), snapshot data")
	fmt.Println("  • Differential Depth: Variable number of changes, delta/update data")
	fmt.Println("  • Both use same data structure but different semantic meaning")
}
