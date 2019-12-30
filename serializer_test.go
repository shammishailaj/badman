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

func serializerCommonTest(t *testing.T, ser badman.Serializer) {
	t1, t2, t3 := time.Now(), time.Now(), time.Now()
	entities := []badman.BadEntity{
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
	ch := make(chan *badman.BadEntityMessage, 1)
	go func() {
		for i := range entities {
			ch <- &badman.BadEntityMessage{Entity: &entities[i]}
		}
		close(ch)
	}()

	err := ser.Serialize(ch, buf)
	require.NoError(t, err)

	raw := buf.Bytes()
	assert.NotEqual(t, 0, len(raw))

	reader := bytes.NewReader(raw)
	readCh := ser.Deserialize(reader)

	m1 := <-readCh
	require.NoError(t, m1.Error)
	assert.Equal(t, "blue", m1.Entity.Name)
	assert.Equal(t, "tester1", m1.Entity.Src)
	assert.Equal(t, t1.Unix(), m1.Entity.SavedAt.Unix())

	m2 := <-readCh
	require.NoError(t, m2.Error)
	assert.Equal(t, "orange", m2.Entity.Name)
	assert.Equal(t, "tester1", m2.Entity.Src)
	assert.Equal(t, t2.Unix(), m2.Entity.SavedAt.Unix())

	m3 := <-readCh
	require.NoError(t, m3.Error)
	assert.Equal(t, "red", m3.Entity.Name)
	assert.Equal(t, "tester1", m3.Entity.Src)
	assert.Equal(t, t3.Unix(), m3.Entity.SavedAt.Unix())

	_, ok := <-readCh
	assert.False(t, ok)
}
