package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	guildID := os.Getenv("DISCORD_GUILD_ID")
	channelName := os.Getenv("DISCORD_CHANNEL_NAME")

	if botToken == "" || guildID == "" || channelName == "" {
		fmt.Println("Missing required environment variables: DISCORD_BOT_TOKEN, DISCORD_GUILD_ID, DISCORD_CHANNEL_NAME")
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		os.Exit(1)
	}

	channels, err := dg.GuildChannels(guildID)
	if err != nil {
		fmt.Println("Error retrieving channels:", err)
		os.Exit(1)
	}

	var oldChannel *discordgo.Channel
	for _, ch := range channels {
		if ch.Name == channelName {
			oldChannel = ch
			break
		}
	}

	if oldChannel == nil {
		fmt.Printf("Channel '%s' not found\n", channelName)
		os.Exit(1)
	}

	newChannel, err := dg.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:                 oldChannel.Name,
		Type:                 oldChannel.Type,
		Topic:                oldChannel.Topic,
		Position:             oldChannel.Position,
		PermissionOverwrites: oldChannel.PermissionOverwrites,
		ParentID:             oldChannel.ParentID,
		RateLimitPerUser:     oldChannel.RateLimitPerUser,
	})
	if err != nil {
		fmt.Println("Error duplicating channel:", err)
		os.Exit(1)
	}

	_, err = dg.ChannelDelete(oldChannel.ID)
	if err != nil {
		fmt.Println("Error deleting old channel:", err)
		os.Exit(1)
	}

	fmt.Printf("Channel '%s' duplicated and original deleted. New channel ID: %s\n", newChannel.Name, newChannel.ID)
}
