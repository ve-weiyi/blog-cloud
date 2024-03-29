package apiparser

import (
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

func TestSwaggerParser(t *testing.T) {
	ac := NewSwaggerParser()
	apis, err := ac.ParseApiDocsByRoots(global.GetRuntimeRoot() + "server/docs")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.ObjectToJsonIndent(apis))

	models, err := ac.ParseModelDocsByRoots(global.GetRuntimeRoot() + "server/docs")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.ObjectToJsonIndent(models))
}

func TestAstParser(t *testing.T) {
	cfg := AstParserConfig{
		ApiBase: "/api/v1",
	}
	ac := NewAstParser(cfg)
	//apis, err := ac.ParseApiDocsByRoots(global.GetRuntimeRoot() + "server/api/controller")
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//t.Log(jsonconv.ObjectToJsonIndent(apis))

	models, err := ac.ParseModelDocsByRoots(global.GetRuntimeRoot() + "server/api/model")
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonconv.ObjectToJsonIndent(models))
}
