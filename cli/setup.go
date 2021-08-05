package cli

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/viper"
)

var DB *pg.DB

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
	return CreateSchema()
}

func ConnectDB() error {
	opt, err := pg.ParseURL(viper.GetString("DATABASE_URL"))
	if err != nil {
		return err
	}
	DB = pg.Connect(opt)
	return nil
}

func CloseDB() error {
	return DB.Close()
}

func CreateSchema() error {
	models := []interface{}{
		(*Client)(nil),
	}
	for _, model := range models {
		err := DB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
