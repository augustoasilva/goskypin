package repost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/augustoasilva/goskyrepost/pkg/atprotocol"
)

type createRecordParams struct {
	DIDResponse *atprotocol.DIDResponse
	Resource    string
	URI         string
	CID         string
}

func createRecord(r createRecordParams) error {
	body := atprotocol.RequestRecordBody{
		LexiconTypeID: r.Resource,
		Collection:    r.Resource,
		Repo:          r.DIDResponse.DID,
		Record: atprotocol.RequestRecord{
			RecordSubject: atprotocol.RecordSubject{
				URI: r.URI,
				CID: r.CID,
			},
			CreatedAt: time.Now(),
		},
	}

	jsonBody, _ := json.Marshal(body)

	url := fmt.Sprintf("%s/com.atproto.repo.createRecord", atprotocol.BskyXrpcURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		slog.Error("erro ao criar objeto de request para criar um record", "error", err, "r.Resource", r.Resource)
		return nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.DIDResponse.AccessJwt))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("erro ao fazer a requisição para criar um record", "error", err, "r.Resource", r.Resource)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("status code inesperado ao tentar criar um record", "status", resp, "r.Resource", r.Resource)
		return nil
	}

	slog.Info("record criado com sucesso", "resource", r.Resource)

	return nil
}
