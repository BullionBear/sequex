package chronicler

import "go.mongodb.org/mongo-driver/bson"

type History struct {
	openTime int64  `bson:"open_time"`
	data     bson.M `bson:"data"`
	wallet   bson.M `bson:"wallet"`
}

func NewHistory(openTime int64, data bson.M, wallet bson.M) *History {
	return &History{
		openTime: openTime,
		data:     data,
		wallet:   wallet,
	}
}
