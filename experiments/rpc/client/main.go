package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/nats-io/nats.go"
)

// Request represents a generic RPC request
type Request struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// Response represents a generic RPC response
type Response struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// AddRequest represents an addition operation request
type AddRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// AddResponse represents an addition operation response
type AddResponse struct {
	Sum int `json:"sum"`
}

// StringRequest represents a string operation request
type StringRequest struct {
	Text string `json:"text"`
}

// StringResponse represents a string operation response
type StringResponse struct {
	Length int    `json:"length"`
	Upper  string `json:"upper"`
}

// RPCClient represents an RPC client
type RPCClient struct {
	nc     *nats.Conn
	logger log.Logger
}

// NewRPCClient creates a new RPC client
func NewRPCClient(nc *nats.Conn, logger log.Logger) *RPCClient {
	return &RPCClient{
		nc:     nc,
		logger: logger.With(log.String("component", "rpc_client")),
	}
}

// generateID generates a unique request ID
func (c *RPCClient) generateID() string {
	return fmt.Sprintf("req_%d", rand.Int63())
}

// Call makes an RPC call and waits for the response
func (c *RPCClient) Call(method string, params interface{}) (*Response, error) {
	req := &Request{
		ID:     c.generateID(),
		Method: method,
		Params: params,
	}

	c.logger.Debug("Making RPC call",
		log.String("request_id", req.ID),
		log.String("method", method),
	)

	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Make request and wait for response
	resp, err := c.nc.Request("rpc.requests", reqData, 5*time.Second)
	if err != nil {
		c.logger.Error("Error making RPC request",
			log.String("request_id", req.ID),
			log.String("method", method),
			log.Error(err),
		)
		return nil, fmt.Errorf("error making request: %v", err)
	}

	var response Response
	if err := json.Unmarshal(resp.Data, &response); err != nil {
		c.logger.Error("Error unmarshaling RPC response",
			log.String("request_id", req.ID),
			log.Error(err),
		)
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	c.logger.Debug("RPC response received",
		log.String("request_id", req.ID),
		log.String("method", method),
		log.Bool("has_error", response.Error != ""),
	)

	return &response, nil
}

// CallAdd makes an addition RPC call
func (c *RPCClient) CallAdd(a, b int) (*AddResponse, error) {
	params := AddRequest{A: a, B: b}
	resp, err := c.Call("add", params)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("RPC error: %s", resp.Error)
	}

	// Convert result to AddResponse
	resultData, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("error marshaling result: %v", err)
	}

	var addResp AddResponse
	if err := json.Unmarshal(resultData, &addResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling add response: %v", err)
	}

	return &addResp, nil
}

// CallString makes a string operation RPC call
func (c *RPCClient) CallString(text string) (*StringResponse, error) {
	params := StringRequest{Text: text}
	resp, err := c.Call("string", params)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("RPC error: %s", resp.Error)
	}

	// Convert result to StringResponse
	resultData, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("error marshaling result: %v", err)
	}

	var strResp StringResponse
	if err := json.Unmarshal(resultData, &strResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling string response: %v", err)
	}

	return &strResp, nil
}

func main() {
	// Initialize structured logger
	logger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithTimeRotation("./logs", "rpc_client.log", 24*time.Hour, 7),
	)

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Fatal("Error connecting to NATS",
			log.String("nats_url", nats.DefaultURL),
			log.Error(err),
		)
	}
	defer nc.Close()

	logger.Info("Connected to NATS server",
		log.String("nats_url", nats.DefaultURL),
	)

	// Create RPC client
	client := NewRPCClient(nc, logger)

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Example 1: Addition operation
	logger.Info("Making addition RPC call")
	addResp, err := client.CallAdd(10, 20)
	if err != nil {
		logger.Error("Error calling add",
			log.Int("operand_a", 10),
			log.Int("operand_b", 20),
			log.Error(err),
		)
	} else {
		logger.Info("Addition operation completed",
			log.Int("operand_a", 10),
			log.Int("operand_b", 20),
			log.Int("result", addResp.Sum),
		)
	}

	// Example 2: String operation
	logger.Info("Making string RPC call")
	strResp, err := client.CallString("Hello, NATS RPC!")
	if err != nil {
		logger.Error("Error calling string",
			log.String("input_text", "Hello, NATS RPC!"),
			log.Error(err),
		)
	} else {
		logger.Info("String operation completed",
			log.String("input_text", "Hello, NATS RPC!"),
			log.Int("text_length", strResp.Length),
			log.String("result", strResp.Upper),
		)
	}

	// Example 3: Multiple calls
	logger.Info("Making multiple RPC calls")
	for i := 1; i <= 3; i++ {
		addResp, err := client.CallAdd(i*10, i*5)
		if err != nil {
			logger.Error("Error in RPC call",
				log.Int("call_number", i),
				log.Int("operand_a", i*10),
				log.Int("operand_b", i*5),
				log.Error(err),
			)
		} else {
			logger.Info("RPC call completed",
				log.Int("call_number", i),
				log.Int("operand_a", i*10),
				log.Int("operand_b", i*5),
				log.Int("result", addResp.Sum),
			)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Example 4: Error handling - unknown method
	logger.Info("Testing error handling with unknown method")
	resp, err := client.Call("unknown_method", map[string]string{"test": "data"})
	if err != nil {
		logger.Error("Error calling unknown method",
			log.String("method", "unknown_method"),
			log.Error(err),
		)
	} else if resp.Error != "" {
		logger.Warn("Expected error received from unknown method",
			log.String("method", "unknown_method"),
			log.String("error", resp.Error),
		)
	}

	logger.Info("RPC client demo completed")
}
