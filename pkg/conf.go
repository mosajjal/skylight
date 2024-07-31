package skylight

import "time"

type Config struct {
	LogLevel  string `hcl:"log_level"`
	StateFile string `hcl:"state_file"`

	Feeds []FeedConfig `hcl:"feed,block"`
}

type FeedConfig struct {
	Timestamp       time.Time
	Name            string `hcl:"name,label"`
	Url             string `hcl:"url"`
	WebhookURL      string `hcl:"webhook_url"`
	Interval        uint   `hcl:"interval"`
	MaxItems        int    `hcl:"max_items"`
	IgnoreOlderThen uint   `hcl:"ignore_items_older_than"` // in hours
}
