package cli

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var DbClient *mongo.Client
var Ctx context.Context
var CancelFunc context.CancelFunc

func Setup() error {
	viper.SetConfigFile(".env")
	configErr := viper.ReadInConfig()
	if configErr != nil {
		return configErr
	}
	dbErr := ConnectDB()
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func ConnectDB() error {
	var err error
	DbClient, err = mongo.NewClient(options.Client().ApplyURI(viper.GetString("DATABASE_URL")))
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

func CloseDB() error {
	return DbClient.Disconnect(Ctx)
}
