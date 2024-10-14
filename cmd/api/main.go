package main

import (
	"log/slog"

	"github.com/augustoasilva/goskyrepost/internal/handler"
	"github.com/augustoasilva/goskyrepost/pkg/atprotocol"
)

func main() {
	slog.Info("iniciando o servidor do bot de repost")
	if firehoseErr := atprotocol.Firehose(handler.EventFn); firehoseErr != nil {
		slog.Error("erro ao escutar a firehose do bluesky", "error", firehoseErr)
		panic(firehoseErr)
	}
}
