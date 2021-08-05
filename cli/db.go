package cli

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
