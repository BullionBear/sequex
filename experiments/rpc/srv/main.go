package main

import (
	"encoding/json"
	"fmt"
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
	nc     *nats.Conn
	logger log.Logger
}

// NewRPCHandler creates a new RPC handler
func NewRPCHandler(nc *nats.Conn, logger log.Logger) *RPCHandler {
	return &RPCHandler{
		nc:     nc,
		logger: logger.With(log.String("component", "rpc_handler")),
	}
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

	h.logger.Debug("Addition operation completed",
		log.String("request_id", req.ID),
		log.Int("operand_a", addReq.A),
		log.Int("operand_b", addReq.B),
		log.Int("result", result.Sum),
	)

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

	h.logger.Debug("String operation completed",
		log.String("request_id", req.ID),
		log.String("input_text", strReq.Text),
		log.Int("text_length", result.Length),
		log.String("result", result.Upper),
	)

	return &Response{
		ID:     req.ID,
		Result: result,
	}, nil
}

// handleRequest processes incoming RPC requests
func (h *RPCHandler) handleRequest(msg *nats.Msg) {
	var req Request
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		h.logger.Error("Error unmarshaling request",
			log.Error(err),
			log.String("component", "rpc_handler"),
		)
		return
	}

	h.logger.Info("Received RPC request",
		log.String("request_id", req.ID),
		log.String("method", req.Method),
		log.String("reply_subject", msg.Reply),
	)

	var resp *Response
	var err error

	switch req.Method {
	case "add":
		resp, err = h.handleAdd(&req)
	case "string":
		resp, err = h.handleString(&req)
	default:
		h.logger.Warn("Unknown method requested",
			log.String("request_id", req.ID),
			log.String("method", req.Method),
		)
		resp = &Response{
			ID:    req.ID,
			Error: "Unknown method: " + req.Method,
		}
	}

	if err != nil {
		h.logger.Error("Error handling request",
			log.String("request_id", req.ID),
			log.String("method", req.Method),
			log.Error(err),
		)
		resp = &Response{
			ID:    req.ID,
			Error: "Internal server error",
		}
	}

	// Send response
	respData, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("Error marshaling response",
			log.String("request_id", req.ID),
			log.Error(err),
		)
		return
	}

	if err := h.nc.Publish(msg.Reply, respData); err != nil {
		h.logger.Error("Error publishing response",
			log.String("request_id", req.ID),
			log.String("reply_subject", msg.Reply),
			log.Error(err),
		)
		return
	}

	h.logger.Info("RPC response sent successfully",
		log.String("request_id", req.ID),
		log.String("method", req.Method),
		log.Bool("has_error", resp.Error != ""),
	)
}

func main() {
	// Initialize structured logger
	logger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithTimeRotation("./logs", "rpc_server.log", 24*time.Hour, 7),
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

	// Create RPC handler
	handler := NewRPCHandler(nc, logger)

	// Subscribe to RPC requests
	sub, err := nc.Subscribe("rpc.requests", handler.handleRequest)
	if err != nil {
		logger.Fatal("Error subscribing to RPC requests",
			log.String("subject", "rpc.requests"),
			log.Error(err),
		)
	}
	defer sub.Unsubscribe()

	logger.Info("RPC Server started",
		log.String("subject", "rpc.requests"),
		log.String("available_methods", "add, string"),
	)

	// Keep the server running
	select {}
}
