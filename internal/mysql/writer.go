package mysql

import (
	"context"
	"fmt"
	"log"
	"myMod/internal/clickhouseStat"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type writer struct {
	db        driver.Conn
	tableName string
}

func NewClickhouseWriter(host string, port uint16, database, table, user, password string) (*writer, error) {

	DatabaseMigration(host, port, database, table, user, password)
	TableMigration(host, port, database, table, user, password)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return nil, err
	}

	return &writer{
		db:        conn,
		tableName: table,
	}, nil
}

func DatabaseMigration(host string, port uint16, database, table, user, password string) error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: user,
			Password: password,
		},
		Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	OldQuery := `CREATE DATABASE IF NOT EXISTS rotator;`
	ctx := context.TODO()
	err = conn.Exec(ctx, OldQuery)

	if err != nil {
		log.Fatalf("Failed to create migrations table: %s", err)
	} else {
		fmt.Println("Migrations table created successfully or already exists.")
	}

	return nil
}

func TableMigration(host string, port uint16, database, table, user, password string) error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	NewQuery := `CREATE TABLE IF NOT EXISTS statistics (
			timestamp    String,
			country      String,
			os   		 String,
			browser   	 String,
			requests   	 Int64,
			impressions  Int64) ENGINE = MergeTree() ORDER BY timestamp;`
	ctx := context.TODO()
	err = conn.Exec(ctx, NewQuery)

	if err != nil {
		log.Fatalf("Failed to create migrations table: %s", err)
	} else {
		fmt.Println("Migrations table created successfully or already exists.")
	}

	return nil
}

func (w *writer) Insert(rows clickhouseStat.Rows) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	batch, err := w.db.PrepareBatch(ctx, fmt.Sprintf(insertQuery, w.tableName))
	if err != nil {
		return err
	}

	for k, v := range rows {
		err := batch.Append(
			time.Unix(k.Timestamp, 0),
			k.Country,
			k.Os,
			k.Browser,
			v.Requests,
			v.Impressions,
		)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return batch.Send()
}

var insertQuery = `INSERT INTO %s (timestamp, country, os, browser, requests, impressions)`
