package internal

import (
	"context"
	"fmt"
	"github.com/dinedal/astrotime"
	"github.com/jinzhu/gorm"
	"net"
	"net/url"
	"strings"
	"time"
)

type Camera struct {
	gorm.Model
	URL       string
	Lat       float64
	Lng       float64
	Sunrise   time.Time // Next Sunrise
	Sunset    time.Time // Next Sunset
	FailCount uint
}

func NewCamera(ctx context.Context, _url string) (*Camera, error) {
	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	h := strings.Split(u.Host, ":")

	host := h[0]
	if h == nil || host == "" {
		return nil, fmt.Errorf("Bad host")
	}
	ip := net.ParseIP(h[0])
	db, ok := GeoDBFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Could not get GeoDB reader from context")
	}
	record, err := db.City(ip)
	if err != nil {
		return nil, err
	}
	lat := record.Location.Latitude
	lng := record.Location.Longitude
	sunrise := astrotime.NextSunrise(time.Now(), lat, lng)
	sunset := astrotime.NextSunset(time.Now(), lat, lng)
	return &Camera{
		URL:     _url,
		Lat:     lat,
		Lng:     lng,
		Sunrise: sunrise,
		Sunset:  sunset,
	}, nil
}
