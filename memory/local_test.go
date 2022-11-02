package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserExists(t *testing.T) {
	db, err := NewLocalDB("../memory.db")
	require.Nil(t, err)

	userPubkeyHash := "BAEB92BC6E8144F2D15E977A878ABFAC"

	err = db.SaveUser(User{
		PubkeyHash: userPubkeyHash,
		Nickname:   "Tester",
		LastSeen:   time.Now(),
	})
	require.Nil(t, err)

	isExists, err := db.IsUserExists(userPubkeyHash)
	require.Nil(t, err)
	assert.Equal(t, true, isExists)
}
