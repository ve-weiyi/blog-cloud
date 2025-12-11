package nacos

import (
	"fmt"
	"testing"
)

type ConfigReloaderMock struct {
}

func (sel *ConfigReloaderMock) Reload(data string) error {
	fmt.Println(data)
	return nil
}

func TestNewNacos(t *testing.T) {

}
