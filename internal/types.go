package internal

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/oschwald/geoip2-golang"
	"log"
)

type key int

var dbKey key = 0
var geoDBKey key = 1
var feedManagerKey key = 2

func NewDBToContext(ctx context.Context, dbDsn string) context.Context {
	db, err := InitDatabase(dbDsn)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	return context.WithValue(ctx, dbKey, db)
}

func DBFromContext(ctx context.Context) (*gorm.DB, bool) {
	db, ok := ctx.Value(dbKey).(*gorm.DB)
	return db, ok
}

func NewGeoDBToContext(ctx context.Context, geoDBDsn string) context.Context {
	db, err := geoip2.Open(geoDBDsn)
	if err != nil {
		log.Fatal(err)
	}
	//TODO cleanu defer db.Close()
	return context.WithValue(ctx, geoDBKey, db)
}

func GeoDBFromContext(ctx context.Context) (*geoip2.Reader, bool) {
	db, ok := ctx.Value(geoDBKey).(*geoip2.Reader)
	return db, ok
}

func NewFeedManagerToContext(ctx context.Context) context.Context {
	fm := NewFeedManager(ctx)
	return context.WithValue(ctx, feedManagerKey, fm)
}

func FeedManagerFromContext(ctx context.Context) (*FeedManager, bool) {
	fm, ok := ctx.Value(feedManagerKey).(*FeedManager)
	return fm, ok
}
