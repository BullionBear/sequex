package main

import (
	"encoding/json"
	"fmt"
	"log"

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

// UnmarshalParams unmarshals request parameters into the target struct
func UnmarshalParams(params interface{}, target interface{}) error {
	paramsData, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("error marshaling params: %v", err)
	}

	return json.Unmarshal(paramsData, target)
}

// RPCHandler handles RPC requests
type RPCHandler struct {
	nc *nats.Conn
}

// NewRPCHandler creates a new RPC handler
func NewRPCHandler(nc *nats.Conn) *RPCHandler {
	return &RPCHandler{nc: nc}
}

// handleAdd handles addition requests
func (h *RPCHandler) handleAdd(req *Request) (*Response, error) {
	var addReq AddRequest
	if err := UnmarshalParams(req.Params, &addReq); err != nil {
		return &Response{
			ID:    req.ID,
			Error: "Invalid request parameters",
		}, nil
	}

	result := AddResponse{
		Sum: addReq.A + addReq.B,
	}

	return &Response{
		ID:     req.ID,
		Result: result,
	}, nil
}

// handleString handles string operations
func (h *RPCHandler) handleString(req *Request) (*Response, error) {
	var strReq StringRequest
	if err := UnmarshalParams(req.Params, &strReq); err != nil {
		return &Response{
			ID:    req.ID,
			Error: "Invalid request parameters",
		}, nil
	}

	result := StringResponse{
		Length: len(strReq.Text),
		Upper:  fmt.Sprintf("%s (length: %d)", strReq.Text, len(strReq.Text)),
	}

	return &Response{
		ID:     req.ID,
		Result: result,
	}, nil
}

// handleRequest processes incoming RPC requests
func (h *RPCHandler) handleRequest(msg *nats.Msg) {
	var req Request
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		log.Printf("Error unmarshaling request: %v", err)
		return
	}

	log.Printf("Received request: %s - %s", req.ID, req.Method)

	var resp *Response
	var err error

	switch req.Method {
	case "add":
		resp, err = h.handleAdd(&req)
	case "string":
		resp, err = h.handleString(&req)
	default:
		resp = &Response{
			ID:    req.ID,
			Error: "Unknown method: " + req.Method,
		}
	}

	if err != nil {
		log.Printf("Error handling request: %v", err)
		resp = &Response{
			ID:    req.ID,
			Error: "Internal server error",
		}
	}

	// Send response
	respData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	if err := h.nc.Publish(msg.Reply, respData); err != nil {
		log.Printf("Error publishing response: %v", err)
	}
}

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Connected to NATS server")

	// Create RPC handler
	handler := NewRPCHandler(nc)

	// Subscribe to RPC requests
	sub, err := nc.Subscribe("rpc.requests", handler.handleRequest)
	if err != nil {
		log.Fatalf("Error subscribing to RPC requests: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("RPC Server started. Listening for requests on 'rpc.requests'")
	log.Println("Available methods: add, string")

	// Keep the server running
	select {}
}
