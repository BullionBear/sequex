package feedclient

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/BullionBear/crypto-trade/api/gen/feed"
	"github.com/BullionBear/crypto-trade/domain/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FeedClient struct {
	conn   *grpc.ClientConn
	client *feed.FeedClient
}

func NewFeedClient(host string, port int) *FeedClient {
	srvAddr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClient(srvAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	feedClient := feed.NewFeedClient(conn)
	return &FeedClient{
		conn:   conn,
		client: &feedClient,
	}
}

func (fc *FeedClient) Close() {
	fc.conn.Close()
}

func (fc *FeedClient) SubscribeKlines(handler func(event *models.Kline)) error {
	stream, err := (*fc.client).SubscribeKline(context.Background(), &emptypb.Empty{})
	if err != nil {
		logrus.Errorf("could not subscribe to kline: %v", status.Convert(err).Message())
		return err
	}
	for {
		pbkline, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				logrus.Infoln("Stream closed by server")
				return nil
			} else {
				logrus.Errorf("Error receiving from kline stream: %v", status.Convert(err).Message())
			}
		}
		kline := models.NewKlineFromPb(pbkline.Kline)
		handler(kline)
	}
}

func (fc *FeedClient) LoadHistoricalKlines(handler func(event *models.Kline), start, end int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := (*fc.client).ReadHistoricalKline(ctx, &feed.ReadKlineRequest{Start: start, End: end})
	if err != nil {
		logrus.Errorf("could not read historical kline: %v", status.Convert(err).Message())
		return err
	}
	for {
		pbkline, err := stream.Recv()
		if err != nil {
			logrus.Errorf("Error receiving from historical kline stream: %v", status.Convert(err).Message())
			break
		}
		kline := models.NewKlineFromPb(pbkline.Kline)
		handler(kline)
	}
	return nil
}
