package cli

import (
	"github.com/anupshk/stock/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllClients() []Client {
	var clients []Client
	clientCollection := DB.Collection("clients")
	cursor, err := clientCollection.Find(Ctx, bson.M{})
	//defer cursor.Close(Ctx)
	if err != nil {
		panic(err)
	}
	if err = cursor.All(Ctx, &clients); err != nil {
		panic(err)
	}
	return clients
}

func InsertClient(client *Client) (res *mongo.InsertOneResult, err error) {
	clientCollection := DB.Collection("clients")
	res, err = clientCollection.InsertOne(Ctx, client)
	return
}

func AddUniqueIndex(collection string, indexKeys interface{}) error {
	col := DB.Collection(collection)
	index := mongo.IndexModel{
		Keys:    indexKeys,
		Options: options.Index().SetUnique(true),
	}
	indexName, err := col.Indexes().CreateOne(Ctx, index)
	if err != nil {
		util.ErrorLogger.Println("Error creating index ", indexKeys)
		return err
	}
	util.InfoLogger.Println("Created index ", indexName)
	return nil
}
