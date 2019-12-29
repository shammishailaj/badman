package badman

import (
	"bufio"
	"encoding/csv"
	"io"
	"net/http"
	"net/url"
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
	NewMVPS(),
	NewURLhausRecent(),
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

func getHttpBody(url string, ch chan *BadEntityMessage) io.Reader {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- &BadEntityMessage{
			Error: errors.Wrapf(err, "Fail to craete new MalwareDomains HTTP request to: %s", url),
		}
		return nil
	}

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		ch <- &BadEntityMessage{
			Error: errors.Wrapf(err, "Fail to send HTTP request to: %s", url),
		}
		return nil
	}
	if resp.StatusCode != 200 {
		ch <- &BadEntityMessage{
			Error: errors.Wrapf(err, "Unexpected status code (%d): %s", resp.StatusCode, url),
		}
		return nil
	}

	return resp.Body
}

// Download of MalwareDomains downloads domains.txt and parses to extract domain names.
func (x *MalwareDomains) Download() chan *BadEntityMessage {
	ch := make(chan *BadEntityMessage, defaultSourceChanSize)

	go func() {
		defer close(ch)

		now := time.Now()
		body := getHttpBody(x.URL, ch)
		if body == nil {
			return
		}

		scanner := bufio.NewScanner(body)
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
		body := getHttpBody(x.URL, ch)
		if body == nil {
			return
		}

		scanner := bufio.NewScanner(body)
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

// URLhausRecent downloads blacklist from https://urlhaus.abuse.ch/downloads/csv_recent/
// The blacklist has only URLs in recent 30 days.
type URLhausRecent struct {
	URL string
}

// NewURLhausRecent is constructor of URLhausRecent
func NewURLhausRecent() *URLhausRecent {
	return &URLhausRecent{
		URL: "https://urlhaus.abuse.ch/downloads/csv_recent/",
		//		"https://urlhaus.abuse.ch/downloads/csv_online/",
	}
}

// Download of URLhausRecent downloads domains.txt and parses to extract domain names.
func (x *URLhausRecent) Download() chan *BadEntityMessage {
	ch := make(chan *BadEntityMessage, defaultSourceChanSize)

	go func() {
		defer close(ch)

		body := getHttpBody(x.URL, ch)
		if body == nil {
			return
		}

		reader := csv.NewReader(body)
		reader.Comment = []rune("#")[0]

		for {
			row, err := reader.Read()
			if err == io.EOF {
				return
			} else if err != nil {
				ch <- &BadEntityMessage{Error: errors.Wrapf(err, "Fail to read CSV of URLhaus")}
				return
			}

			if len(row) != 8 {
				continue
			}

			url, err := url.Parse(row[2])
			if err != nil {
				ch <- &BadEntityMessage{Error: errors.Wrapf(err, "Fail to parse URL in URLhaus CSV")}
				return
			}

			ts, err := time.Parse("2006-01-02 15:04:05", row[1])
			if err != nil {
				ch <- &BadEntityMessage{Error: errors.Wrapf(err, "Fail to parse tiemstamp in URLhaus CSV")}
				return
			}

			ch <- &BadEntityMessage{
				Entity: &BadEntity{
					Name:    url.Hostname(),
					SavedAt: ts,
					Src:     "URLhaus",
					Reason:  row[4],
				},
			}
		}
	}()

	return ch
}
