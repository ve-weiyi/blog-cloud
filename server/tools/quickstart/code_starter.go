package quickstart

import (
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/inject"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent"
)

type CodeStarter struct {
	Config

	inventMetas []*invent.TemplateMeta
	injectMetas []*inject.AstInjectMeta
}

func NewCodeStarter(config Config) *CodeStarter {
	return &CodeStarter{
		Config: config,
	}
}

func (s *CodeStarter) AddInventMetas(meta ...*invent.TemplateMeta) {
	s.inventMetas = append(s.inventMetas, meta...)
}

func (s *CodeStarter) AddInjectMetas(meta ...*inject.AstInjectMeta) {
	s.injectMetas = append(s.injectMetas, meta...)
}

func (s *CodeStarter) Execute() error {

	for _, item := range s.inventMetas {
		if s.IsIgnoreKey(item.Key) {
			continue
		}

		err := item.CreateTempFile()
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, item := range s.injectMetas {
		if s.IsIgnoreKey(item.Key) {
			continue
		}
		err := item.Inject()
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (s *CodeStarter) RollBack() error {

	for _, item := range s.inventMetas {
		if s.IsIgnoreKey(item.Key) {
			continue
		}
		err := item.RollBack()
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, item := range s.injectMetas {
		if s.IsIgnoreKey(item.Key) {
			continue
		}
		err := item.RollBack()
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
