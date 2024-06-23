package feedclient

import (
	"context"
	"io"
	"log"

	"github.com/BullionBear/crypto-trade/api/gen/feed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FeedClient struct {
	conn *grpc.ClientConn
}

func NewFeedClient(conn *grpc.ClientConn) *FeedClient {
	return &FeedClient{
		conn: conn,
	}
}

func (f *FeedClient) GetConfig() (*feed.Config, error) {
	client := feed.NewFeedClient(f.conn)
	return client.GetConfig(context.Background(), &feed.Empty{})
}

func (f *FeedClient) SubscribeKlines(handler func(event *Kline)) error {
	stream, err := f.SubscribeKline(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Errorf("could not subscribe to kline: %v", status.Convert(err).Message())
		return err
	}
	for {
		kline, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Infoln("Stream closed by server")
				return nil
			} else {
				log.Errorf("Error receiving from kline stream: %v", status.Convert(err).Message())
			}
		}
		log.Infof("Received kline: %+v", kline)
	}
}
