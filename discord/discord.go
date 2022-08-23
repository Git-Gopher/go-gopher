package discord

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

var (
	//nolint: gosec
	// XXX: Feel free to send memes to our logging channel.
	Token     = "TVRBd05qQXhOVEkzTnpBMk16azVPVFV3T1EuRzNtLWFHLnNac3RUNi1icnBMWTEtZzVlMnl6Q25nVVdvN28wX2NlakdNSkNz"
	ChannelID = "1006004369990369352"
)

// Just so that nobody can run strings on our binary.
func DecodeToken() (string, error) {
	data, err := base64.StdEncoding.DecodeString(Token)
	if err != nil {
		return "", fmt.Errorf("failed to decode token: %w", err)
	}

	return string(data), nil
}

func SendLog(filename string) error {
	token, err := DecodeToken()
	if err != nil {
		return fmt.Errorf("failed to decode token: %w", err)
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
