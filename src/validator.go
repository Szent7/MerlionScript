package src

import (
	"MerlionScript/services/common"
	"fmt"
	"strings"
)

// Поле Article в PositionsERP - это не артикул, а номер для идентификации!!!
func CompareArticle(PositionsERP *[]common.ItemList, PositionService *common.ItemList) (common.ItemList, error) {
	// Берем артикул из названия
	articleService := findSubInSlice(PositionService.PositionName, PositionService.Article)

	// Если артикул в названии не найден
	if articleService == "" {
		articleService = PositionService.Article
	}

	for i := range *PositionsERP {
		// Если полностью совпадает
		if СontainsSubstring((*PositionsERP)[i].PositionName, articleService) {
			// Но при этом артикул пустой
			if (*PositionsERP)[i].Article == "" {
				return common.ItemList{}, fmt.Errorf("проблема пустого артикула | article = %s\n", articleService)
			}
			// Не пустой
			return (*PositionsERP)[i], nil // полное совпадение
		}
		// Если начинается на DH-/DHI-
		serviceArticle := IgnoreDHManufacturer(articleService)
		erpArticle := findSubInSlice((*PositionsERP)[i].PositionName, serviceArticle)
		if serviceArticle == IgnoreDHManufacturer(erpArticle) {
			return common.ItemList{}, fmt.Errorf("проблема DH-/DHI- | article = %s\n", articleService)
		}
	}

	return common.ItemList{}, nil
}

func findSubInSlice(name string, find string) string {
	substrings := strings.Fields(name)
	for _, subS := range substrings {
		if strings.Contains(subS, find) {
			return subS
		}
	}
	return ""
}

// Сравнение артикула в названии
func СontainsSubstring(s string, substr string) bool {
	n := len(substr)
	if n == 0 || s == "" || n > len(s) {
		return false
	}

	for i := 0; i <= len(s)-n; i++ {
		if strings.HasPrefix(s[i:], substr) {
			// Проверяем, что после подстроки идет либо пробел, либо конец строки
			if i+n == len(s) || s[i+n] == ' ' {
				return true
			}
		}
	}

	return false
}

// Игнорируем DH- и DHI-
func IgnoreDHManufacturer(manufacturer string) string {
	if strings.HasPrefix(manufacturer, "DH-") {
		return manufacturer[3:]
	} else if strings.HasPrefix(manufacturer, "DHI-") {
		return manufacturer[4:]
	}
	return manufacturer
}
