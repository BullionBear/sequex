package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

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
	nc *nats.Conn
}

// NewRPCClient creates a new RPC client
func NewRPCClient(nc *nats.Conn) *RPCClient {
	return &RPCClient{nc: nc}
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

	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Make request and wait for response
	resp, err := c.nc.Request("rpc.requests", reqData, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	var response Response
	if err := json.Unmarshal(resp.Data, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

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
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Connected to NATS server")

	// Create RPC client
	client := NewRPCClient(nc)

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Example 1: Addition operation
	log.Println("Making addition RPC call...")
	addResp, err := client.CallAdd(10, 20)
	if err != nil {
		log.Printf("Error calling add: %v", err)
	} else {
		log.Printf("Addition result: %d + %d = %d", 10, 20, addResp.Sum)
	}

	// Example 2: String operation
	log.Println("Making string RPC call...")
	strResp, err := client.CallString("Hello, NATS RPC!")
	if err != nil {
		log.Printf("Error calling string: %v", err)
	} else {
		log.Printf("String result: length=%d, text='%s'", strResp.Length, strResp.Upper)
	}

	// Example 3: Multiple calls
	log.Println("Making multiple RPC calls...")
	for i := 1; i <= 3; i++ {
		addResp, err := client.CallAdd(i*10, i*5)
		if err != nil {
			log.Printf("Error in call %d: %v", i, err)
		} else {
			log.Printf("Call %d result: %d + %d = %d", i, i*10, i*5, addResp.Sum)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Example 4: Error handling - unknown method
	log.Println("Testing error handling with unknown method...")
	resp, err := client.Call("unknown_method", map[string]string{"test": "data"})
	if err != nil {
		log.Printf("Error calling unknown method: %v", err)
	} else if resp.Error != "" {
		log.Printf("Expected error received: %s", resp.Error)
	}

	log.Println("RPC client demo completed")
}
