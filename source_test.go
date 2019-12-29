package badman_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/m-mizutani/badman"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyHTTPClient struct {
	Req  *http.Request
	Resp *http.Response
}

func (x *dummyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return x.Resp, nil
}

func TestMalwareDomains(t *testing.T) {
	fd, err := os.Open("sample/malwaredomains/domains.txt")
	require.NoError(t, err)

	dummy := &dummyHTTPClient{
		Resp: &http.Response{
			StatusCode: 200,
			Body:       fd,
		},
	}

	badman.InjectNewHTTPClient(dummy)
	defer badman.FixNewHTTPClient()

	var entities []*badman.BadEntity
	src := badman.NewMalwareDomains()
	for msg := range src.Download() {
		require.NoError(t, msg.Error)
		entities = append(entities, msg.Entity)
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
	fd, err := os.Open("sample/mvps/hosts.txt")
	require.NoError(t, err)
	dummy := &dummyHTTPClient{
		Resp: &http.Response{
			StatusCode: 200,
			Body:       fd,
		},
	}

	badman.InjectNewHTTPClient(dummy)
	defer badman.FixNewHTTPClient()

	var entities []*badman.BadEntity
	for msg := range badman.NewMVPS().Download() {
		require.NoError(t, msg.Error)
		entities = append(entities, msg.Entity)
	}

	assert.Equal(t, 2, len(entities))
	assert.Equal(t, "blue.example.com", entities[0].Name)
	assert.Equal(t, "MVPs", entities[0].Src)
	assert.Equal(t, "", entities[0].Reason)

	assert.Equal(t, "orange.example.net", entities[1].Name)
	assert.Equal(t, "MVPs", entities[1].Src)
	assert.Equal(t, "", entities[1].Reason)
}
