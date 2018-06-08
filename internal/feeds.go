package internal

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/dinedal/astrotime"
	uuid "github.com/satori/go.uuid"
	//	"os/exec"

	"time"
)

type FeedManager struct {
	Feeds map[string]*Feed
}

func NewFeedManager(ctx context.Context) *FeedManager {
	var fm FeedManager
	fm.Feeds = make(map[string]*Feed)

	go fm.Run(ctx)
	return &fm
}

func (fm *FeedManager) NewFeed(ctx context.Context, name string, isSunrise bool) {
	feed := NewFeed(isSunrise)
	fm.Feeds[name] = feed
	go feed.Run(ctx)
}

func (fm *FeedManager) Run(ctx context.Context) {
	// TODO
	// Calculate Next Sunrise and Next Sunset for past sunrise & sunset
	c, ok := SchedulerChanFromContext(ctx)
	if !ok {
		log.Fatalf("Could not obtain Scheduler chan from context")
	}
	c <- callback(func(ctx context.Context) error {
		now := time.Now()
		cams, err := GetPastCameras(ctx, now)
		if err != nil {
			return fmt.Errorf("FeedManager Runner: %v", err)
		}
		for _, cam := range cams {
			if cam.Sunrise.Before(now) {
				cam.Sunrise = astrotime.NextSunrise(now, cam.Lat, cam.Lng)
			}
			if cam.Sunset.Before(now) {
				cam.Sunset = astrotime.NextSunset(now, cam.Lat, cam.Lng)
			}
			UpdateCam(ctx, cam)
		}
		return nil
	})
}

func (fm *FeedManager) GetFeed(ctx context.Context, name string) (*Feed, error) {
	feed, ok := fm.Feeds[name]
	if !ok {
		return nil, fmt.Errorf("Feed %q not found", name)
	}
	return &Feed{CurrentURLs: feed.CurrentURLs}, nil
	//return feed, nil
}

// Holds currents samples / urls
type Feed struct {
	isSunrise   bool
	CurrentURLs []string
}

func NewFeed(sunrise bool) *Feed {
	return &Feed{isSunrise: sunrise}
}

// Cache Next cameras to display in memory
func (f *Feed) Run(ctx context.Context) {
	// Main idea:
	// Every X minutes
	// Get urls to display
	// Feed URLs to sampler
	// Sampler populate feeder with expiration time
	c, ok := SchedulerChanFromContext(ctx)
	if !ok {
		log.Fatalf("Could not obtain Scheduler chan from context")
	}
	c <- callback(func(ctx context.Context) error {
		urls := f.GetNextCurrentUrls(ctx)
		f.CurrentURLs = urls
		log.Printf("URLS(%v): %v", len(urls), urls)
		return nil
	})
}

func (f *Feed) GetNextCurrentUrls(ctx context.Context) []string {
	now := time.Now()
	duration := 30 * time.Minute
	end := now.Add(duration)
	cameras, err := GetCameras(ctx, f.isSunrise, now, end)
	if err != nil {
		log.Printf("Error fetching samples from database: %v", err)
		return nil
	}
	camerasUniq := make(map[string]struct{})
	ret := make([]string, len(cameras)) //max cameras
	for _, e := range cameras {
		if _, ok := camerasUniq[e.URL]; ok {
			continue
		}
		camerasUniq[e.URL] = struct{}{}
		// Below code to get sample and put it backend
		sampleTime := 10 * time.Second // TODO move to configuration
		// allow 2 times the maximum sample time
		//		ctx, cancel := context.WithTimeout(context.Background(), sampleTime*3)
		//		defer cancel()

		url := e.URL
		out := fmt.Sprintf("%s", uuid.Must(uuid.NewV4()))
		duration := Time(sampleTime).String()
		// TODO put get_sample.sh to configuration
		if err := exec.CommandContext(context.Background(), "/get_sample.sh", url, out, duration).Run(); err != nil {
			log.Printf("Failed to get sample for %v", url)
			continue
		}
		//retURL <- out
		//log.Println("Current url", e.CurrentSample, e.URL)
		ret = append(ret, out)
	}
	return ret
}
