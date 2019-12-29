package badman

import (
	"bufio"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Source is interface of BlackList.
type Source interface {
	Download() chan *BadEntityMessage
}

// DefaultSources is default set of blacklist source that is maintained by badman.
var DefaultSources = []Source{
	NewMalwareDomains(),
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

// MalwareDomains downloads blacklist from http://www.malwaredomains.com/
type MalwareDomains struct {
	URL string
}

// NewMalwareDomains is constructor of MalwareDomains
func NewMalwareDomains() *MalwareDomains {
	return &MalwareDomains{
		URL: "http://mirror1.malwaredomains.com/files/domains.txt",
	}
}

// Download of MalwareDomains downloads domains.txt and parses to extract domain names.
func (x *MalwareDomains) Download() chan *BadEntityMessage {
	ch := make(chan *BadEntityMessage, defaultSourceChanSize)

	go func() {
		defer close(ch)

		now := time.Now()

		req, err := http.NewRequest("GET", x.URL, nil)
		if err != nil {
			ch <- &BadEntityMessage{
				Error: errors.Wrapf(err, "Fail to craete new MalwareDomains HTTP request to: %s", x.URL),
			}
			return
		}

		client := newHTTPClient()
		resp, err := client.Do(req)
		if err != nil {
			ch <- &BadEntityMessage{
				Error: errors.Wrapf(err, "Fail to send HTTP request to: %s", x.URL),
			}
			return
		}
		if resp.StatusCode != 200 {
			ch <- &BadEntityMessage{
				Error: errors.Wrapf(err, "Unexpected status code (%d): %s", resp.StatusCode, x.URL),
			}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue // Comment line
			}

			row := strings.Split(line, "\t")
			ch <- &BadEntityMessage{
				Entity: &BadEntity{
					Name:    row[2],
					SavedAt: now,
					Src:     "MalwareDomains",
					Reason:  row[3],
				},
			}
		}
	}()

	return ch
}

// MVPS downloads blacklist from http://winhelp2002.mvps.org/hosts.txt
type MVPS struct {
	URL string
}

// NewMVPS is constructor of MVPS
func NewMVPS() *MVPS {
	return &MVPS{
		URL: "http://winhelp2002.mvps.org/hosts.txt",
	}
}

// Download of MVPS downloads domains.txt and parses to extract domain names.
func (x *MVPS) Download() chan *BadEntityMessage {
	ch := make(chan *BadEntityMessage, defaultSourceChanSize)

	go func() {
		defer close(ch)

		now := time.Now()

		req, err := http.NewRequest("GET", x.URL, nil)
		if err != nil {
			ch <- &BadEntityMessage{
				Error: errors.Wrapf(err, "Fail to craete new MVPS HTTP request to: %s", x.URL),
			}
			return
		}

		client := newHTTPClient()
		resp, err := client.Do(req)
		if err != nil {
			ch <- &BadEntityMessage{
				Error: errors.Wrapf(err, "Fail to send HTTP request to: %s", x.URL),
			}
			return
		}
		if resp.StatusCode != 200 {
			ch <- &BadEntityMessage{
				Error: errors.Wrapf(err, "Unexpected status code (%d): %s", resp.StatusCode, x.URL),
			}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			row := strings.Split(line, " ")

			if len(row) == 2 && row[0] == "0.0.0.0" {
				ch <- &BadEntityMessage{
					Entity: &BadEntity{
						Name:    row[1],
						SavedAt: now,
						Src:     "MVPs",
					},
				}
			}
		}
	}()

	return ch
}
