package repost

import (
	"log/slog"

	"github.com/augustoasilva/goskyrepost/pkg/atprotocol"
)

func Repost(p *atprotocol.PostRecord) error {
	token, err := atprotocol.CreateSession()
	if err != nil {
		slog.Error("Error getting token", "error", err)
		return err
	}

	recordParams := createRecordParams{
		DIDResponse: token,
		Resource:    "app.bsky.feed.repost",
		URI:         p.Reply.Root.URI,
		CID:         p.Reply.Root.CID,
	}

	err = createRecord(recordParams)
	if err != nil {
		return err
	}

	recordParams.Resource = "app.bsky.feed.like"
	err = createRecord(recordParams)
	if err != nil {
		return err
	}

	return nil
}
