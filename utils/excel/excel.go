package excel

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Workbook struct {
	*excelize.File
	Path string
}

// Open открывает существующий файл Excel и возвращает Workbook
func Open(path string) (*Workbook, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	return &Workbook{File: f, Path: path}, nil
}

// Close закрывает файл
func (wb *Workbook) Close() error {
	return wb.Close()
}

// Save сохраняет файл в оригинальное место
func (wb *Workbook) Save() error {
	return wb.Save()
}

// SaveAs сохраняет файл в другой путь
func (wb *Workbook) SaveAs(newPath string) error {
	return wb.SaveAs(newPath)
}

// ---------------------------------------
// Ячейки
// ---------------------------------------

// CellValue читаёт значение из ячейки ("Sheet1", "A1") и возвращает его как строку
func (wb *Workbook) CellValue(sheet, cell string) (string, error) {
	val, err := wb.GetCellValue(sheet, cell)
	if err != nil {
		return "", err
	}
	return val, nil
}

// SetCellValue записывает значение в ячейку
func (wb *Workbook) SetCellValue(sheet, cell string, value any) error {
	if err := wb.SetCellValue(sheet, cell, value); err != nil {
		return err
	}
	return nil
}

// ---------------------------------------
// Строки/столбцы
// ---------------------------------------

// Rows возвращает все строки указанного листа
func (wb *Workbook) Rows(sheet string) ([][]string, error) {
	rows, err := wb.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("получить строки из %s: %w", sheet, err)
	}
	return rows, nil
}

// Row возвращает n‑ую строку. Если ряд не существует – возвращает nil
func (wb *Workbook) Row(sheet string, rowIndex int) ([]string, error) {
	if rowIndex < 1 {
		return nil, errors.New("индекс строки должен быть >= 1")
	}
	rows, err := wb.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	if rowIndex > len(rows) {
		return nil, nil
	}
	return rows[rowIndex-1], nil
}

// Column возвращает все значения из указанного столбца (буквенно, например "B")
func (wb *Workbook) Column(sheet, column string) ([]string, error) {
	if column == "" {
		return nil, errors.New("имя столбца не может быть пустым")
	}
	rows, err := wb.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	var res []string
	pos, err := excelize.ColumnNameToNumber(column)
	if err != nil {
		return nil, fmt.Errorf("неверный столбец %q: %w", column, err)
	}
	for _, r := range rows {
		if pos-1 < len(r) {
			res = append(res, r[pos-1])
		} else {
			res = append(res, "")
		}
	}
	return res, nil
}

// ---------------------------------------
// Строки
// ---------------------------------------

// AppendRows добавляет строки в конец листа
func (wb *Workbook) AppendRow(sheet string, row []any) error {
	// Проверяем, существует ли лист
	if _, err := wb.GetSheetIndex(sheet); err != nil {
		return err
	}

	// Получаем количество строк в листе
	rows, err := wb.GetRows(sheet)
	if err != nil {
		return err
	}
	rowsCount := len(rows)
	if rowsCount == 0 {
		rowsCount = 1
	} else {
		rowsCount++
	}

	// Добавляем строку
	if err := wb.SetSheetRow(sheet, "A"+strconv.Itoa(rowsCount), &row); err != nil {
		return err
	}

	return nil
}

// InsertRows вставляет строки начиная с rowIndex
func (wb *Workbook) InsertRows(sheet string, rowIndex int, row []any) error {
	// Проверяем, существует ли лист
	if _, err := wb.GetSheetIndex(sheet); err != nil {
		return err
	}

	// Получаем количество строк в листе
	rows, err := wb.GetRows(sheet)
	if err != nil {
		return err
	}
	rowsCount := len(rows)
	if rowsCount == 0 {
		rowsCount = 1
	} else {
		rowsCount++
	}

	// Добавляем строку
	if err := wb.SetSheetRow(sheet, "A"+strconv.Itoa(rowsCount), &row); err != nil {
		return err
	}

	return nil
}

// ---------------------------------------
// Лист
// ---------------------------------------

// CreateSheet создаёт новый лист
func (wb *Workbook) CreateSheet(name string) error {
	if index, _ := wb.GetSheetIndex(name); index > 0 {
		return nil
	}
	if _, err := wb.NewSheet(name); err != nil {
		return err
	}
	return nil
}

// DeleteSheet удаляет лист
func (wb *Workbook) DeleteSheet(name string) error {
	if index, _ := wb.GetSheetIndex(name); index <= 0 {
		return nil
	}
	if err := wb.DeleteSheet(name); err != nil {
		return err
	}
	return nil
}
