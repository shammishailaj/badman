package source

import (
	"io"
	"net/http"

	"github.com/m-mizutani/badman"
	"github.com/pkg/errors"
)

// DefaultSet is default set of blacklist source that is maintained by badman.
var DefaultSet = []badman.Source{
	NewMalwareDomains(),
	NewMVPS(),
	NewURLhausRecent(),
	NewURLhausOnline(),
}

const defaultSourceChanSize = 128

// httpClient interface is used to inject own client for testing.
// InjectNewHTTPClient in export_test.go replace constructor to use
// dummy HTTP client and FixNewHTTPClient in export_test.go reverts it.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func newNormalHTTPClient() httpClient { return &http.Client{} }

var newHTTPClient = newNormalHTTPClient

func getHTTPBody(url string, ch chan *badman.BadEntityMessage) io.Reader {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- &badman.BadEntityMessage{
			Error: errors.Wrapf(err, "Fail to craete new MalwareDomains HTTP request to: %s", url),
		}
		return nil
	}

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		ch <- &badman.BadEntityMessage{
			Error: errors.Wrapf(err, "Fail to send HTTP request to: %s", url),
		}
		return nil
	}
	if resp.StatusCode != 200 {
		ch <- &badman.BadEntityMessage{
			Error: errors.Wrapf(err, "Unexpected status code (%d): %s", resp.StatusCode, url),
		}
		return nil
	}

	return resp.Body
}
