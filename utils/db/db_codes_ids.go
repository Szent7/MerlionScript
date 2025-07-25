package db

import (
	"MerlionScript/utils/db/typesDB"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

func (instance *DB) InsertCodesIDs(record typesDB.CodesIDs) (bool, error) {
	exists, err := instance.checkRecordExists(record.Article, typesDB.IDsTable)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	query := fmt.Sprintf("INSERT INTO %s (ms_own_id, moy_sklad, article, manufacturer) VALUES (?, ?, ?, ?) ON CONFLICT(article) DO NOTHING",
		typesDB.IDsTable)
	statement, err := instance.Prepare(query)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	res, err := statement.Exec(record.MsOwnId, record.MoySkladCode, record.Article, record.Manufacturer)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (instance *DB) GetCodesIDs(article string) (typesDB.CodesIDs, error) {
	query := fmt.Sprintf("SELECT (id, ms_own_id, moy_sklad, article, manufacturer) FROM %s WHERE article = ?", typesDB.IDsTable)
	record := typesDB.CodesIDs{}

	err := instance.QueryRow(query, article).Scan(&record.Id, &record.MsOwnId, &record.MoySkladCode, &record.Article, &record.Manufacturer)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.CodesIDs{}, nil
		}
		return typesDB.CodesIDs{}, err
	}
	return record, nil
}

func (instance *DB) GetCodesIDsByMS(msCode string) (*[]typesDB.CodesIDs, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE moy_sklad = ?", typesDB.IDsTable)
	var records []typesDB.CodesIDs

	rows, err := instance.Query(query, msCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record typesDB.CodesIDs
		rows.Scan(&record.Id, &record.MsOwnId, &record.MoySkladCode, &record.Article, &record.Manufacturer)
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DB) GetCodesIDsFilledMS() (*[]typesDB.CodesIDs, error) {
	var records []typesDB.CodesIDs
	query := fmt.Sprintf("SELECT * FROM %s WHERE moy_sklad != '')", typesDB.IDsTable)
	rows, err := instance.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record typesDB.CodesIDs
		rows.Scan(&record.Id, &record.MsOwnId, &record.MoySkladCode, &record.Article, &record.Manufacturer)
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DB) GetCodesFilledMS(tableName string) (*[]typesDB.General, error) {
	var records []typesDB.General
	query := fmt.Sprintf(`
	SELECT %s.*, %s.*
	FROM %s 
	INNER JOIN %s ON %s.article = %s.article
	WHERE %s.moy_sklad != '')`, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName, typesDB.IDsTable)
	rows, err := instance.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record typesDB.General
		rows.Scan(&record.Codes.Id, &record.Codes.MsOwnId, &record.Codes.MoySkladCode, &record.Codes.Article, &record.Codes.Manufacturer,
			&record.Service.Id, &record.Service.Article, &record.Service.ServiceCode, &record.Service.TryUploadImage)
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DB) GetCodesFilledMSNoImage(tableName string) (*[]typesDB.General, error) {
	var records []typesDB.General
	query := fmt.Sprintf(`
	SELECT %s.*, %s.*
	FROM %s 
	INNER JOIN %s ON %s.article = %s.article
	WHERE %s.moy_sklad != '' AND %s.try_upload_imgage = 0)`, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName, typesDB.IDsTable, tableName)
	rows, err := instance.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record typesDB.General
		rows.Scan(&record.Codes.Id, &record.Codes.MsOwnId, &record.Codes.MoySkladCode, &record.Codes.Article, &record.Codes.Manufacturer,
			&record.Service.Id, &record.Service.Article, &record.Service.ServiceCode, &record.Service.TryUploadImage)
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DB) UpdateCodesIDs(record typesDB.CodesIDs) error {
	query := fmt.Sprintf("UPDATE %s SET ms_own_id = ?, moy_sklad = ?, article = ?, manufacturer = ? WHERE id = ?",
		typesDB.IDsTable)

	statement, err := instance.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(record.MsOwnId, record.MoySkladCode, record.Article, record.Manufacturer, record.Id)
	return err
}

func (instance *DB) GetLastOwnIdMS() (int64, error) {
	var maxValue sql.NullInt64
	query := fmt.Sprintf("SELECT MAX(ms_own_id) FROM %s", typesDB.IDsTable)

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
