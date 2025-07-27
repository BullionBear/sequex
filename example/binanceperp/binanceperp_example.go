package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	/*
		// Example: Create a LIMIT BUY order
		createOrderReq := binanceperp.CreateOrderRequest{
			Symbol:      "BTCUSDT",
			Side:        binanceperp.OrderSideBuy,
			Type:        binanceperp.OrderTypeLimit,
			TimeInForce: binanceperp.TimeInForceGTC,
			Quantity:    "0.001",
			Price:       "25000.00",  // Set a low price so it won't fill immediately
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
	*/
	fmt.Println("Create Order example is commented out - creates real orders!")

	// Example 11: Cancel Order (TRADE - signed request - requires API credentials)
	fmt.Println("\n--- Cancel Order (TRADE - Signed Request) ---")
	/*
		// Cancel the order created above using order ID
		cancelOrderReq := binanceperp.CancelOrderRequest{
			Symbol:  "BTCUSDT",
			OrderId: createResp.Data.OrderId,  // Use the order ID from create order
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
	*/
	fmt.Println("Cancel Order example is commented out - requires existing order!")

	// Example 12: Cancel All Orders (TRADE - signed request - requires API credentials)
	fmt.Println("\n--- Cancel All Orders (TRADE - Signed Request) ---")
	/*
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
	*/
	fmt.Println("Cancel All Orders example is commented out - cancels ALL open orders!")

	fmt.Println("\nTo test trading functions:")
	fmt.Println("1. Set your API credentials in environment variables")
	fmt.Println("2. Uncomment the trading example code")
	fmt.Println("3. Use testnet for safe testing: https://testnet.binancefuture.com")
	fmt.Println("4. Be careful with real trading - orders can lose money!")
	fmt.Println("5. Cancel All Orders will close ALL your open orders for the symbol!")
}
