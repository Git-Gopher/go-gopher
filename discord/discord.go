package discord

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

var (
	// Link time variable.
	DiscordToken = ""
	ChannelID    = "1006004369990369352"
)

func SendLog(filename string) error {
	dg, err := discordgo.New("Bot " + DiscordToken)
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
