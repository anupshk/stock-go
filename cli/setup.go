package cli

import (
	"context"
	"time"

	"github.com/anupshk/stock/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var DbClient *mongo.Client
var Ctx context.Context
var CancelFunc context.CancelFunc

func Setup() error {
	dbErr := ConnectDB()
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func ConnectDB() error {
	var err error
	DbClient, err = mongo.NewClient(options.Client().ApplyURI(util.DATABASE_URL))
	if err != nil {
		return err
	}
	Ctx, CancelFunc = context.WithTimeout(context.Background(), 10*time.Second)
	err = DbClient.Connect(Ctx)
	if err != nil {
		return err
	}
	err = DbClient.Ping(Ctx, nil)
	if err != nil {
		return err
	}
	DB = DbClient.Database("stock")
	return nil
}

func createIndexes() error {
	cErr := AddUniqueIndex("clients", bson.M{"ident": 1})
	return cErr
}

func CloseDB() error {
	return DbClient.Disconnect(Ctx)
}
