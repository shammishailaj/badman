package badman_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/m-mizutani/badman"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONSerializer(t *testing.T) {
	ser := badman.NewJSONSerializer()
	serializerCommonTest(t, ser)
}

func TestGzipJSONSerializer(t *testing.T) {
	ser := badman.NewGzipJSONSerializer()
	serializerCommonTest(t, ser)
}

func TestMsgpackSerializer(t *testing.T) {
	ser := badman.NewMsgpackSerializer()
	serializerCommonTest(t, ser)
}

func TestGzipMsgpackSerializer(t *testing.T) {
	ser := badman.NewGzipMsgpackSerializer()
	serializerCommonTest(t, ser)
}

func serializerCommonTest(t *testing.T, ser badman.Serializer) {
	t1, t2, t3 := time.Now(), time.Now(), time.Now()
	entities := []*badman.BadEntity{
		{
			Name:    "blue",
			SavedAt: t1,
			Src:     "tester1",
		},
		{
			Name:    "orange",
			SavedAt: t2,
			Src:     "tester1",
		},
		{
			Name:    "red",
			SavedAt: t3,
			Src:     "tester1",
		},
	}

	buf := &bytes.Buffer{}
	ch := make(chan *badman.EntityQueue, 1)
	go func() {
		ch <- &badman.EntityQueue{Entities: entities}
		close(ch)
	}()

	err := ser.Serialize(ch, buf)
	require.NoError(t, err)

	raw := buf.Bytes()
	assert.NotEqual(t, 0, len(raw))

	reader := bytes.NewReader(raw)
	readCh := ser.Deserialize(reader)

	var recvEntities []*badman.BadEntity
	for q := range readCh {
		require.NoError(t, q.Error)
		for _, e := range q.Entities {
			recvEntities = append(recvEntities, e)
		}
	}

	assert.Equal(t, "blue", recvEntities[0].Name)
	assert.Equal(t, "tester1", recvEntities[0].Src)
	assert.Equal(t, t1.Unix(), recvEntities[0].SavedAt.Unix())

	assert.Equal(t, "orange", recvEntities[1].Name)
	assert.Equal(t, "tester1", recvEntities[1].Src)
	assert.Equal(t, t2.Unix(), recvEntities[1].SavedAt.Unix())

	assert.Equal(t, "red", recvEntities[2].Name)
	assert.Equal(t, "tester1", recvEntities[2].Src)
	assert.Equal(t, t3.Unix(), recvEntities[2].SavedAt.Unix())
}
