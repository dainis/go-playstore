package playstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDownload(t *testing.T) {
	c := createClient()

	d, _ := c.PackageDetails("com.twitter.android")

	pkg, err := c.DownloadPackage(d.Id, d.VersionCode)

	assert.Nil(t, err)

	assert.NotNil(t, pkg)
	assert.Equal(t, len(pkg) > 10000, true)
}
