package badman_test

import (
	"io/ioutil"
	"net/http"
	"strings"
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
	sampleData := `## This is sample data
		blue.example.com	phishing	openphish.com	20171117	20160527	20160108
		orange.example.net	exploit	xxx.com	20171117	20160527	20160108
`
	dummy := &dummyHTTPClient{
		Resp: &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(sampleData)),
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
