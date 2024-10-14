package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/augustoasilva/goskyrepost/internal/repost"
	"github.com/augustoasilva/goskyrepost/pkg/atprotocol"
	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
	carv2 "github.com/ipld/go-car/v2"
)

var terms = []string{"#repostbot"}

var EventFn = func(evt atprotocol.RepoCommitEvent) error {
	for _, op := range evt.Ops {
		if op.Action == "create" {
			if len(evt.Blocks) > 0 {
				err := handleCARBlocks(evt.Blocks, op)
				if err != nil {
					slog.Error("erro ao processar os blocos do CAR", "error", err)
					return err
				}
			}
		}
	}

	return nil
}

func handleCARBlocks(blocks []byte, op atprotocol.RepoOperation) error {
	if len(blocks) == 0 {
		return errors.New("não existem blocos para processar")
	}

	reader, err := carv2.NewBlockReader(bytes.NewReader(blocks))
	if err != nil {
		slog.Error("erro ao criar reader para ler os blocos do CAR", "error", err)
		return err
	}

	for {
		block, readErr := reader.Next()
		if readErr == io.EOF {
			break
		}

		if readErr != nil {
			slog.Error("erro ao ler um bloco do CAR", "error", readErr)
			break
		}

		c, cidErr := getCidFromOp(op)
		if cidErr != nil {
			slog.Error("erro ao pegar o CID da operação do evento", "error", cidErr)
			continue
		}

		if block.Cid().Equals(*c) {
			var post atprotocol.PostRecord
			if unmarshalErr := cbor.Unmarshal(block.RawData(), &post); unmarshalErr != nil {
				slog.Error("erro ao decodificar o bloco CAR usando CBOR", "error", unmarshalErr)
				continue
			}

			if post.Text == "" || post.Reply == nil {
				continue
			}

			if filterTerms(post.Text) {
				_ = repost.Repost(&post)
			}
		}

	}

	return nil
}

func getCidFromOp(op atprotocol.RepoOperation) (*cid.Cid, error) {
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
