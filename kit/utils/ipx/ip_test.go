package ipx

import (
	"log"
	"testing"
)

func TestIP(t *testing.T) {

	var ip = "119.23.144.144"
	//var ip = "24.48.0.1"
	baidu, err := GetIpInfoByBaidu(ip)
	log.Println(err)
	log.Printf("%+v", baidu)

	api, err := GetIpInfoByApi(ip)
	log.Println(err)
	log.Printf("%+v", api)
}
