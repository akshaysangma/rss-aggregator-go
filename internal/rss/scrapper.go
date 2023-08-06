package rss

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/akshaysangma/rss-aggregator-go/internal/database"
	"github.com/google/uuid"
)

func StartScrapper(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scrapping on %v goroutines with %v time between requests", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error while getting feeds to fetch: %v", err)
			continue
		}
		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, &wg, feed)
		}
		wg.Wait()
	}

}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	log.Printf("Scrapping feed %v", feed.Url)

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error while scrapping feed %v: %v", feed.Url, err)
		return
	}

	log.Printf("Feed %v has %v items", feed.Url, len(rssFeed.Channel.Items))

	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description.Valid = true
			description.String = item.Description
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error while parsing date %v : %v", item.PubDate, err)
		}

		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				FeedID:      feed.ID,
				Title:       item.Title,
				Url:         item.Link,
				Description: description,
				PublishedAt: pubAt,
			})
		if err != nil {
			if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Error while inserting posts: %v", err)
		}

		_, err = db.MarkFeedAsFetched(context.Background(), feed.ID)
		if err != nil {
			log.Printf("Error while updating feed last fetched: %v", err)
			return
		}
	}
}
