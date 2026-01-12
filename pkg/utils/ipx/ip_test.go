package ipx

import (
	"log"
	"testing"
)

func TestIP(t *testing.T) {

	var ip = "46.243.122.48"
	//var ip = "24.48.0.1"
	baidu, err := GetIpInfoByBaidu(ip)
	log.Println(err)
	log.Printf("%+v", baidu)

	api, err := GetIpInfoByApi(ip)
	log.Println(err)
	log.Printf("%+v", api)
}
