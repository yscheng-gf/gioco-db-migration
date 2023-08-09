package helper

import (
	"context"

	"github.com/yscheng-gf/gioco-db-migration/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FetchAllOperator(mongo *mongo.Client) (*[]models.Operator, error) {
	cursor, err := mongo.Database("GIOCO_PLUS").Collection("sys_operators").Find(context.TODO(), bson.M{"status": "1"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	ops := new([]models.Operator)
	err = cursor.All(context.TODO(), ops)
	if err != nil {
		return nil, err
	}
	return ops, nil
}
