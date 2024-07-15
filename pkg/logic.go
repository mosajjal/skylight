package skylight

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/knadh/koanf"
	"github.com/mmcdole/gofeed"
)

type Feed struct {
	Timestamp time.Time
}

var state struct {
	initiated  bool // used to check if the state has been loaded from the file
	feedsState map[string]*Feed
}

func RSSItem2TGMessage(item *gofeed.Item) error {
	//TODO
	return nil

}

// run this as a goroutine forever
func (f *Feed) HandleFeed(url string, refreshInterval time.Duration) {
	// fetch the feed
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	var latestTimestamp time.Time
	for {
		// iterate over the items
		for i, item := range feed.Items {
			// check if the item is new
			if f.Timestamp.Before(*item.PublishedParsed) {
				// do something with the item
				if i == 0 {
					latestTimestamp = *item.PublishedParsed
				}
				fmt.Println(item.Title)
				fmt.Println(RSSItem2TGMessage(item))
			}
			if i >= config.General.MaxEpisodesPerPodcast {
				break
			}
		}
		// update the timestamp
		f.Timestamp = latestTimestamp
		// wait for the refresh interval
		time.Sleep(refreshInterval)
	}
}

// helper function to get a state filepath and updates the state struct (global variable)
func StateHandler(stateFilepath string) {

	// read or create the state file
	stateFile, err := os.Open(stateFilepath)
	if err == os.ErrNotExist {
		os.WriteFile(stateFilepath, []byte{}, 0644)
		stateFile, _ = os.Open(stateFilepath)
	}

	// if the state has already been initiated, no need to load it again, but we need to write back the state to the file
	if state.initiated {
		go func() {
			for {
				time.Sleep(5 * time.Minute)
				// truncate the file before re-writing the state
				stateFile, _ := os.Create(stateFilepath)
				encoder := gob.NewEncoder(stateFile)
				encoder.Encode(&state)
			}
		}()
	}

	// now use the statefile as a GOB
	decoder := gob.NewDecoder(stateFile)
	if decoder.Decode(&state) != nil {
		// if there's an error decoding the state, we need to create a new state on disk
		state.feedsState = make(map[string]*Feed)
	}
	// no matter if the state was loaded or created, we need to set the initiated flag to true
	state.initiated = true
}

func Run(stateFilepath string, k *koanf.Koanf) {
	// run the state handler
	StateHandler(stateFilepath)

	// get the feeds from the config
	//TODO: build the config struct and run a goroutine per feed

}
