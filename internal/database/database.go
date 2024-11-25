package database

import (
	"context"
	"errors"
	"fmt"
	"semaright/internal/entities"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DB struct {
	connectionString string
	conn             *mongo.Client
}

func New(connectionString string) (*DB, error) {
	return &DB{
		connectionString: connectionString,
	}, nil
}

func (d *DB) Connect() error {
	client, err := mongo.Connect(options.Client().ApplyURI(d.connectionString))
	if err != nil {
		return err
	}

	d.conn = client

	return nil
}

func (d *DB) SaveTransaction(ctx context.Context, transaction *entities.Transaction) error {
	session, err := d.conn.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(ctx context.Context) error {
		_, err = d.conn.Database("semaright").Collection("transactions").InsertOne(ctx, transaction)
		if err != nil {
			return err
		}

		todaysDate := time.Now()
		month := fmt.Sprintf("%02d", todaysDate.Month())
		year := todaysDate.Year()

		_, err := d.conn.Database("semaright").Collection("usecase_totals").UpdateOne(ctx, map[string]interface{}{
			"_id": fmt.Sprintf("%s_%s_%d", transaction.UseCaseId, month, year),
		}, map[string]interface{}{
			"$inc": map[string]interface{}{
				"spend": transaction.Spend,
			},
		}, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}

		err = session.CommitTransaction(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetCurrentSpend(ctx context.Context, useCaseId string) (float64, error) {
	todaysDate := time.Now()
	month := fmt.Sprintf("%02d", todaysDate.Month())
	year := todaysDate.Year()

	record := d.conn.Database("semaright").Collection("usecase_totals").FindOne(ctx, map[string]interface{}{
		"_id": fmt.Sprintf("%s_%s_%d", useCaseId, month, year),
	})

	if record.Err() != nil {
		return 0, ErrNoEntryFound
	}

	t := entities.Total{}
	err := record.Decode(&t)
	if err != nil {
		return 0, err
	}

	return t.Spend, nil
}

var ErrNoEntryFound = errors.New("no entry found")
