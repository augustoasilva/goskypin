package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli"
	"github.com/augustoasilva/goskyrepost/internal/handler"
)

func main() {
	ctx := context.Background()
	xrpcURL := os.Getenv("XRPC_URL")
	wsURL := os.Getenv("WS_URL")
	client := lazuli.NewClient(xrpcURL, wsURL)
	slog.Debug("created lazuli client", "client", client)

	identifier := os.Getenv("IDENTIFIER")
	password := os.Getenv("PASSWORD")
	_, sessErr := client.CreateSession(ctx, identifier, password)
	if sessErr != nil {
		slog.Error("error creating session", "error", sessErr)
		panic(sessErr)
	}

	slog.Debug("started session to client", "client", client)
	slog.Info("starting bot to listen to bluesky and repost", "client", client)

	if firehoseErr := client.ConsumeFirehose(ctx, handler.GetEventFn(ctx, client)); firehoseErr != nil {
		slog.Error("error listening bluesky firehose websocket", "error", firehoseErr)
		panic(firehoseErr)
	}
}
