package mysql

import (
	"context"
	"fmt"
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
			return err
		}
	}

	return batch.Send()
}

var insertQuery = `INSERT INTO %s (ts, country, os, browser, campaign_id, requests, impressions)`
