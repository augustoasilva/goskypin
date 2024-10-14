package atprotocol

import "time"

type Reply struct {
	Parent Parent `json:"parent"`
	Root   Root   `json:"root"`
}

type Parent struct {
	CID string `json:"cid"`
	URI string `json:"uri"`
}

type Root struct {
	CID string `json:"cid"`
	URI string `json:"uri"`
}

type PostRecord struct {
	LexiconTypeID string `json:"$type"`
	Text          string `json:"text"`
	Reply         *Reply `json:"reply"`
}

type RecordSubject struct {
	URI string `json:"uri"`
	CID string `json:"cid"`
}

type RequestRecord struct {
	RecordSubject `json:"subject"`
	CreatedAt     time.Time `json:"createdAt"`
}

type RequestRecordBody struct {
	LexiconTypeID string        `json:"$type"`
	Collection    string        `json:"collection"`
	Repo          string        `json:"repo"`
	Record        RequestRecord `json:"record"`
}
