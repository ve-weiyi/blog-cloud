package ipx

import (
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

const (
	FieldsAll = "66846719"
	FieldsStd = "61439"
)

type Location struct {
	Query         string  `json:"query"`         // 查询的 IP 地址
	Status        string  `json:"status"`        // 请求状态（例如 "success" 表示成功，"fail" 表示失败）
	Continent     string  `json:"continent"`     // 大洲名称（例如 "亚洲"）
	ContinentCode string  `json:"continentCode"` // 大洲代码（例如 "AS"）
	Country       string  `json:"country"`       // 国家名称（例如 "美国"）
	CountryCode   string  `json:"countryCode"`   // 国家代码（例如 "US"）
	Region        string  `json:"region"`        // 地区或州名称（例如 "加利福尼亚州"）
	RegionName    string  `json:"regionName"`    // 地区或州名称（例如 "加利福尼亚州"）
	City          string  `json:"city"`          // 城市名称（例如 "旧金山"）
	District      string  `json:"district"`      // 区或地区名称（例如 "密西昂区"）
	Zip           string  `json:"zip"`           // 邮政编码（例如 "94110"）
	Lat           float64 `json:"lat"`           // 纬度坐标（例如 37.7749）
	Lon           float64 `json:"lon"`           // 经度坐标（例如 -122.4194）
	Timezone      string  `json:"timezone"`      // 时区（例如 "美国/洛杉矶"）
	Offset        int     `json:"offset"`        // 时区偏移（单位为秒，例如 -25200 表示 UTC-8）
	Currency      string  `json:"currency"`      // 货币代码（例如 "USD"）
	Isp           string  `json:"isp"`           // 网络服务提供商（例如 "AT&T 服务公司"）
	Org           string  `json:"org"`           // 组织名称（例如 "谷歌有限责任公司"）
	As            string  `json:"as"`            // 自治系统号码（例如 "AS15169"）
	Asname        string  `json:"asname"`        // 自治系统名称（例如 "谷歌有限责任公司"）
	Reverse       string  `json:"reverse"`       // 反向 DNS 查询结果（例如 "google.com"）
	Mobile        bool    `json:"mobile"`        // 表示 IP 是否来自移动网络（true 表示是，false 表示否）
	Proxy         bool    `json:"proxy"`         // 表示 IP 是否为已知代理（true 表示是，false 表示否）
	Hosting       bool    `json:"hosting"`       // 表示 IP 是否来自托管提供商（true 表示是，false 表示否）
}

func GetIpInfoByApi(ip string) (*Location, error) {
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN&fields=%s", ip, FieldsStd)

	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var location Location
	err = jsoniter.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}
