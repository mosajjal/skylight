package skylight

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/phuslu/log"
)

var runtimeState struct {
	initiated  bool // used to check if the state has been loaded from the file
	FeedsState map[string]*FeedConfig
}

type FeedNews struct {
	*rss.Item
	FeedName string
}

// run this as a goroutine forever
func HandleFeed(f *FeedConfig) {
	log.Info().Msgf("Feed %s initiated", f.Name)
	// fetch the feed. first create a http client with 10 second timeout
	// TODO: move this inside the loop with ETAG and Last-Modified headers to make sure we don't over-fetch
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		fp, err := rss.FetchByClient(f.Url, httpClient)
		if err != nil {
			log.Error().Msgf("Error fetching feed %s: %s", f.Name, err)
		}
		// iterate over the items
		if fp != nil && err == nil {
			log.Info().Msgf("Fetched feed %s with %d items", f.Name, len(fp.Items))
			for i, item := range fp.Items {

				// check and fix the item's published date so every item has a date
				if item.Date.IsZero() {
					log.Warn().Msgf("Item %s has no published date, skipping", item.Title)
					// // check the feed updated date instead
					// if fp.Refresh != time.Unix(0, 0) {
					// 	item.Date = fp.Refresh
					// } else {
					// 	log.Warn().Msgf("Feed %s has no updated date", f.Name)
					// 	continue
					// }
					continue
				}

				// if the item is published older than ignore_items_older_than hours, skip it
				if f.IgnoreOlderThen != 0 {
					// check to see if the feed has a published date
					if time.Since(item.Date) > time.Duration(f.IgnoreOlderThen)*time.Hour {
						continue
					}
				}

				// check if the fetched feed item is newer than the latest timestamp recorded in the state
				if f.Timestamp.Before(item.Date) {

					fmt.Println(item.Title)
					if err := SendToWebhook(f.WebhookURL, &FeedNews{item, f.Name}); err != nil {
						log.Error().Msgf("Error sending to webhook: %s", err)
					}
					// update the timestamp
					log.Info().Msgf("Updating latest timestamp of feed %s to %s", f.Name, f.Timestamp)
					if !item.Date.IsZero() {
						if item.Date.After(f.Timestamp) {
							f.Timestamp = item.Date
						}
					}
				}
				if i >= f.MaxItems {
					break
				}
			}
		}
		// wait for the refresh interval
		log.Debug().Msgf("Waiting for %d seconds", f.Interval)
		time.Sleep(time.Duration(f.Interval * uint(time.Second)))
	}
}

// helper function to get a state filepath and updates the state struct (global variable)
func StateHandler(stateFilepath string) {

	// read or create the state file
	stateFile, err := os.Open(stateFilepath)
	if errors.Is(err, os.ErrNotExist) {
		os.WriteFile(stateFilepath, []byte{}, 0644)
		stateFile, _ = os.Open(stateFilepath)
	} else if err != nil {
		log.Error().Msgf("Error opening state file: %s", err)
	}

	// if the state has already been initiated, no need to load it again, but we need to write back the state to the file
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			if runtimeState.initiated {
				// truncate the file before re-writing the state
				stateFile2, _ := os.Create(stateFilepath)
				j, _ := json.Marshal(runtimeState)
				stateFile2.Write(j)
				stateFile2.Close()
			}
		}
	}()

	// now use the statefile as a JSON
	statefileBytes := make([]byte, 10240)
	if n, err := stateFile.Read(statefileBytes); err != nil {
		log.Error().Msgf("Error reading state file: %s", err)
		runtimeState.FeedsState = make(map[string]*FeedConfig)
	} else {
		if err := json.Unmarshal(statefileBytes[:n], &runtimeState); err != nil {
			log.Error().Msgf("Error decoding state file: %s", err)
			// if there's an error decoding the state, we need to create a new state on disk
			runtimeState.FeedsState = make(map[string]*FeedConfig)
		}
	}
	// no matter if the state was loaded or created, we need to set the initiated flag to true
	runtimeState.initiated = true
	stateFile.Close()
	select {}
}

func Run(c Config) {
	// run the state handler
	go StateHandler(c.StateFile)

	time.Sleep(1 * time.Second) // wait for the state to be loaded

	// get the feeds from the config
	for _, feed := range c.Feeds {
		// create the state if it doesn't exist
		if _, ok := runtimeState.FeedsState[feed.Name]; !ok {
			log.Debug().Msgf("Feed %s not found in state", feed.Name)
			feedCfg := feed
			// BUG: if the config changes, the state takes priority over the config, meaning the config becomes useless for that item.
			runtimeState.FeedsState[feed.Name] = &feedCfg
		}
		go HandleFeed(runtimeState.FeedsState[feed.Name])
	}
}
