package internal

import (
	"context"
	"fmt"
	"html"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	op "github.com/Magicking/secure-sunrise/restapi/operations"
	"github.com/dinedal/astrotime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/oschwald/geoip2-golang"
)

type Cam struct {
	Url    *url.URL
	Record *geoip2.City
}

func NewCam(ctx context.Context, url *url.URL) (*Cam, error) {
	h := strings.Split(url.Host, ":")

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
	return &Cam{
		Url:    url,
		Record: record,
	}, nil
}

type Feed struct {
	Cams map[string][]*Cam
}

func (f *Feed) AddCam(cam *Cam) {
	// save cam with minitz hhmms0
	t := astrotime.NextSunrise(time.Now(), cam.Record.Location.Latitude, cam.Record.Location.Longitude)
	minitz := fmt.Sprintf("%02d%02d", t.Hour(), t.Minute())

	if f.Cams == nil {
		f.Cams = make(map[string][]*Cam)
	}
	f.Cams[minitz] = append(f.Cams[minitz], cam)
}

func (f *Feed) GetCurrentUrls() []string {
	now := time.Now()
	duration := 30 * time.Minute
	end := now.Add(duration)
	var ret []string
	for idx := now; idx.Before(end); idx = idx.Add(1 * time.Minute) {
		minitz := fmt.Sprintf("%02d%02d", idx.Hour(), idx.Minute())
		cams, ok := f.Cams[minitz]
		if !ok {
			continue
		}
		for _, e := range cams {
			ret = append(ret, e.Url.String())
		}
	}
	return ret
}

var feed Feed

func AddUrls(ctx context.Context, params op.AddUrlsParams) middleware.Responder {
	//	return op.NewAddUrlsDefault(500).WithPayload(&models.Error{Message: &err_str})
	f := &feed
	for _, _url := range params.Urls {
		_url = html.UnescapeString(_url)
		u, err := url.Parse(_url)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		cam, err := NewCam(ctx, u)
		if err != nil {
			log.Printf("Could not add url: %v", _url)
			continue
		}
		f.AddCam(cam)
	}
	return op.NewAddUrlsOK()
}

func Getfeeds(ctx context.Context, params op.GetfeedsParams) middleware.Responder {
	//	return op.NewGetfeedsDefault(500).WithPayload(&models.Error{Message: &err_str})
	f := &feed
	return op.NewGetfeedsOK().WithPayload(f.GetCurrentUrls())
}
