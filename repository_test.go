package badman_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/badman"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRepository(t *testing.T) {
	repo := badman.NewInMemoryRepository()
	repositoryCommonTest(repo, t)
}

func repositoryCommonTest(repo badman.Repository, t *testing.T) {
	e1 := badman.BadEntity{
		Name:    "10.0.0.1",
		SavedAt: time.Now(),
		Src:     "tester1",
	}
	e2 := badman.BadEntity{
		Name:    "blue.example.com",
		SavedAt: time.Now(),
		Src:     "tester2",
	}
	e3 := badman.BadEntity{
		Name:    "blue.example.com",
		SavedAt: time.Now(),
		Src:     "tester3",
	}
	e4 := badman.BadEntity{
		Name:    "orange.example.net",
		SavedAt: time.Now(),
		Src:     "tester3",
	}

	// No entity in repository
	r1, err := repo.Get("10.0.0.1")
	require.NoError(t, err)
	assert.Nil(t, r1)
	r2, err := repo.Get("blue.example.com")
	require.NoError(t, err)
	assert.Nil(t, r2)

	// Insert entities
	require.NoError(t, repo.Put(e1))
	require.NoError(t, repo.Put(e2))
	require.NoError(t, repo.Put(e3))
	require.NoError(t, repo.Put(e4))

	// Get operations
	r3, err := repo.Get("10.0.0.1")
	require.NoError(t, err)
	assert.NotNil(t, r3)
	require.Equal(t, 1, len(r3))
	assert.Equal(t, "10.0.0.1", r3[0].Name)

	r4, err := repo.Get("blue.example.com")
	require.NoError(t, err)
	assert.NotNil(t, r4)
	require.Equal(t, 2, len(r4))
	assert.Equal(t, "blue.example.com", r4[0].Name)
	assert.Equal(t, "blue.example.com", r4[1].Name)
	if r4[0].Src == "tester2" {
		assert.Equal(t, "tester3", r4[1].Src)
	} else {
		assert.Equal(t, "tester2", r4[1].Src)
	}

	// Delete operation
	r5, err := repo.Get("orange.example.net")
	require.NoError(t, err)
	assert.NotNil(t, r5)
	require.Equal(t, 1, len(r5))
	assert.Equal(t, "orange.example.net", r5[0].Name)

	err = repo.Del("orange.example.net")
	require.NoError(t, err)
	r6, err := repo.Get("orange.example.net")
	require.NoError(t, err)
	assert.Equal(t, 0, len(r6))

	// Dump operation
	counter := map[string]int{}
	for msg := range repo.Dump() {
		require.NoError(t, msg.Error)
		counter[msg.Entity.Name]++
	}
	assert.Equal(t, 1, counter["10.0.0.1"])
	assert.Equal(t, 2, counter["blue.example.com"])
	assert.Equal(t, 0, counter["orange.example.net"])

	// Clear operation
	require.NoError(t, repo.Clear())
	r7, err := repo.Get("blue.example.com")
	require.NoError(t, err)
	assert.Equal(t, 0, len(r7))
}
