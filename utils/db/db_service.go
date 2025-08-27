package db

import (
	"MerlionScript/services/common"
	"MerlionScript/utils/db/typesDB"
	"database/sql"
	"fmt"
)

// InsertService вставляет новую запись в таблицу, если записи с таким article не существует
func (instance *DB) InsertService(record typesDB.CodesService, tableName string) (bool, error) {
	// Проверка существования записи в таблице
	exists, err := instance.checkRecordExists(record.Article, tableName)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	// Формирование запроса
	query := fmt.Sprintf("INSERT INTO %s (article, service, try_upload_image) VALUES (?, ?, ?) ON CONFLICT(article) DO NOTHING",
		tableName)
	statement, err := instance.Prepare(query)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	// Выполнение запроса
	res, err := statement.Exec(record.Article, record.ServiceCode, record.TryUploadImage)
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

// GetService возвращает запись по article из таблицы
func (instance *DB) GetService(article string, tableName string) (typesDB.CodesService, error) {
	// Формирование запроса
	query := fmt.Sprintf("SELECT * FROM %s WHERE article = ?", tableName)

	// Выполнение запроса
	record := typesDB.CodesService{}
	err := instance.QueryRow(query, article).Scan(&record.Id, &record.Article, &record.ServiceCode, &record.TryUploadImage)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.CodesService{}, nil
		}
		return typesDB.CodesService{}, err
	}
	return record, nil
}

// GetServiceIDsNoMS возвращает список записей сервиса, у которых отсутствует код ERP системы
func (instance *DB) GetServiceIDsNoMS(tableName string) (*[]typesDB.CodesService, error) {
	// Формирование запроса
	query := fmt.Sprintf("SELECT * FROM %s WHERE article IN (SELECT article FROM codes_ids WHERE moy_sklad = '')", tableName)

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
	var records []typesDB.CodesService
	for rows.Next() {
		var record typesDB.CodesService
		rows.Scan(&record.Id, &record.Article, &record.ServiceCode, &record.TryUploadImage)
		records = append(records, record)
	}
	return &records, nil
}

// GetServiceIDsFilledMS возвращает список записей сервиса, у которых указан код ERP системы
func (instance *DB) GetServiceIDsFilledMS(tableName string) (*[]typesDB.CodesService, error) {
	// Формирование запроса
	query := fmt.Sprintf("SELECT * FROM %s WHERE article IN (SELECT article FROM codes_ids WHERE moy_sklad != '')", tableName)

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
	var records []typesDB.CodesService
	for rows.Next() {
		var record typesDB.CodesService
		rows.Scan(&record.Id, &record.Article, &record.ServiceCode, &record.TryUploadImage)
		records = append(records, record)
	}
	return &records, nil
}

// GetServiceIDsFilledMSWithNoImage возвращает список записей сервиса, у которых указан код ERP системы и не загружены изображения
func (instance *DB) GetServiceIDsFilledMSWithNoImage(tableName string) (*[]typesDB.CodesService, error) {
	// Формирование запроса
	query := fmt.Sprintf("SELECT * FROM %s WHERE article IN (SELECT article FROM codes_ids WHERE moy_sklad != '') AND try_upload_image = 0", tableName)

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
	var records []typesDB.CodesService
	for rows.Next() {
		var record typesDB.CodesService
		rows.Scan(&record.Id, &record.Article, &record.ServiceCode, &record.TryUploadImage)
		records = append(records, record)
	}
	return &records, nil
}

// GetAllServicesByArticle возвращает список записей всех сервисов, у которых один и тот же артикул
func (instance *DB) GetAllServicesByArticle(article string) (*[]typesDB.CodesService, error) {
	var records []typesDB.CodesService
	for _, v := range common.RegisteredServices { // Цикл по зарегистрированным сервисам
		record, err := instance.GetService(article, v.DBTableName) // Получение записей из указанного сервиса
		if err != nil {
			return nil, err
		}
		if record != (typesDB.CodesService{}) {
			records = append(records, record)
		}
	}
	return &records, nil
}

// UpdateService обновляет запись в таблице сервиса
func (instance *DB) UpdateService(record typesDB.CodesService, tableName string) error {
	// Формирование запроса
	query := fmt.Sprintf("UPDATE %s SET article = ?, service = ?, try_upload_image = ? WHERE id = ?",
		tableName)
	statement, err := instance.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	// Выполнение запроса
	_, err = statement.Exec(record.Article, record.ServiceCode, record.TryUploadImage, record.Id)
	return err
}
