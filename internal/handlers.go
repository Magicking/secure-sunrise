package internal

import (
	"context"
	"html"
	"log"

	models "github.com/Magicking/secure-sunrise/models"
	op "github.com/Magicking/secure-sunrise/restapi/operations"
	middleware "github.com/go-openapi/runtime/middleware"
)

func AddUrls(ctx context.Context, params op.AddUrlsParams) middleware.Responder {
	go func() {
		for _, _url := range params.Urls {
			_url = html.UnescapeString(_url)
			cam, err := NewCamera(ctx, _url)
			if err != nil {
				log.Printf("Could not add url: %v", _url)
				continue
			}
			if err = InsertCamera(ctx, cam); err != nil {
				log.Printf("Could not add camera %q: %v", _url, err)
			}
		}
		log.Printf("Added %d cameras", len(params.Urls))
	}()
	return op.NewAddUrlsOK()
}

func Getfeeds(ctx context.Context, params op.GetfeedsParams) middleware.Responder {
	fm, ok := FeedManagerFromContext(ctx)
	if !ok {
		err_str := "Could not obtain FeedManager from context"
		log.Println(err_str)
		return op.NewGetfeedsDefault(500).WithPayload(&models.Error{Message: &err_str})
	}
	sunrise, err := fm.GetFeed(params.Name)
	if err != nil {
		log.Println(err)
		err_str := err.Error()
		return op.NewGetfeedsDefault(500).WithPayload(&models.Error{Message: &err_str})
	}
	urls := sunrise.CurrentURLs
	return op.NewGetfeedsOK().WithPayload(urls)
}
