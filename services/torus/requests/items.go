package requests

import (
	"MerlionScript/utils/excel"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var quantityRe = regexp.MustCompile(`\d+`)

// Функция для парсинга количества
func parseQuantity(s string) (int, error) {
	// Ищем первое вхождение числа в строке
	match := quantityRe.FindStringSubmatch(s)
	if len(match) == 0 {
		return 0, fmt.Errorf("нет числа в строке: %s", s)
	}

	// Конвертируем найденное число в int
	return strconv.Atoi(match[0])
}

func GetItems(table *excel.Workbook) (*[][]string, error) {
	if table == nil {
		return nil, nil
	}

	rows, err := table.Rows("TDSheet")
	if err != nil {
		return nil, err
	}

	return &rows, nil
}

// фильтр по производителю
func GetItemsFormatted(table *excel.Workbook) (*map[string]int, error) {
	rawItems, err := GetItems(table)
	if err != nil {
		return nil, err
	}

	items := make(map[string]int, 600)
	manufacturerFound := false

	for _, row := range *rawItems {
		rowLen := len(row)
		switch {
		case rowLen == 1:
			{
				lower := strings.ToLower(row[0])
				if strings.Contains(lower, "dahua") {
					manufacturerFound = true
				} else {
					manufacturerFound = false
				}
				break
			}
		case (4 <= rowLen && rowLen < 7):
			{
				if manufacturerFound {
					items[row[3]] = 0
				}
				break
			}
		case rowLen >= 7:
			{
				if manufacturerFound {
					if row[6] == "" {
						items[row[3]] = 0
					} else {
						quantity, err := parseQuantity(row[6])
						if err != nil {
							fmt.Printf("ошибка при парсинге из таблицы torus: %s\n", err)
							continue
						}
						items[row[3]] = quantity
					}
				}
				break
			}
		}
	}

	return &items, nil
}
