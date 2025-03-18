package main

import (
	"context"
	"testing"

	pb "github.com/BullionBear/sequex/pkg/protobuf/greet"
	"github.com/stretchr/testify/assert"
)

func TestSayHello(t *testing.T) {
	s := &server{}
	req := &pb.HelloRequest{Name: "test"}
	resp, err := s.SayHello(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Hello test", resp.Message)
}
