package chronicler

import "go.mongodb.org/mongo-driver/bson"

type History struct {
	OpenTime int64  `bson:"open_time"`
	Data     bson.M `bson:"data"`
	Wallet   bson.M `bson:"wallet"`
}

func NewHistory(openTime int64, data bson.M, wallet bson.M) *History {
	return &History{
		OpenTime: openTime,
		Data:     data,
		Wallet:   wallet,
	}
}
