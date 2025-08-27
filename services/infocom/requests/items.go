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

	rows, err := table.Rows("Лист_1")
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

	items := make(map[string]int, 200)
	manufacturerFound := false

	for _, row := range *rawItems {
		switch len(row) {
		case 1:
			{
				lower := strings.ToLower(row[0])
				if strings.Contains(lower, "dahua") {
					manufacturerFound = true
				} else {
					manufacturerFound = false
				}
				break
			}
		case 2:
			{
				if manufacturerFound {
					quantity, err := parseQuantity(row[1])
					if err != nil {
						fmt.Printf("ошибка при парсинге из таблицы infocom: %s\n", err)
						continue
					}
					items[row[0]] = quantity
				}
				break
			}
		}
	}

	return &items, nil
}
