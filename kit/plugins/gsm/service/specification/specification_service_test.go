package specification

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeviceList(t *testing.T) {
	specs, err := GetSpecification("apple_iphone_16_pro_max-13123")

	assert.Equal(t, nil, err)
	assert.Equal(t, "Acer", specs.Brand)
	assert.Equal(t, "Acer Liquid Z4", specs.DeviceName)
	assert.Equal(t, "Phone", specs.DeviceType)

	js, _ := json.MarshalIndent(specs, "", " ")
	t.Logf("Devices: %v", string(js))
}
