package excel

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestExport(t *testing.T) {
	var err error
	e := NewExcelExporter()

	err = e.NewActiveSheet("Sheet1")
	t.Log(err)

	err = e.SetSheetTitle([]any{"姓名", "年龄", "性别"})
	t.Log(err)

	err = e.AddRowValue([]any{"张三", 18, "男"})
	t.Log(err)

	err = e.AddRowValue([]any{"李四", 18, "男"})
	t.Log(err)

	err = e.NewActiveSheet("Sheet2")
	t.Log(err)

	err = e.SetSheetTitle([]any{"姓名", "年龄", "性别"})
	t.Log(err)

	err = e.AddRowValue([]any{"张三", 18, "男"})
	t.Log(err)

	err = e.AddRowValue([]any{"李四", 18, "男"})
	t.Log(err)

	err = e.ExportFile("test.xlsx")
	if err != nil {
		t.Log(err)
	}
}

func TestImport(t *testing.T) {
	var err error
	i := NewExcelImporter()

	err = i.OpenFile("test.xlsx")
	t.Log(err)

	err = i.SetActiveSheet("Sheet1")
	t.Log(err)

	row, err := i.GetRowValue(0)
	t.Log(row, err)

	row, err = i.GetRowValue(1)
	t.Log(row, err)

	row, err = i.GetRowValue(2)
	t.Log(row, err)

	err = i.SetActiveSheet("Sheet2")
	t.Log(err)

	row, err = i.GetRowValue(1)
	t.Log(row, err)

	row, err = i.GetRowValue(2)
	t.Log(row, err)
}

func TestImport2(t *testing.T) {
	var err error
	i := NewExcelImporter()

	err = i.OpenFile("/Users/weiyi/Downloads/beta-sleep_data.xlsx")
	t.Log(err)

	err = i.SetActiveSheet("Sheet1")
	t.Log(err)

	title, err := i.GetSheetTitle()
	t.Log(title, err)

	row, err := i.GetRowValue(2)
	t.Log(row, err)

}

func TestNewExcelExporter(t *testing.T) {
	f := excelize.NewFile()
	ex := &ExcelExportImpl{
		File: f,
	}

	t.Log(ex.point(1, 1))
}
