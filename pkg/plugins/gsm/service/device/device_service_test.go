package device

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/gsm/service/brand"
)

func TestGetDeviceList(t *testing.T) {
	// get all brands first to get a slug
	_, err := brand.GetAllBrands()
	assert.Equal(t, nil, err)

	// get device list
	devices, err := GetDeviceList("apple-phones-48", 1)

	js, _ := json.MarshalIndent(devices, "", " ")
	t.Logf("Devices: %v", string(js))
}
