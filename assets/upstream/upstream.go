package upstream

import (
	"io"
	urlpkg "net/url"
	"time"
)

// Upstream interface is an abstraction that supports initiating request and
// returning information such as data streams.
type Upstream interface {
	// Tag function get the description of current upstream.
	Tag() (tag string)

	// Request function initiates a request, returns the data stream and last
	// modification time, or an error information.
	Request() (stream io.ReadCloser, lastModify time.Time, err error)
}

func NewLocalAsset(tag string, path string) Upstream {
	return &localAsset{
		tag:  tag,
		path: path,
	}
}

func NewRemoteAsset(url string, proxy *urlpkg.URL) Upstream {
	return &remoteAsset{}
}
