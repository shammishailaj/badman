package source

import (
	"encoding/csv"
	"io"
	"net/url"
	"time"

	"github.com/m-mizutani/badman"
	"github.com/pkg/errors"
)

func downloadURLhasu(csvURL string, ch chan *badman.BadEntityMessage) {
	defer close(ch)

	body := getHTTPBody(csvURL, ch)
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
			ch <- &badman.BadEntityMessage{
				Error: errors.Wrapf(err, "Fail to read CSV of URLhaus"),
			}
			return
		}

		if len(row) != 8 {
			continue
		}

		url, err := url.Parse(row[2])
		if err != nil {
			ch <- &badman.BadEntityMessage{
				Error: errors.Wrapf(err, "Fail to parse URL in URLhaus CSV"),
			}
			return
		}

		ts, err := time.Parse("2006-01-02 15:04:05", row[1])
		if err != nil {
			ch <- &badman.BadEntityMessage{
				Error: errors.Wrapf(err, "Fail to parse tiemstamp in URLhaus CSV"),
			}
			return
		}

		ch <- &badman.BadEntityMessage{
			Entity: &badman.BadEntity{
				Name:    url.Hostname(),
				SavedAt: ts,
				Src:     "URLhaus",
				Reason:  row[4],
			},
		}
	}
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
	}
}

// Download of URLhausRecent downloads domains.txt and parses to extract domain names.
func (x *URLhausRecent) Download() chan *badman.BadEntityMessage {
	ch := make(chan *badman.BadEntityMessage, defaultSourceChanSize)
	go downloadURLhasu(x.URL, ch)
	return ch
}

// URLhausOnline downloads blacklist from https://urlhaus.abuse.ch/downloads/csv_recent/
// The blacklist has only online URLs.
type URLhausOnline struct {
	URL string
}

// NewURLhausOnline is constructor of URLhausOnline
func NewURLhausOnline() *URLhausOnline {
	return &URLhausOnline{
		URL: "https://urlhaus.abuse.ch/downloads/csv_online/",
	}
}

// Download of URLhausOnline downloads domains.txt and parses to extract domain names.
func (x *URLhausOnline) Download() chan *badman.BadEntityMessage {
	ch := make(chan *badman.BadEntityMessage, defaultSourceChanSize)
	go downloadURLhasu(x.URL, ch)
	return ch
}
