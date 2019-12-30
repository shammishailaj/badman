package source

import (
	"bufio"
	"strings"
	"time"

	"github.com/m-mizutani/badman"
)

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
func (x *MVPS) Download() chan *badman.BadEntityMessage {
	ch := make(chan *badman.BadEntityMessage, defaultSourceChanSize)

	go func() {
		defer close(ch)

		now := time.Now()
		body := getHTTPBody(x.URL, ch)
		if body == nil {
			return
		}

		scanner := bufio.NewScanner(body)
		for scanner.Scan() {
			line := scanner.Text()
			row := strings.Split(line, " ")

			if len(row) == 2 && row[0] == "0.0.0.0" {
				ch <- &badman.BadEntityMessage{
					Entity: &badman.BadEntity{
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
