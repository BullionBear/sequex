package feedclient

import (
	"context"
	"log"

	"github.com/BullionBear/crypto-trade/api/gen/feed"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type FeedClient struct {
	conn *grpc.ClientConn
}

func NewFeedClient(conn *grpc.ClientConn) *FeedClient {
	return &FeedClient{
		conn: conn,
	}
}

func (fc *FeedClient) GetConfig(ctx context.Context) error {
	config, err := fc.GetConfig(ctx, &feed.Empty{})
	if err != nil {
		logrus.Errorf("could not get status: %v", status.Convert(err).Message())
		return err
	}
	log.Printf("Status: %v", config)
	return nil
}

// func (f *FeedClient) SubscribeKlines(handler func(event *Kline)) error {
// 	stream, err := f.conn.SubscribeKline(context.Background(), &emptypb.Empty{})
// 	if err != nil {
// 		logrus.Errorf("could not subscribe to kline: %v", status.Convert(err).Message())
// 		return err
// 	}
// 	for {
// 		kline, err := stream.Recv()
// 		if err != nil {
// 			if err == io.EOF {
// 				logrus.Infoln("Stream closed by server")
// 				return nil
// 			} else {
// 				logrus.Errorf("Error receiving from kline stream: %v", status.Convert(err).Message())
// 			}
// 		}
// 		logrus.Infof("Received kline: %+v", kline)
// 	}
// }
//
