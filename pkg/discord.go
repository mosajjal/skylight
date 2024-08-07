package skylight

import (
	"errors"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/phuslu/log"
)

// this file takes a gofeed `item` and converts it to a discord `embed` object

func itemToEmbed(item *FeedNews) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		URL:   item.Link,
		Title: item.Title,
	}

	if item.Image != nil {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: item.Image.URL,
		}
	} else {
		log.Debug().Msgf("No image found for %s", item.Title)
	}

	if item.Author != nil {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name: item.FeedName,
		}
	} else {
		log.Debug().Msgf("No author found for %s", item.Title)
	}

	if item.PublishedParsed != nil {
		embed.Timestamp = item.PublishedParsed.Format(time.RFC3339)
	} else {
		log.Debug().Msgf("No published date found for %s", item.Title)
	}

	return embed
}

func parseWebhookURL(webhookURL string) (id string, token string, err error) {
	// split the webhook url into the id and token
	// the url should be in the format: https://discord.com/api/webhooks/ID/TOKEN
	// we need to extract the ID and TOKEN
	// the ID is the second last part of the url
	// the TOKEN is the last part of the url
	parts := strings.Split(webhookURL, "/")
	if len(parts) < 6 {
		return "", "", errors.New("invalid webhook url")
	}
	return parts[len(parts)-2], parts[len(parts)-1], nil
}

// SendToWebhook sends an embedded message to a discord webhook
func SendToWebhook(webhookURL string, item *FeedNews) error {
	embed := itemToEmbed(item)

	// create a new discord session
	session, err := discordgo.New("Bot " + webhookURL)
	if err != nil {
		log.Warn().Msgf("Failed to create discord session: %s", err)
		return err
	}

	// create the webhook message
	message := discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	}

	// parse the webhook url
	id, token, err := parseWebhookURL(webhookURL)
	if err != nil {
		return err
	}

	_, err = session.WebhookExecute(id, token, false, &message)
	if err != nil {
		return err
	}

	return nil
}
