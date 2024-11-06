package excel

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

type ExcelImporter interface {
	// 打开文件
	OpenFile(fileName string) error
	// 打开reader
	OpenReader(r io.Reader) error
	// 设置活动的sheet页
	SetActiveSheet(sheetName string) error
	// 获取所有sheet页的标题
	GetSheetTitle() ([]string, error)
	// 获取某一行的值
	GetRowValue(rowIndex int) ([]string, error)
	// 获取所有行的值
	GetAllRowValue() ([][]string, error)
	// 获取所有行数
	RowCount() int
}

type ExcelImportImpl struct {
	activeSheetName string
	File            *excelize.File
}

func NewExcelImporter() ExcelImporter {
	f := excelize.NewFile()
	return &ExcelImportImpl{
		File: f,
	}
}

func (s *ExcelImportImpl) OpenFile(fileName string) error {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return err
	}

	s.File = f
	return nil
}

func (s *ExcelImportImpl) OpenReader(r io.Reader) error {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	s.File = f
	return nil
}

func (s *ExcelImportImpl) SetActiveSheet(sheetName string) error {

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

func (s *ExcelImportImpl) GetSheetTitle() ([]string, error) {
	return s.GetRowValue(1)
}

// excel 是从1开始的，为了方便rows对齐excel。所以这里的  rows下标=rowIndex-1
func (s *ExcelImportImpl) GetRowValue(rowIndex int) ([]string, error) {
	if rowIndex <= 0 {
		return nil, fmt.Errorf("index must be greater than 0")
	}

	rows, err := s.File.GetRows(s.activeSheetName)
	if err != nil {
		return nil, err
	}

	if rowIndex >= len(rows) {
		return nil, fmt.Errorf("rowIndex out of range")
	}

	return rows[rowIndex-1], nil
}

func (s *ExcelImportImpl) GetAllRowValue() ([][]string, error) {
	rows, err := s.File.GetRows(s.activeSheetName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *ExcelImportImpl) RowCount() int {
	rows, err := s.File.GetRows(s.activeSheetName)
	if err != nil {
		return 0
	}

	return len(rows)
}
