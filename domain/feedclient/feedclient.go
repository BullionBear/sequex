package feedclient


import (
	"github.com/BullionBear/crypto-trade/api/gen/feed"
)

struct FeedClient {
	conn *grpc.ClientConn
}

func NewFeedClient(conn *grpc.ClientConn) *FeedClient {
	return &FeedClient{
		conn: conn,
	}
}