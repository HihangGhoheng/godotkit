package gdk_mongo

import (
	"context"
	gdk_helpers "github.com/HihangGhoheng/godotkit/helpers"
	gdk_types "github.com/HihangGhoheng/godotkit/types"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoDb struct {
	Database *mongo.Database
}

func (c *MongoDb) ShutdownConnection(ctx context.Context, logger *logrus.Logger) {
	if err := c.Database.Client().Disconnect(ctx); err != nil {
		gdk_helpers.FatalOnError(err, "Failed to shutdown connection of mongodb")
	} else {
		logger.Infof("Successfully close mongodb connection!")
	}
}

func Connect(ctx context.Context, opt gdk_types.MongoConfig, log *logrus.Logger) *MongoDb {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client()

	switch opt.ReadPreference {
	case "primary":
		clientOptions.SetReadPreference(readpref.Primary())
	case "preimaryPreferred":
		clientOptions.SetReadPreference(readpref.PrimaryPreferred())
	case "secondary":
		clientOptions.SetReadPreference(readpref.Secondary())
	case "secondaryPreferred":
		clientOptions.SetReadPreference(readpref.SecondaryPreferred())
	default:
		clientOptions.SetReadPreference(readpref.Primary())
		log.Warnf("Unknown readPreferrence! We'll set to primary")
	}

	clientOptions.ApplyURI(opt.Dsn).SetServerAPIOptions(serverApi)
	clientOptions.SetMinPoolSize(uint64(opt.MinPoolSize))
	clientOptions.SetMaxPoolSize(uint64(opt.MaxPoolSize))
	clientOptions.SetMaxConnIdleTime(time.Duration(uint(opt.MaxConnectionIdleTime) * uint(time.Millisecond)))

	monitor := event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			log.Info(e.Command.String())
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			log.Infof(
				"Command: %s | Reply: %s | Duration: %s",
				e.CommandName,
				e.Reply,
				e.Duration.String(),
			)
		},
		Failed: func(_ context.Context, e *event.CommandFailedEvent) {
			log.Errorf(
				"Command: %s | Failed: %s | Duration: %s",
				e.CommandName,
				e.Failure,
				e.Duration.String(),
			)
		},
	}

	clientOptions.SetMonitor(&monitor)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect mongdb: %+v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to connect mongdb: %+v", err)
	}

	log.Info("Success connecting to mongodb")

	return &MongoDb{
		Database: client.Database(opt.DatabaseName),
	}
}
