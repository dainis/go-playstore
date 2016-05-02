package playstore

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func createClient() *Playstore {
	c, err := New(os.Getenv("PLAYSTORE_EMAIL"), os.Getenv("PLAYSTORE_PASSWORD"), os.Getenv("ANDROID_ID"))
	if err != nil {
		log.Panicf("Failed to create client %e\n", err)
	}
	return c
}

func TestLogin(t *testing.T) {
	_, err := New(os.Getenv("PLAYSTORE_EMAIL"), os.Getenv("PLAYSTORE_PASSWORD"), os.Getenv("ANDROID_ID"))

	assert.Nil(t, err)

	_, err = New("trolololo@test.com", "trolololo", "6666666")

	assert.Error(t, err)
}
