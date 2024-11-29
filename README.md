# Skylight

Skylight is a Go application that monitors RSS feeds and sends updates to Discord channels via webhooks. It allows you to configure multiple RSS feeds, set how often to check for updates, and customize notification settings.

## Features

- Monitor multiple RSS feeds.
- Send updates to Discord via webhooks.
- Customize check intervals and notification settings.
- Filter feed items based on age.
- Persistent state management to avoid duplicate notifications.

## Installation

Ensure you have Go installed on your system. Then, download and build Skylight:

```bash
git clone https://github.com/yourusername/skylight.git
cd skylight
go build -o skylight
```

## Usage

Run Skylight with a configuration file:

```bash
./skylight -config path/to/config.hcl
```

Options:

- `-config`: Path to the configuration file.
- `-defaultconfig`: Print the default configuration to stdout.

## Configuration

Skylight uses a configuration file in HCL format. Below is an example 

config.hcl:

```hcl
log_level = "info"
state_file = "state.json"

feed "Example Feed" {
    url = "https://example.com/rss"
    interval = 300  # in seconds
    max_items = 10
    ignore_items_older_than = 24  # in hours
    webhook_url = "https://discord.com/api/webhooks/your_webhook_id/your_webhook_token"
}
```

### Parameters

- `log_level`: Logging level (`debug`, `info`, `warn`, `error`).
- `state_file`: File to store the state between runs.

#### Feed Configuration

Each `feed` block represents an RSS feed to monitor.

- `name`: Unique name for the feed.
- `url`: RSS feed URL.
- `interval`: How often to check for updates (in seconds).
- `max_items`: Maximum number of items to keep track of.
- `ignore_items_older_than`: Ignore items older than this value (in hours). Set to `0` to keep all items.
- `webhook_url`: Discord webhook URL for notifications.

## Discord Webhook Setup

To receive notifications, set up a Discord webhook:

1. Open your Discord server settings.
2. Navigate to **Integrations** > **Webhooks**.
3. Click **New Webhook** and customize it.
4. Copy the **Webhook URL** and use it in your 


## Examples

### Multiple Feeds

```hcl
log_level = "info"
state_file = "state.json"

feed "Hacker News" {
    url = "https://news.ycombinator.com/rss"
    interval = 300
    max_items = 10
    ignore_items_older_than = 24
    webhook_url = "https://discord.com/api/webhooks/your_webhook_id/your_webhook_token"
}

feed "Tech Crunch" {
    url = "http://feeds.feedburner.com/TechCrunch/"
    interval = 600
    max_items = 15
    ignore_items_older_than = 48
    webhook_url = "https://discord.com/api/webhooks/your_webhook_id/your_webhook_token"
}
```
