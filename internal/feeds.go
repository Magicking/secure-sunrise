package internal

import (
	"context"
	"fmt"
	"log"
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
	go feed.Run(ctx) // TODO Add channel to close
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
				//UPDATE ASTROTIME
			}
			if cam.Sunset.Before(now) {
				//UPDATE ASTROTIME
			}
			// SAVE CAM
		}
		return nil
	})
}

func (fm *FeedManager) GetFeed(name string) (*Feed, error) {
	feed, ok := fm.Feeds[name]
	if !ok {
		return nil, fmt.Errorf("Feed %q not found", name)
	}
	return feed, nil
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
		urls := f.getNextCurrentUrls(ctx)
		f.CurrentURLs = urls
		log.Printf("URLS(%v): %v", len(urls), urls)
		return nil
	})
}

func (f *Feed) getNextCurrentUrls(ctx context.Context) []string {
	now := time.Now()
	duration := 30 * time.Minute
	end := now.Add(duration)
	cameras, err := GetCameras(ctx, f.isSunrise, now, end) // TODO Do Sunset a New Feed level
	if err != nil {
		log.Printf("Error fetching samples from database: %v", err)
		return nil
	}
	ret := make([]string, len(cameras))
	for i, e := range cameras {
		ret[i] = e.URL
	}
	return ret
}
