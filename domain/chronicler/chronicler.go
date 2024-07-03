package chronicler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chronicler struct {
	mongo *mongo.Client

	historyC       chan *History
	historicalData []*History
	mutex          sync.Mutex
}

func NewChronicler(mongo *mongo.Client, name string) *Chronicler {
	chronicler := &Chronicler{
		mongo: mongo,

		historyC:       make(chan *History, 86400),
		historicalData: make([]*History, 0),
		mutex:          sync.Mutex{},
	}
	// Create a page with the given name
	now := time.Now()
	layout := "20060102150405"
	formattedTime := now.Format(layout)
	pageName := fmt.Sprintf("%s.%s", name, formattedTime)
	err := chronicler.createPage(pageName)
	if err != nil {
		logrus.Errorf("failed to create page: %v", err)
	}
	go chronicler.batchWriter(pageName, 3600) // 1 hour
	return chronicler
}

func (r *Chronicler) createPage(page string) error {
	// Create a new collection with the given name
	ctx := context.TODO()
	database := r.mongo.Database(CollectionHistory)

	// Check if the collection already exists
	collections, err := database.ListCollectionNames(ctx, bson.M{"name": page})
	if err != nil {
		return err
	}

	if len(collections) > 0 {
		return fmt.Errorf("collection %s already exists", page)
	}

	// Create the collection
	err = database.CreateCollection(ctx, page, options.CreateCollection())
	if err != nil {
		return err
	}

	// Create an index on the "open_time" field
	collection := database.Collection(page)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"open_time": 1}, // Index on the "open_time" field in ascending order
		Options: options.Index().SetName("open_time_index"),
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index on open_time: %v", err)
	}

	return nil
}

func (r *Chronicler) Record(history *History) {
	r.historyC <- history
}

func (r *Chronicler) batchWriter(page string, batchSize int) {
	database := r.mongo.Database(CollectionHistory)
	collection := database.Collection(page)

	for history := range r.historyC {
		r.historicalData = append(r.historicalData, history)
		if len(r.historicalData) >= batchSize {
			r.flushData(collection)
		}
	}
	if len(r.historicalData) >= batchSize {
		r.flushData(collection)
	}
}

func (r *Chronicler) flushData(collection *mongo.Collection) {
	data := make([]interface{}, len(r.historicalData))
	for i, v := range r.historicalData {
		data[i] = v
	}

	if _, err := collection.InsertMany(context.TODO(), data); err != nil {
		logrus.Errorf("failed to insert data: %v", err)
	}
	r.historicalData = r.historicalData[:0]
}

func (r *Chronicler) Close() {
	close(r.historyC)
}
