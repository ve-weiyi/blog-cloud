package tempx

import (
	"bytes"
	"text/template"
)

// 解析text模版并填充实体数据
func TempParseString(tempString string, tempContent interface{}) (string, error) {
	temple, err := template.New("temp").Parse(tempString) // （2）解析模板
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = temple.Execute(buf, tempContent) //（3）数据驱动模板，将tempContent的值填充到模板中，存入buf
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 解析file模版并填充实体数据,注意需要传入文件绝对路径
func TempParseFile(filenames string, tempContent interface{}) (string, error) {
	temple, err := template.ParseFiles(filenames) // （2）解析模板
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = temple.Execute(buf, tempContent) //（3）数据驱动模板，将tempContent的值填充到模板中，存入buf
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
