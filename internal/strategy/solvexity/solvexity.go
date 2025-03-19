package solvexity

import (
	"context"
	"errors"
	"time"

	"github.com/BullionBear/sequex/internal/payload"
	"github.com/BullionBear/sequex/internal/strategy"
	pb "github.com/BullionBear/sequex/pkg/protobuf/solvexity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ strategy.Strategy = (*Solvexity)(nil)

type Solvexity struct {
	c pb.SolvexityClient
}

func NewSolvexity(client pb.SolvexityClient) *Solvexity {
	return &Solvexity{
		c: client,
	}
}

func (s *Solvexity) OnKLineUpdate(meta payload.KLineUpdate) error {
	ts := time.Unix(meta.Timestamp, 0).UTC()
	resp, err := s.c.Solve(context.Background(), &pb.SolveRequest{
		Symbol:    meta.Symbol,
		Timestamp: timestamppb.New(ts),
	})
	if err != nil {
		return err
	}
	if resp.GetStatus() != pb.StatusType_SUCCESS {
		return errors.New(resp.GetMessage())
	}
	return nil
}
