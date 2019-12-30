package source_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/m-mizutani/badman"
	"github.com/m-mizutani/badman/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyHTTPClient struct {
	Req  *http.Request
	Resp *http.Response
}

func (x *dummyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	x.Req = req
	return x.Resp, nil
}

func TestMalwareDomains(t *testing.T) {
	fd, err := os.Open("test/malwaredomains/domains.txt")
	require.NoError(t, err)

	dummy := &dummyHTTPClient{
		Resp: &http.Response{
			StatusCode: 200,
			Body:       fd,
		},
	}

	source.InjectNewHTTPClient(dummy)
	defer source.FixNewHTTPClient()

	var entities []*badman.BadEntity
	src := source.NewMalwareDomains()
	for q := range src.Download() {
		require.NoError(t, q.Error)
		entities = append(entities, q.Entities...)
	}

	assert.Equal(t, 2, len(entities))
	assert.Equal(t, "blue.example.com", entities[0].Name)
	assert.Equal(t, "MalwareDomains", entities[0].Src)
	assert.Equal(t, "phishing", entities[0].Reason)

	assert.Equal(t, "orange.example.net", entities[1].Name)
	assert.Equal(t, "MalwareDomains", entities[1].Src)
	assert.Equal(t, "exploit", entities[1].Reason)
}

func TestMVPs(t *testing.T) {
	fd, err := os.Open("test/mvps/hosts.txt")
	require.NoError(t, err)
	dummy := &dummyHTTPClient{
		Resp: &http.Response{
			StatusCode: 200,
			Body:       fd,
		},
	}

	source.InjectNewHTTPClient(dummy)
	defer source.FixNewHTTPClient()

	var entities []*badman.BadEntity
	for q := range source.NewMVPS().Download() {
		require.NoError(t, q.Error)
		entities = append(entities, q.Entities...)
	}

	assert.Equal(t, 2, len(entities))
	assert.Equal(t, "blue.example.com", entities[0].Name)
	assert.Equal(t, "MVPs", entities[0].Src)
	assert.Equal(t, "", entities[0].Reason)

	assert.Equal(t, "orange.example.net", entities[1].Name)
	assert.Equal(t, "MVPs", entities[1].Src)
	assert.Equal(t, "", entities[1].Reason)
}

func TestURLhausRecent(t *testing.T) {
	fd, err := os.Open("test/urlhaus/test.csv")
	require.NoError(t, err)
	dummy := &dummyHTTPClient{
		Resp: &http.Response{
			StatusCode: 200,
			Body:       fd,
		},
	}

	source.InjectNewHTTPClient(dummy)
	defer source.FixNewHTTPClient()

	var entities []*badman.BadEntity
	for q := range source.NewURLhausRecent().Download() {
		require.NoError(t, q.Error)
		entities = append(entities, q.Entities...)
	}

	assert.Equal(t, "/downloads/csv_recent/", dummy.Req.URL.Path)

	assert.Equal(t, 2, len(entities))
	assert.Equal(t, "blue.example.com", entities[0].Name)
	assert.Equal(t, "URLhaus", entities[0].Src)
	assert.Equal(t, "malware_download", entities[0].Reason)

	assert.Equal(t, "orange.example.net", entities[1].Name)
	assert.Equal(t, "URLhaus", entities[1].Src)
	assert.Equal(t, "malware_download", entities[1].Reason)
}

func TestURLhausOnline(t *testing.T) {
	fd, err := os.Open("test/urlhaus/test.csv")
	require.NoError(t, err)
	dummy := &dummyHTTPClient{
		Resp: &http.Response{
			StatusCode: 200,
			Body:       fd,
		},
	}

	source.InjectNewHTTPClient(dummy)
	defer source.FixNewHTTPClient()

	var entities []*badman.BadEntity
	for q := range source.NewURLhausOnline().Download() {
		require.NoError(t, q.Error)
		entities = append(entities, q.Entities...)
	}

	assert.Equal(t, "/downloads/csv_online/", dummy.Req.URL.Path)

	assert.Equal(t, 2, len(entities))
	assert.Equal(t, "blue.example.com", entities[0].Name)
	assert.Equal(t, "URLhaus", entities[0].Src)
	assert.Equal(t, "malware_download", entities[0].Reason)

	assert.Equal(t, "orange.example.net", entities[1].Name)
	assert.Equal(t, "URLhaus", entities[1].Src)
	assert.Equal(t, "malware_download", entities[1].Reason)
}
