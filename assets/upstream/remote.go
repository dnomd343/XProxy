package upstream

import (
	"io"
	"net/http"
	"time"
)

type remoteAsset struct {
	client http.Client
}

func (r *remoteAsset) GetTag() string {
	return ""
}

func (r *remoteAsset) Request() (io.ReadCloser, time.Time, error) {
	return nil, time.Now(), nil
}
