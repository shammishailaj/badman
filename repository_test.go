package badman_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/m-mizutani/badman"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	fd, err := os.Open("data/hosts.txt")
	require.NoError(t, err)
	r := badman.NewRepository()
	r.Read(fd)
	fmt.Println("read")
	time.Sleep(100 * time.Second)
}
