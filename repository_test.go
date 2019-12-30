package badman_test

import (
	"encoding/binary"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/badman"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRepository(t *testing.T) {
	repo := badman.NewInMemoryRepository()
	repositoryCommonTest(repo, t)
}

func repositoryCommonTest(repo badman.Repository, t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, rnd.Uint32())

	addr1 := ip.String()
	domain1 := uuid.New().String() + ".blue.example.com"
	domain2 := uuid.New().String() + ".orange.example.com"

	e1 := badman.BadEntity{
		Name:    addr1,
		SavedAt: time.Now(),
		Src:     "tester1",
	}
	e2 := badman.BadEntity{
		Name:    domain1,
		SavedAt: time.Now(),
		Src:     "tester2",
	}
	e3 := badman.BadEntity{
		Name:    domain1,
		SavedAt: time.Now(),
		Src:     "tester3",
	}
	e4 := badman.BadEntity{
		Name:    domain2,
		SavedAt: time.Now(),
		Src:     "tester3",
	}

	// No entity in repository
	r1, err := repo.Get(addr1)
	require.NoError(t, err)
	assert.Nil(t, r1)
	r2, err := repo.Get(domain1)
	require.NoError(t, err)
	assert.Nil(t, r2)

	// Insert entities
	require.NoError(t, repo.Put([]*badman.BadEntity{&e1, &e2, &e3, &e4}))

	// Get operations
	r3, err := repo.Get(addr1)
	require.NoError(t, err)
	assert.NotNil(t, r3)
	require.Equal(t, 1, len(r3))
	assert.Equal(t, addr1, r3[0].Name)

	r4, err := repo.Get(domain1)
	require.NoError(t, err)
	assert.NotNil(t, r4)
	require.Equal(t, 2, len(r4))
	assert.Equal(t, domain1, r4[0].Name)
	assert.Equal(t, domain1, r4[1].Name)
	if r4[0].Src == "tester2" {
		assert.Equal(t, "tester3", r4[1].Src)
	} else {
		assert.Equal(t, "tester2", r4[1].Src)
	}

	// Delete operation
	r5, err := repo.Get(domain2)
	require.NoError(t, err)
	assert.NotNil(t, r5)
	require.Equal(t, 1, len(r5))
	assert.Equal(t, domain2, r5[0].Name)

	err = repo.Del(domain2)
	require.NoError(t, err)
	r6, err := repo.Get(domain2)
	require.NoError(t, err)
	assert.Equal(t, 0, len(r6))

	// Dump operation
	counter := map[string]int{}
	for q := range repo.Dump() {
		require.NoError(t, q.Error)
		for _, e := range q.Entities {
			counter[e.Name]++
		}
	}
	assert.Equal(t, 1, counter[addr1])
	assert.Equal(t, 2, counter[domain1])
	assert.Equal(t, 0, counter[domain2])
}
