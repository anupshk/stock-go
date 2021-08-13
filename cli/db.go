package cli

import (
	"fmt"
	"strings"
	"time"

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

func (client Client) GetLatestShares() (shares []Share, err error) {
	latest, err := client.GetLastImportDate()
	var parsedDate string
	if err != nil {
		return
	}
	latest_date := latest["last_date"]
	if latest_date != nil {
		parsedDate = fmt.Sprintf("%v", latest_date)
	} else {
		parsedDate = util.GetCurrentTime().String()
	}

	shareCollection := DB.Collection("shares")

	filterCursor, err := shareCollection.Find(Ctx, bson.M{
		"client":        client.Id,
		"imported_date": parsedDate,
	})

	if err != nil {
		return
	}
	if err = filterCursor.All(Ctx, &shares); err != nil {
		return
	}
	return
}

func (client Client) GetShares(date string) (shares []Share, err error) {
	var d time.Time
	d, err = time.Parse(util.DATE_FORMAT, date)
	if err != nil {
		return
	}
	parsedDate := d.Format(util.DATE_FORMAT)

	shareCollection := DB.Collection("shares")

	pipeline := bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key:   "client",
			Value: client.Id,
		}, {
			Key: "$expr",
			Value: bson.M{
				"$eq": bson.A{
					parsedDate,
					bson.D{
						{
							Key: "$dateToString",
							Value: bson.D{
								{Key: "date", Value: "$imported_at"},
								{Key: "format", Value: "%Y-%m-%d"},
							},
						},
					},
				},
			},
		}},
	}}

	// q, _ := json.Marshal(pipeline)
	// fmt.Printf("q %s", q)

	filterCursor, err := shareCollection.Aggregate(Ctx, mongo.Pipeline{pipeline})

	if err != nil {
		return
	}
	if err = filterCursor.All(Ctx, &shares); err != nil {
		return
	}
	return
}

func (client Client) GetLastImportDate() (latest map[string]interface{}, err error) {
	shareCollection := DB.Collection("shares")
	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key:   "client",
			Value: client.Id,
		}},
	}}
	groupStage := bson.D{{
		Key: "$group",
		Value: bson.D{{
			Key:   "_id",
			Value: "$client",
		}, {
			Key: "last_date",
			Value: bson.D{{
				Key:   "$max",
				Value: "$imported_date",
			}},
		}},
	}}
	cursor, err := shareCollection.Aggregate(Ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return
	}
	var res []bson.M
	if err = cursor.All(Ctx, &res); err != nil {
		return
	}
	if len(res) > 0 {
		latest = res[0]
	}
	return
}

func (client Client) GetStockValueSummary(scrip string) (res []ShareValue, err error) {
	shareCollection := DB.Collection("shares")

	matchStage := bson.D{{}}
	groupStage := bson.D{{}}

	if scrip == "" {
		matchStage = bson.D{{
			Key: "$match",
			Value: bson.D{{
				Key:   "client",
				Value: client.Id,
			}},
		}}
		groupStage = bson.D{{
			Key: "$group",
			Value: bson.D{{
				Key:   "_id",
				Value: "$imported_at",
			}, {
				Key: "total",
				Value: bson.D{{
					Key: "$sum",
					Value: bson.D{{
						Key:   "$multiply",
						Value: bson.A{"$last_tran_price", "$balance"},
					}},
				}},
			}},
		}}
	} else {
		scrip = strings.ToUpper(scrip)
		matchStage = bson.D{{
			Key: "$match",
			Value: bson.D{{
				Key: "$and",
				Value: bson.A{
					bson.D{{
						Key:   "client",
						Value: client.Id,
					}},
				},
			}},
		}}
		groupStage = bson.D{{
			Key: "$group",
			Value: bson.D{{
				Key:   "_id",
				Value: "$imported_at",
			}, {
				Key: "total",
				Value: bson.D{{
					Key: "$sum",
					Value: bson.D{{
						Key:   "$multiply",
						Value: bson.A{"$last_tran_price", "$balance"},
					}},
				}},
			}},
		}}
	}

	sortStage := bson.D{{
		Key: "$sort",
		Value: bson.D{{
			Key:   "_id",
			Value: 1,
		}},
	}}
	// q, _ := json.Marshal(groupStage)
	// fmt.Printf("q %s", q)
	cursor, err := shareCollection.Aggregate(Ctx, mongo.Pipeline{matchStage, groupStage, sortStage})
	if err != nil {
		return
	}
	if err = cursor.All(Ctx, &res); err != nil {
		return
	}
	return
}
