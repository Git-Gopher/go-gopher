package discord

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrDiscordTokenNotSet = errors.New("DISCORD_TOKEN is not set")
	ChannelID             = "1006004369990369352"
)

func SendLog(filename string) error {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return ErrDiscordTokenNotSet
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	reader, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	_, err = dg.ChannelFileSend(ChannelID, filename, reader)

	if err != nil {
		return fmt.Errorf("failed to send file: %w", err)
	}

	return nil
}
