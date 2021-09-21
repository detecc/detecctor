package database

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"github.com/detecc/detecctor/config"
)

// InitDatabase check if database exists, if not, create a database.
func InitDatabase() {
	CreateStatisticsIfNotExists()
}

// Connect to the MongoDb instance specified in the configuration.
func Connect() {
	config := config.GetServerConfiguration()
	credentials := config.Mongo

	// connect to mongodb
	mongoDbConnection := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		credentials.Username,
		credentials.Password,
		credentials.Host,
		credentials.Port,
	)

	log.Println("Connecting to MongoDB at", credentials.Host, credentials.Port)
	dbOptions := options.Client().ApplyURI(mongoDbConnection)

	err := mgm.SetDefaultConfig(nil, credentials.Database, dbOptions)
	if err != nil {
		log.Fatal(err)
	}

	InitDatabase()
}
