package chronicler

import "go.mongodb.org/mongo-driver/bson"

type History struct {
	OpenTime int64   `bson:"open_time"`
	Price    float64 `bson:"price"`
	Data     bson.M  `bson:"data"`
	Wallet   bson.M  `bson:"wallet"`
}

func NewHistory(openTime int64, price float64, data bson.M, wallet bson.M) *History {
	return &History{
		OpenTime: openTime,
		Price:    price,
		Data:     data,
		Wallet:   wallet,
	}
}
