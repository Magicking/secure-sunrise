package internal

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/oschwald/geoip2-golang"
	"log"
	"time"
)

type key int

var ctxValKey key = 0
var dbKey key = 1
var geoDBKey key = 2
var feedManagerKey key = 3
var schedulerKey key = 4

type Time time.Duration

func (t Time) Hours() int64 {
	return int64(time.Duration(t) / time.Hour)
}

func (t Time) Minutes() int64 {
	return int64((time.Duration(t) % time.Hour) / time.Minute)
}

func (t Time) Seconds() int64 {
	return int64((time.Duration(t) % time.Minute) / time.Second)
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hours(), t.Minutes(), t.Seconds())
}

type ctxValues struct {
	m map[key]interface{}
}

func (v ctxValues) Set(k key, val interface{}) {
	v.m[k] = val
}

func (v ctxValues) Get(k key) interface{} {
	val, ok := v.m[k]
	if !ok {
		log.Fatalf("Could not find key: %v", k)
	}
	return val
}

func NewCtxValues() *ctxValues {
	mm := make(map[key]interface{})
	cv := &ctxValues{
		m: mm,
	}
	return cv
}

func getContextValue(ctx context.Context, k key) interface{} {
	v, ok := ctx.Value(ctxValKey).(*ctxValues)
	if !ok {
		log.Fatalf("Could not obtain map context values, key: %v", k)
	}
	return v.Get(k)
}

func setContextValue(ctx context.Context, k key, val interface{}) {
	v, ok := ctx.Value(ctxValKey).(*ctxValues)
	if !ok {
		log.Fatalf("Could not obtain map context values")
	}
	v.Set(k, val)
}

func InitContext(ctx context.Context) context.Context {
	values := NewCtxValues()
	return context.WithValue(ctx, ctxValKey, values)
}

func NewDBToContext(ctx context.Context, dbDsn string) {
	db, err := InitDatabase(dbDsn)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	setContextValue(ctx, dbKey, db)
}

func DBFromContext(ctx context.Context) (*gorm.DB, bool) {
	db, ok := getContextValue(ctx, dbKey).(*gorm.DB)
	return db, ok
}

func NewGeoDBToContext(ctx context.Context, geoDBDsn string) {
	db, err := geoip2.Open(geoDBDsn)
	if err != nil {
		log.Fatal(err)
	}
	//TODO cleanu defer db.Close()
	setContextValue(ctx, geoDBKey, db)
}

func GeoDBFromContext(ctx context.Context) (*geoip2.Reader, bool) {
	db, ok := getContextValue(ctx, geoDBKey).(*geoip2.Reader)
	return db, ok
}

func NewFeedManagerToContext(ctx context.Context) {
	fm := NewFeedManager(ctx)
	go fm.Run(ctx) // TODO death pill
	setContextValue(ctx, feedManagerKey, fm)
}

func FeedManagerFromContext(ctx context.Context) (*FeedManager, bool) {
	fm, ok := getContextValue(ctx, feedManagerKey).(*FeedManager)
	return fm, ok
}

func NewSchedulerToContext(ctx context.Context, tick time.Duration) {
	c := NewScheduler(ctx, tick)
	setContextValue(ctx, schedulerKey, c)
}

func SchedulerChanFromContext(ctx context.Context) (chan callback, bool) {
	c, ok := getContextValue(ctx, schedulerKey).(chan callback)
	return c, ok
}
