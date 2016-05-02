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
