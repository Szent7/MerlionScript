package merlion

import (
	"MerlionScript/cache"
	skladReq "MerlionScript/services/sklad/requests"
	skladTypes "MerlionScript/types/restTypes/sklad"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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
		//fmt.Println("Подстрока найдена:", matched)
	} else {
		return false
		//fmt.Println("Подстрока не найдена.")
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

// Парсинг counter из артикула
func ExtractCounterFromOwnID(ownId string) (int64, error) {
	// Регулярное выражение для поиска числовой части в формате I00000
	re := regexp.MustCompile(`I(\d+)`)

	// Находим все совпадения
	matches := re.FindStringSubmatch(ownId)
	if len(matches) != 2 {
		return -1, fmt.Errorf("invalid code format")
	}

	// Извлекаем числовую часть из совпадений
	numStr := matches[1]

	// Преобразуем строку в int64
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to convert number: %w", err)
	}

	return num, nil
}

func GetItemAndCache(code string) (skladTypes.Rows, error) {
	itemMS, err := skladReq.GetItem(code)
	if err != nil {
		return skladTypes.Rows{}, err
	}

	if itemMS.Id == "" {
		return skladTypes.Rows{}, fmt.Errorf("пустой ID для кода: %s", code)
	}

	if err := cache.CacheRecord(code, itemMS); err != nil {
		return skladTypes.Rows{}, err
	}

	return itemMS, nil
}
