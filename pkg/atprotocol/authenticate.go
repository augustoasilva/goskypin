package atprotocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type SessionRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func CreateSession() (*DIDResponse, error) {
	request := SessionRequest{
		Identifier: os.Getenv("BLUESKY_IDENTIFIER"),
		Password:   os.Getenv("BLUESKY_PASSWORD"),
	}
	requestBody, _ := json.Marshal(request)

	url := fmt.Sprintf("%s/com.atproto.server.createSession", BskyXrpcURL)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("status code inesperado", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("status code inesperado: %d", resp.StatusCode)
	}

	var didResponse DIDResponse
	if jsonDecoderErr := json.NewDecoder(resp.Body).Decode(&didResponse); jsonDecoderErr != nil {
		slog.Error(
			"error ao decodificar a resposta do servidor",
			"error", jsonDecoderErr,
		)
		return nil, jsonDecoderErr
	}

	return &didResponse, nil
}
