package tmpl

const NotEditMark = `
`

const Header = NotEditMark + `
package {{.LowerStartCamelName}}

import(	
	{{range .ImportPkgPaths}}{{.}}` + "\n" + `{{end}}
)
`

const KeyApi = "api"
const KeyRouter = "router"
const KeyController = "controller"
const KeyService = "service"
const KeyRepository = "repository"
const KeyModel = "model"
