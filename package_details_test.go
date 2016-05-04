package playstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetails(t *testing.T) {
	c := createClient()

	d, _ := c.PackageDetails("com.twitter.android")

	assert.Equal(t, "Twitter", d.Title)
	assert.Equal(t, "com.twitter.android", d.Id)
}

func TestPhonyDetails(t *testing.T) {
	c := createClient()

	d, err := c.PackageDetails("trolololo.trolololo.com")

	assert.Nil(t, d)
	assert.Error(t, err)
}
