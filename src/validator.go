package src

import (
	"MerlionScript/services/common"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//Проверяет артикул из сервиса и результатами поиска в ERP
//PositionERP.Article - в МС это не артикул товара, а код
func CompareArticle(PositionsERP *[]common.ItemList, PositionService *common.ItemList) (common.ItemList, error) {
	article := ""
	var foundedPosition common.ItemList
	var founded bool = false
	var dhiProblem bool = false
	var emptyArticleProblem bool = false
	var sProblem bool = false
	var serviceSnum int = -10
	var erpSnum int = -11
	var bProblem bool = false
	var serviceBnum int = -10
	var erpBnum int = -11
	// Вытягиваем -S0/-0000B номер из названия позиции сервиса, если есть
	substringsService := strings.Fields(PositionService.PositionName)
	for _, subS := range substringsService {
		serviceSnum = ExtractNumberFromS(subS)
		if serviceSnum >= 0 {
			break
		}
		serviceBnum = ExtractNumberFromB(subS)
		if serviceBnum >= 0 {
			break
		}
	}
	if serviceSnum == -1 {
		return common.ItemList{}, fmt.Errorf("ошибка при парсинге -S0 номера | article = %s\n", PositionService.Article)
	}
	if serviceBnum == -1 {
		return common.ItemList{}, fmt.Errorf("ошибка при парсинге -0000B номера | article = %s\n", PositionService.Article)
	}
	//TODO позиция с 3 артикулами, добавить проверку на уникальность кода мс
	for i := range *PositionsERP {
		// Если полностью совпадает
		if СontainsSubstring((*PositionsERP)[i].PositionName, PositionService.Article) {
			founded = true
			foundedPosition = (*PositionsERP)[i]
			// Но при этом артикул пустой
			if (*PositionsERP)[i].Article == "" {
				emptyArticleProblem = true
				break
			}
			// Не пустой
			article = (*PositionsERP)[i].Article
			break
		}
		// Если начинается на DH-/DHI-, то убираем их из сравнения
		erpManufacturer := IgnoreDHManufacturer((*PositionsERP)[i].PositionName)
		serviceManufacturer := IgnoreDHManufacturer(PositionService.Article)
		if СontainsSubstring(erpManufacturer, serviceManufacturer) {
			dhiProblem = true
			break
		}
		// Если заканчивается на -S0
		if CheckSManufacturer((*PositionsERP)[i].PositionName, PositionService.Article) && serviceSnum >= 0 {
			// Вытягиваем -S0 номер из мс, если есть
			substringsMS := strings.Fields((*PositionsERP)[i].PositionName)
			for _, subS := range substringsMS {
				erpSnum = ExtractNumberFromS(subS)
				if erpSnum >= 0 {
					break
				}
			}
			// Если номера совпали
			if erpSnum == serviceSnum {
				sProblem = true
				founded = true
				foundedPosition = (*PositionsERP)[i]
				// Но при этом артикул пустой
				if (*PositionsERP)[i].Article == "" {
					emptyArticleProblem = true
					break
				}
				// Не пустой
				article = (*PositionsERP)[i].Article
				break
			}
		}

		// Если заканчивается на -0000B
		if CheckBManufacturer((*PositionsERP)[i].PositionName, PositionService.Article) && serviceBnum >= 0 {
			// Вытягиваем -0000B номер из мс, если есть
			substringsMS := strings.Fields((*PositionsERP)[i].PositionName)
			for _, subS := range substringsMS {
				erpBnum = ExtractNumberFromB(subS)
				if erpBnum >= 0 {
					break
				}
			}
			// Если номера совпали
			if erpBnum == serviceBnum {
				bProblem = true
				founded = true
				foundedPosition = (*PositionsERP)[i]
				// Но при этом артикул пустой
				if (*PositionsERP)[i].Article == "" {
					emptyArticleProblem = true
					break
				}
				// Не пустой
				article = (*PositionsERP)[i].Article
				break
			}
		}
	}

	if erpSnum == -1 {
		return common.ItemList{}, fmt.Errorf("ошибка при парсинге -S0 номера из ERP | article = %s\n", PositionService.Article)
	}
	if erpBnum == -1 {
		return common.ItemList{}, fmt.Errorf("ошибка при парсинге -0000B номера из ERP | article = %s\n", PositionService.Article)
	}
	if dhiProblem {
		return common.ItemList{}, fmt.Errorf("проблема DH-/DHI- | article = %s\n", PositionService.Article)
	}
	if emptyArticleProblem {
		return common.ItemList{}, fmt.Errorf("проблема пустого артикула | article = %s\n", PositionService.Article)
	}
	if sProblem {
		return common.ItemList{}, fmt.Errorf("проблема окончания -S0 | article = %s соответствие на мс = %s\n", PositionService.Article, article)
	}
	if bProblem {
		return common.ItemList{}, fmt.Errorf("проблема окончания -0000B | article = %s соответствие на мс = %s\n", PositionService.Article, article)
	}
	if !founded {
		return common.ItemList{}, nil
	}

	return foundedPosition, nil
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

// Проверка на постфикс -S0
func CheckSManufacturer(s string, substr string) bool {
	re := regexp.MustCompile(regexp.QuoteMeta(substr) + `-S\d`)
	matched := re.FindString(s)
	if matched != "" {
		return true
		//fmt.Println("Подстрока найдена:", matched)
	} else {
		return false
		//fmt.Println("Подстрока не найдена.")
	}
}

// Извлечение числа из постфикса -S0
// Возвращает -1 при ошибке парсинга
// Возвращает -2 если постфикса -S0 нет
func ExtractNumberFromS(s string) int {
	re := regexp.MustCompile(`-S(\d+)`)
	matches := re.FindStringSubmatch(s)

	if len(matches) > 1 {
		// Преобразование найденного числа в int
		number, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1
		}
		return number
	}

	return -2
}

// Проверка на постфикс -0000B
func CheckBManufacturer(s string, substr string) bool {
	re := regexp.MustCompile(regexp.QuoteMeta(substr) + `-\d{4}B`)
	matched := re.FindString(s)
	if matched != "" {
		return true
	} else {
		return false
	}
}

// Извлечение числа из постфикса -0000B
// Возвращает -1 при ошибке парсинга
// Возвращает -2 если постфикса -S0 нет
func ExtractNumberFromB(s string) int {
	re := regexp.MustCompile(`-(\d{4})B`)
	matches := re.FindStringSubmatch(s)

	if len(matches) > 1 {
		// Преобразование найденного числа в int
		number, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1
		}
		return number
	}

	return -2
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
