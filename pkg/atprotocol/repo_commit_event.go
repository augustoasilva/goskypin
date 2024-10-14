package atprotocol

type RepoCommitEvent struct {
	Repo   string          `cbor:"repo"`
	Rev    string          `cbor:"rev"`
	Seq    int64           `cbor:"seq"`
	Since  string          `cbor:"since"`
	Time   string          `cbor:"time"`
	TooBig bool            `cbor:"tooBig"`
	Prev   interface{}     `cbor:"prev"`
	Rebase bool            `cbor:"rebase"`
	Blocks []byte          `cbor:"blocks"`
	Ops    []RepoOperation `cbor:"ops"`
}

type RepoOperation struct {
	Action string      `cbor:"action"`
	Path   string      `cbor:"path"`
	Reply  *Reply      `cbor:"reply"`
	Text   []byte      `cbor:"text"`
	CID    interface{} `cbor:"cid"`
}
