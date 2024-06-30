package reporter

import "go.mongodb.org/mongo-driver/mongo"

type Reporter struct {
	mongo *mongo.Client
}

func NewReporter(mongo *mongo.Client) *Reporter {
	return &Reporter{
		mongo: mongo,
	}
}

func (r *Reporter) Record(key string, value interface{}) error {
	return nil
}
