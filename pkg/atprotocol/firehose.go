package atprotocol

import (
	"bytes"
	"io"
	"log/slog"

	"github.com/fxamacker/cbor/v2"
	"github.com/gorilla/websocket"
)

func Firehose(handleEvent func(evt RepoCommitEvent) error) error {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		slog.Error("falha ao conectar no websocket", "error", err)
		return err
	}
	defer conn.Close()

	slog.Info("conectado ao websocket", "url", wsURL)

	for {
		_, message, errMessage := conn.ReadMessage()
		if errMessage != nil {
			slog.Error("error ao ler a mensagem do websocket", "error", errMessage)
			return errMessage
		}

		decoder := cbor.NewDecoder(bytes.NewReader(message))

		for {
			var evt RepoCommitEvent
			decodeErr := decoder.Decode(&evt)
			if decodeErr != nil {
				if decodeErr == io.EOF {
					break
				}
				slog.Error("error ao decodificar a mensagem de commit do repo", "error", decodeErr)
				return decodeErr
			}

			if handleErr := handleEvent(evt); handleErr != nil {
				panic(handleErr)
			}
		}
	}
}
