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

func GetClient(ident string) (client Client, err error) {
	client = Client{}
	clientCollection := DB.Collection("clients")
	err = clientCollection.FindOne(Ctx, bson.M{"ident": ident}).Decode(&client)
	return
}

func InsertShares(shares []Share) (res *mongo.InsertManyResult, err error) {
	list := make([]interface{}, len(shares))
	for i, v := range shares {
		list[i] = v
	}
	shareCollection := DB.Collection("shares")
	res, err = shareCollection.InsertMany(Ctx, list)
	return
}

func (client Client) GetShares() (shares []Share, err error) {
	shareCollection := DB.Collection("shares")
	filterCursor, err := shareCollection.Find(Ctx, bson.M{"client": client.Id})
	if err != nil {
		return
	}
	if err = filterCursor.All(Ctx, &shares); err != nil {
		return
	}
	return
}
