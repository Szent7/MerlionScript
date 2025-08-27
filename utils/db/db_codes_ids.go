package db

import (
	"MerlionScript/utils/db/typesDB"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

// InsertCodesIDs вставляет новую запись в таблицу, если записи с таким article не существует
func (instance *DB) InsertCodesIDs(record typesDB.CodesIDs) (bool, error) {
	// Проверка существования записи в таблице
	exists, err := instance.checkRecordExists(record.Article, typesDB.IDsTable)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	// Формирование запроса
	query := fmt.Sprintf("INSERT INTO %s (ms_own_id, moy_sklad, article, manufacturer) VALUES (?, ?, ?, ?) ON CONFLICT(article) DO NOTHING",
		typesDB.IDsTable)
	statement, err := instance.Prepare(query)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	// Выполнение запроса
	res, err := statement.Exec(record.MsOwnId, record.MoySkladCode, record.Article, record.Manufacturer)
	if err != nil {
		return false, err
	}

	// Проверка вставки новой записи
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

// GetCodesIDs возвращает запись по article из таблицы
func (instance *DB) GetCodesIDs(article string) (typesDB.CodesIDs, error) {
	// Формирование запроса
	query := fmt.Sprintf("SELECT id, ms_own_id, moy_sklad, manufacturer FROM %s WHERE article = ?", typesDB.IDsTable)

	// Выполнение запроса
	record := typesDB.CodesIDs{Article: article}
	err := instance.QueryRow(query, article).Scan(&record.Id, &record.MsOwnId, &record.MoySkladCode, &record.Manufacturer)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.CodesIDs{}, nil
		}
		return typesDB.CodesIDs{}, err
	}

	return record, nil
}

// GetCodesIDsByMS возвращает запись по коду ERP из таблицы
func (instance *DB) GetCodesIDsByMS(msCode string) (*[]typesDB.CodesIDs, error) {
	// Формирование запроса
	query := fmt.Sprintf("SELECT * FROM %s WHERE moy_sklad = ?", typesDB.IDsTable)

	// Выполнение запроса
	rows, err := instance.Query(query, msCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	// Извлечение результата запроса
	var records []typesDB.CodesIDs
	for rows.Next() {
		var record typesDB.CodesIDs
		rows.Scan(&record.Id, &record.MsOwnId, &record.MoySkladCode, &record.Article, &record.Manufacturer)
		records = append(records, record)
	}
	return &records, nil
}

// GetCodesIDsFilledMS возвращает список записей, у которых указан код ERP системы
func (instance *DB) GetCodesIDsFilledMS() (*[]typesDB.CodesIDs, error) {
	// Формирование запроса
	query := fmt.Sprintf("SELECT * FROM %s WHERE moy_sklad != '')", typesDB.IDsTable)

	// Выполнение запроса
	rows, err := instance.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	// Извлечение результата запроса
	var records []typesDB.CodesIDs
	for rows.Next() {
		var record typesDB.CodesIDs
		rows.Scan(&record.Id, &record.MsOwnId, &record.MoySkladCode, &record.Article, &record.Manufacturer)
		records = append(records, record)
	}
	return &records, nil
}

// GetCodesFilledMS возвращает список записей из главной таблицы и соответствие из таблицы сервиса, у которых указан код ERP системы
func (instance *DB) GetCodesFilledMS(tableName string) (*[]typesDB.General, error) {
	// Формирование запроса
	query := fmt.Sprintf(`
	SELECT %s.*, %s.*
	FROM %s 
	INNER JOIN %s ON %s.article = %s.article
	WHERE %s.moy_sklad != ''`, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName, typesDB.IDsTable)

	// Выполнение запроса
	rows, err := instance.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	// Извлечение результата запроса
	var records []typesDB.General
	for rows.Next() {
		var record typesDB.General
		rows.Scan(&record.Codes.Id, &record.Codes.MsOwnId, &record.Codes.MoySkladCode, &record.Codes.Article, &record.Codes.Manufacturer,
			&record.Service.Id, &record.Service.Article, &record.Service.ServiceCode, &record.Service.TryUploadImage)
		records = append(records, record)
	}
	return &records, nil
}

// GetCodesFilledMSNoImage возвращает список записей из главной таблицы и соответствие из таблицы сервиса, у которых указан код ERP системы и не загружены изображения
func (instance *DB) GetCodesFilledMSNoImage(tableName string) (*[]typesDB.General, error) {
	// Формирование запроса
	query := fmt.Sprintf(`
	SELECT %s.*, %s.*
	FROM %s 
	INNER JOIN %s ON %s.article = %s.article
	WHERE %s.moy_sklad != '' AND %s.try_upload_image = 0`, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName)

	// Выполнение запроса
	rows, err := instance.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	// Извлечение результата запроса
	var records []typesDB.General
	for rows.Next() {
		var record typesDB.General
		rows.Scan(&record.Codes.Id, &record.Codes.MsOwnId, &record.Codes.MoySkladCode, &record.Codes.Article, &record.Codes.Manufacturer,
			&record.Service.Id, &record.Service.Article, &record.Service.ServiceCode, &record.Service.TryUploadImage)
		records = append(records, record)
	}
	return &records, nil
}

// UpdateCodesIDs обновляет запись в главной таблице
func (instance *DB) UpdateCodesIDs(record typesDB.CodesIDs) error {
	// Формирование запроса
	query := fmt.Sprintf("UPDATE %s SET ms_own_id = ?, moy_sklad = ?, article = ?, manufacturer = ? WHERE id = ?",
		typesDB.IDsTable)
	statement, err := instance.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	// Выполнение запроса
	_, err = statement.Exec(record.MsOwnId, record.MoySkladCode, record.Article, record.Manufacturer, record.Id)
	return err
}

// GetLastOwnIdMS возвращает последний идентификатор созданных позиций
func (instance *DB) GetLastOwnIdMS() (int64, error) {
	// Формирование запроса
	query := fmt.Sprintf("SELECT MAX(ms_own_id) FROM %s", typesDB.IDsTable)

	// Выполнение запроса
	var maxValue sql.NullInt64
	err := instance.QueryRow(query).Scan(&maxValue)
	if err != nil {
		return -1, err
	}

	// Если maxValue не имеет значения, устанавливаем его равным 1
	if !maxValue.Valid {
		maxValue.Int64 = 1
	}

	return maxValue.Int64, nil
}

// ExtractCounterFromOwnID извлекает числовую часть из ID формата I00000
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
