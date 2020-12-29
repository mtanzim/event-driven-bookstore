package service

import (
	"context"
	"log"
	"time"

	localDTO "github.com/mtanzim/event-driven-bookstore/warehouse-server/dto"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DLQMonitorService struct {
	paymentDLQCollection *mongo.Collection
	warehouseCollection  *mongo.Collection
}

func NewDLQMonitorService(paymentDLQColl, warehouseColl *mongo.Collection) *DLQMonitorService {
	return &DLQMonitorService{paymentDLQColl, warehouseColl}
}

// TODO: proper usage of context here?
func (s DLQMonitorService) Monitor() {
	for {
		go func() {
			// log.Println("Monitoring Payment DLQ")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			opts := options.Find()
			// opts.SetLimit(10)
			cur, err := s.paymentDLQCollection.Find(ctx, bson.D{}, opts)
			defer cur.Close(ctx)

			if err != nil {
				log.Println(err)
			}
			for cur.Next(ctx) {
				var item localDTO.CartPaymentDLQItem
				curErr := cur.Decode(&item)
				log.Println("In payment DLQ", item)
				if curErr == nil {
					update := bson.D{{"$set", bson.D{{"paid", true}}}}
					filter := bson.M{"_id": item.CartID}
					updateRes, updateErr := s.warehouseCollection.UpdateOne(ctx, filter, update)
					if updateErr != nil {
						log.Println(updateErr)
					}
					if updateRes.ModifiedCount == 1 || updateRes.UpsertedCount == 1 {
						s.paymentDLQCollection.DeleteOne(ctx, bson.M{"cartId": item.CartID})
						log.Println("Successfully updated payment status for cart from DLQ", item.CartID)
					}
				} else {
					log.Println(curErr)
				}

			}
		}()

		time.Sleep(3 * time.Second)
	}
}
