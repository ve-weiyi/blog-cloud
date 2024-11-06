/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"log"

	"github.com/ve-weiyi/ve-blog-golang/server/cmd"
)

//go:generate go env -w GOPATH=./
//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// @title						Swagger Example API
// @version						1.0.0
// @description					APP接口鉴权方式：token采用自加密token的方式，自加密token规则：sha256(${x-timestamp}${base_secret}) 生成的加密token，x-timestamp 时间戳需要保持10分钟内
// @termsOfService				https://swagger.io/terms/
// @contact.name				API Support
// @contact.url					https://www.swagger.io/support/
// @contact.email				support@swagger.io
// @license.name				Apache 2.0
// @license.url					https://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:9090
// @BasePath
// @securityDefinitions.basic	BasicAuth
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						x-token
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	cmd.NewRootCmd().Execute()
}
