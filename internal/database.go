package internal

import (
	"fmt"
	"log"
	"time"

	"context"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var fPruneDB bool

func InsertCamera(ctx context.Context, cam *Camera) error {
	db, ok := DBFromContext(ctx)
	if !ok {
		return fmt.Errorf("Could not obtain DB from Context")
	}
	if err := db.Create(cam).Error; err != nil {
		return err
	}

	return nil
}

func UpdateCam(ctx context.Context, cam *Camera) (bool, error) {
	db, ok := DBFromContext(ctx)
	if !ok {
		return false, fmt.Errorf("Could not obtain DB from Context")
	}
	cur := db.Save(cam)
	if cur.Error != nil {
		return false, cur.Error
	}

	return cur.RowsAffected > 0, nil
}

func GetPastCameras(ctx context.Context, now time.Time) (ret []*Camera, err error) {
	db, ok := DBFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Could not obtain DB from Context")
	}
	cursor := db.Where("sunrise < ? or sunset < ?", now, now)
	if cursor.Error != nil {
		return nil, cursor.Error
	}
	if cursor.Find(&ret).RecordNotFound() {
		log.Println("No samples found in database, astro updater too fast ?")
		return nil, nil
	}
	return ret, nil
}

func GetCameras(ctx context.Context, sunrise bool, begin, end time.Time) (ret []*Camera, err error) {
	db, ok := DBFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Could not obtain DB from Context")
	}
	var whereString string
	if sunrise {
		whereString = "sunrise > ? AND sunrise < ?"
	} else {
		whereString = "sunset > ? AND sunset < ?"
	}
	cursor := db.Where(whereString, begin, end).Select("DISTINCT(URL)")
	if cursor.Error != nil {
		return nil, cursor.Error
	}
	if cursor.Find(&ret).RecordNotFound() {
		return nil, fmt.Errorf("No samples found in database")
	}
	return ret, nil
}

func DelCamera(ctx context.Context, url string) error {
	db, ok := DBFromContext(ctx)
	if !ok {
		return fmt.Errorf("Could not obtain DB from Context")
	}
	cursor := db.Where(Camera{URL: url}).Delete(Camera{})
	if cursor.Error != nil {
		return fmt.Errorf("Error deleting for URL (%v): %v", url, cursor.Error)
	}
	fPruneDB = true
	return nil
}

func RegisterDBCallback(ctx context.Context) error {
	c, ok := SchedulerChanFromContext(ctx)
	if !ok {
		return fmt.Errorf("Could not obtain Scheduler chan from context")
	}

	c <- callback(func(context.Context) error {
		if fPruneDB {
			db, ok := DBFromContext(ctx)
			if !ok {
				return fmt.Errorf("Could not obtain DB from Context")
			}
			var count uint
			db.Unscoped().Find(&Camera{}).Count(&count)
			log.Printf("Should prune %d Cameras", count)
			fPruneDB = false
		}
		return nil
	})
	return nil
}

func InitDatabase(dbDsn string) (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	for i := 1; i < 10; i++ {
		db, err = gorm.Open("postgres", dbDsn)
		if err == nil || i == 10 {
			break
		}
		sleep := (2 << uint(i)) * time.Second
		log.Printf("Could not connect to DB: %v", err)
		log.Printf("Waiting %v before retry", sleep)
		time.Sleep(sleep)
	}
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&Camera{}).Error; err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
