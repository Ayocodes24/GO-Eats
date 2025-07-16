package database

import (
	"GO-Eats/pkg/database/models/user"
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"os"
	"strconv"
	"time"
)

// These are all method calls , which can be directly accessed by any or all services
type Database interface {
	Insert(ctx context.Context, model any) (sql.Result, error)
	Delete(ctx context.Context, tableName string, filter Filter) (sql.Result, error)
	Select(ctx context.Context, model any, columnName string, parameter any) error
	SelectAll(ctx context.Context, tableName string, model any) error
	SelectWithRelation(ctx context.Context, model any, relations []string, Condition Filter) error
	SelectWithMultipleFilter(ctx context.Context, model any, Condition Filter) error
	Raw(ctx context.Context, model any, query string, args ...interface{}) error
	Update(ctx context.Context, tableName string, Set Filter, Condition Filter) (sql.Result, error)
	Count(ctx context.Context, tableName string, ColumnExpression string, columnName string, parameter any) (int64, error)
	Migrate() error
	HealthCheck() bool
	Close() error
}

type Filter map[string]any

type DB struct {
	db *bun.DB
} //INITIALIZING A GLOBAL INSTANCE HERE

func New() Database {
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	databasePort, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatal("Invalid DB Port")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUsername, dbPassword, dbHost, databasePort, dbName)
	database := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(database, pgdialect.New())
	return &DB{db: db}

}

func (d *DB) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	models := []interface{}{
		(*user.User)(nil),
		(*restaurant.Restaurant)(nil),
		(*restaurant.MenuItem)(nil),
		(*review.Review)(nil),
		(*order.Order)(nil),
		(*order.OrderItems)(nil),
		(*cart.Cart)(nil),
		(*cart.CartItems)(nil),
		(*delivery.DeliveryPerson)(nil),
		(*delivery.Deliveries)(nil),
	}

	for _, model := range models {
		if _, err := d.db.NewCreateTable().Model(model).WithForeignKeys().IfNotExists().Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}
