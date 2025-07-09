package excelx

import (
	"fmt"
	"io"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter interface {
	// 导出为文件
	ExportFile(fileName string) error
	// 导出到指定位置
	ExportWriter(w io.Writer) error
	// 新建一个sheet页
	NewActiveSheet(sheetName string) error
	// 设置sheet页的标题
	SetSheetTitle(title []any) error
	// 设置sheet页的某一行的值
	SetRowValue(rowIndex int, rowData []any) error
	// 在末尾添加一列
	AddRowValue(rowData []any) error
	// 设置某个格子值
	SetCellValue(rowIndex int, cloIndex int, value any) error
	// 获取行数
	RowCount() int
}

type ExcelExportImpl struct {
	activeSheetName string
	File            *excelize.File
}

func NewExcelExporter() ExcelExporter {
	f := excelize.NewFile()
	return &ExcelExportImpl{
		File: f,
	}
}

func (s *ExcelExportImpl) NewActiveSheet(sheetName string) error {

	s.activeSheetName = sheetName

	// Create a new sheet.
	index, err := s.File.NewSheet(s.activeSheetName)
	if err != nil {
		return err
	}
	// Set active sheet of the workbook.
	s.File.SetActiveSheet(index)

	return nil
}

func (s *ExcelExportImpl) ExportFile(fileName string) error {
	return s.File.SaveAs(fileName)
}

func (s *ExcelExportImpl) ExportWriter(w io.Writer) error {
	return s.File.Write(w)
}

func (s *ExcelExportImpl) SetSheetTitle(title []any) error {
	return s.SetRowValue(1, title)
}

// excel 是从1开始的，为了方便rows对齐excel。所以这里的  rows下标=rowIndex-1
func (s *ExcelExportImpl) SetRowValue(rowIndex int, rowData []any) error {
	if rowIndex <= 0 {
		return fmt.Errorf("index must be greater than 0")
	}

	for i, v := range rowData {
		err := s.File.SetCellValue(s.activeSheetName, s.point(i+1, rowIndex), v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ExcelExportImpl) AddRowValue(rowData []any) error {
	count := s.RowCount()
	return s.SetRowValue(count+1, rowData)
}

func (s *ExcelExportImpl) SetCellValue(rowIndex int, cloIndex int, value any) error {
	return s.File.SetCellValue(s.activeSheetName, s.point(cloIndex, rowIndex), value)
}

func (s *ExcelExportImpl) RowCount() int {
	rows, err := s.File.GetRows(s.activeSheetName)
	if err != nil {
		return 0
	}
	return len(rows)
}

// (1,1),(2,2),(3,3) -> A1,B2,C3
func (s *ExcelExportImpl) point(charaId int, i int) string {
	row := strconv.Itoa(i)
	return s.charColIndex(charaId) + row
}

// 1,2,3 -> A,B,C
func (s *ExcelExportImpl) charColIndex(i int) string {
	return string(rune(64 + i))
}
