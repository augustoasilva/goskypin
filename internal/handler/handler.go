package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/augustoasilva/go-lazuli/pkg/lazuli"
	"github.com/augustoasilva/go-lazuli/pkg/lazuli/bsky"
	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
	carv2 "github.com/ipld/go-car/v2"
)

var terms = []string{"#repostbot"}

func GetEventFn(ctx context.Context, client lazuli.Client) func(evt bsky.CommitEvent) error {
	handler := func(evt bsky.CommitEvent) error {
		for _, op := range evt.GetOps() {
			if op.Action == "create" {
				if len(evt.GetBlocks()) > 0 {
					err := handleCARBlocks(ctx, evt, op, client)
					if err != nil {
						slog.Error("error processing car blocks", "error", err)
						return err
					}
				}
			}
		}
		return nil
	}
	return handler
}

func handleCARBlocks(ctx context.Context, evt bsky.CommitEvent, op bsky.RepoOperation, client lazuli.Client) error {
	if len(evt.GetBlocks()) == 0 {
		return errors.New("there is no blocks to process")
	}

	reader, err := carv2.NewBlockReader(bytes.NewReader(evt.GetBlocks()))
	if err != nil {
		slog.Error("error creating reader for car blocks", "error", err)
		return err
	}

	for {
		block, readErr := reader.Next()
		if readErr == io.EOF {
			break
		}

		if readErr != nil {
			slog.Error("error reading car block", "error", readErr)
			break
		}

		c, cidErr := getCidFromOp(op)
		if cidErr != nil {
			slog.Error("error getting CID from events operation", "error", cidErr)
			continue
		}

		if block.Cid().Equals(*c) {
			var rawData map[string]any
			_ = cbor.Unmarshal(block.RawData(), &rawData)
			if rawData["$type"] != nil {
				slog.Info("car block read", "op_path", op.Path, "raw_data", rawData)
				if rawData["$type"].(string) == "app.bsky.feed.post" {
					var postRecord bsky.PostRecord
					if unmarshalErr := cbor.Unmarshal(block.RawData(), &postRecord); unmarshalErr != nil {
						slog.Error("error decoding car block using CBOR", "error", unmarshalErr)
					}

					if postRecord.Reply != nil && postRecord.Text != "" {
						if filterTerms(postRecord.Text) {
							postParams := bsky.CreateRecordParams{
								URI: postRecord.Reply.Parent.URI,
								CID: postRecord.Reply.Parent.CID,
							}

							repostErr := client.CreateRepostRecord(ctx, postParams)
							if repostErr != nil {
								slog.Error("error reposting", "error", repostErr)
								continue
							}

							likeErr := client.CreateLikeRecord(ctx, postParams)
							if likeErr != nil {
								slog.Error("error liking", "error", likeErr)
							}
						}
					}

					if postRecord.Text != "" {
						if filterTerms(postRecord.Text) {
							uri := "at://" + evt.GetRepo() + "/" + op.Path
							post, getErr := client.GetPost(ctx, uri)
							if getErr != nil {
								slog.Error("error getting post data", "error", getErr)
								continue
							}

							postParams := bsky.CreateRecordParams{
								URI: post.URI,
								CID: post.CID.(string),
							}

							repostErr := client.CreateRepostRecord(ctx, postParams)
							if repostErr != nil {
								slog.Error("error reposting", "error", repostErr)
								continue
							}

							likeErr := client.CreateLikeRecord(ctx, postParams)
							if likeErr != nil {
								slog.Error("error liking", "error", likeErr)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func getCidFromOp(op bsky.RepoOperation) (*cid.Cid, error) {
	if opTag, ok := op.CID.(cbor.Tag); ok {
		if cidBytes, ok := opTag.Content.([]byte); ok {
			return decodeCID(cidBytes)
		}
	}
	return nil, errors.New("nenhum CID encontrado na operação")
}

func decodeCID(cidBytes []byte) (*cid.Cid, error) {
	c, err := cid.Decode(string(cidBytes))
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar o CID: %w", err)
	}

	return &c, nil
}

func filterTerms(text string) bool {
	for _, term := range terms {
		if strings.Contains(strings.ToLower(text), strings.ToLower(term)) {
			return true
		}
	}
	return false
}
