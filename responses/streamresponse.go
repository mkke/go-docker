package responses

import (
	"bufio"
	"encoding/json"
	"io"

	logif "github.com/mkke/go-log"
	"github.com/pkg/errors"
)

type StreamResponse struct {
	Stream string `json:"stream"`
}

func ParseStreamBody(r io.ReadCloser, log logif.Logger) error {
	defer r.Close()

	s := bufio.NewScanner(r)
	for s.Scan() {
		var resp StreamResponse
		if err := json.Unmarshal(s.Bytes(), &resp); err != nil {
			return errors.Wrapf(err, "could not unmarshal '%s'", s.Text())
		}
		log.Println(resp.Stream)
	}

	return s.Err()
}
