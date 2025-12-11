package brand

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllBrands(t *testing.T) {
	brands, err := GetAllBrands()
	assert.Equal(t, nil, err)

	firstBrandName := "Acer"
	assert.Equal(t, firstBrandName, brands[0].Name)

	js, _ := json.MarshalIndent(brands, "", " ")
	fmt.Println("Brands:", string(js))
}
