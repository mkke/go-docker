package responses

import (
	"bufio"
	"encoding/json"
	"io"

	logif "github.com/mkke/go-log"
	"github.com/pkg/errors"
)

type PushStatus struct {
	Tag    string `json:"Tag"`
	Digest string `json:"Digest"`
}

type PushStatusResponse struct {
	Id       string      `json:"id"`
	Status   string      `json:"status"`
	Progress string      `json:"progress"`
	Aux      *PushStatus `json:"aux"`
}

func ParsePushStatusBody(r io.ReadCloser, log logif.Logger) (*PushStatus, error) {
	var pushStatus *PushStatus

	defer r.Close()

	s := bufio.NewScanner(r)
	prevStatus := ""
	for s.Scan() {
		var resp PushStatusResponse

		if err := json.Unmarshal(s.Bytes(), &resp); err != nil {
			return nil, errors.Wrapf(err, "could not unmarshal '%s'", s.Text())
		}
		if resp.Status != "" && prevStatus != resp.Status {
			log.Println(resp.Status)
			prevStatus = resp.Status
		}
		if resp.Aux != nil {
			pushStatus = resp.Aux
		}
	}

	if s.Err() != nil {
		return nil, s.Err()
	}

	if pushStatus == nil {
		return nil, errors.New("auxiliary information not sent")
	}

	return pushStatus, nil
}
