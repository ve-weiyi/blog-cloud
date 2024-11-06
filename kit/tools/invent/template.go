package invent

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

const (
	ModeCreateOrReplace = 0 //创建或替换
	ModeOnlyCreate      = 1 //只创建
	ModeOnlyReplace     = 2 //只替换
)

// 生成目录，文件名，模板内容，填充数据 data
type TemplateMeta struct {
	Key         string
	Mode        int    //模式 0:创建或替换 1:只创建 2:只替换
	CodeOutPath string //生成的代码路径  test/template.go

	//TemplateFile   string      //模版文件路径   tpl/api.go.tpl
	TemplateString string         //模版文件内容
	FunMap         map[string]any //模版函数
	Data           interface{}    //填充内容
}

func (s *TemplateMeta) Execute() error {

	switch s.Mode {
	case ModeOnlyCreate:
		if fileExist(s.CodeOutPath) {
			return fmt.Errorf("目标文件已存在:%s", s.CodeOutPath)
		}
	case ModeOnlyReplace:
		if !fileExist(s.CodeOutPath) {
			return fmt.Errorf("目标文件不存在:%s", s.CodeOutPath)
		}
	}

	//创建文件夹
	err := os.MkdirAll(filepath.Dir(s.CodeOutPath), 0755)
	if err != nil {
		return err
	}
	//创建.go文件
	f, err := os.Create(s.CodeOutPath)
	if err != nil {
		return err
	}
	defer f.Close()

	//解析模板
	temp, err := s.getTemplate()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = temp.Execute(&buf, s.Data)
	if err != nil {
		return err
	}

	if path.Ext(s.CodeOutPath) == ".go" {
		err = output(s.CodeOutPath, buf.Bytes())
		if err != nil {
			return err
		}
	} else {
		err = os.WriteFile(s.CodeOutPath, buf.Bytes(), 0640)
		if err != nil {
			return err
		}
	}

	fmt.Println("生成文件成功:", s.CodeOutPath)
	return nil
}

func (s *TemplateMeta) RollBack() error {
	//只创建不需要回滚
	if s.Mode == ModeOnlyCreate {
		return nil
	}

	err := deLFile(s.CodeOutPath)
	if err != nil {
		return err
	}
	return nil
}

func (s *TemplateMeta) MoveTempFile(movePath string) error {
	//判断目标文件是否都可以移动
	if movePath == "" {
		return nil
	}

	if fileExist(movePath) {
		return fmt.Errorf("目标文件已存在:%s", movePath)
	}

	err := moveFile(s.CodeOutPath, movePath)
	if err != nil {
		return err
	}

	fmt.Println("file move success:", movePath)
	return nil
}

func (s *TemplateMeta) getTemplate() (*template.Template, error) {
	//if s.TemplateFile != "" {
	//	//解析模板
	//	temp, err := template.ParseFiles(s.TemplateFile)
	//	temp.Funcs(s.FunMap)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return temp, nil
	//}

	if s.TemplateString != "" {
		//解析模板
		temp := template.New("temp")
		temp.Funcs(s.FunMap)
		_, err := temp.Parse(s.TemplateString)
		if err != nil {
			return nil, err
		}
		return temp, nil
	}

	return nil, fmt.Errorf("TemplateString or TemplateFile all null")
}
