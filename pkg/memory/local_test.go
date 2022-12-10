package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserExists(t *testing.T) {
	db, err := NewLocalDB("test.db")
	require.Nil(t, err)

	userPubkeyHash := "BAEB92BC6E8144F2D15E977A878ABFAC"

	err = db.AddUser(User{
		PubkeyHash: userPubkeyHash,
		Nickname:   "Tester",
		LastSeen:   time.Now(),
	})
	require.Nil(t, err)

	isExists, err := db.IsUserExists(User{PubkeyHash: "test"})
	require.Nil(t, err)
	assert.Equal(t, true, isExists)

	isExists, err = db.IsUserExists(User{PubkeyHash: "test"})
	require.Nil(t, err)
	assert.Equal(t, false, isExists)
}
