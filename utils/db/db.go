package db

import (
	"MerlionScript/utils/db/typesDB"
	"MerlionScript/utils/dir"
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Структура для работы с базой данных
type DB struct {
	*sql.DB
}

// Глобальная переменная для хранения текущего экземпляра подключения к базе данных
var instance *DB

// createDBConnection создает новое соединение с базой данных SQLite
func createDBConnection(pathToDB string) (*sql.DB, error) {
	err := dir.CreateDirectoryDefault(pathToDB) // Создает директорию для базы данных по указанному пути
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(pathToDB, "database.db") // Полный путь к файлу базы данных
	db, err := sql.Open("sqlite3", fullPath)           // Открывает соединение с базой данных
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetDB возвращает экземпляр базы данных. Паттерн "одиночка"
func GetDB(pathToDB string) (*DB, error) {
	if instance == nil {
		db, err := createDBConnection(pathToDB)
		if err != nil {
			return nil, err
		}
		instance = &DB{db}
	}
	return instance, nil
}

// CloseDB закрывает текущее соединение с базой данных
func CloseDB() {
	if instance != nil {
		err := instance.Close()
		if err != nil {
			panic(err)
		}
		instance = nil
	}
}

// initSQL формирует SQL-запрос для создания таблиц
func (instance *DB) initSQL(tableNames []string) string {
	tablesQuerry := ""
	Querry := fmt.Sprintf(typesDB.TableIDsSQL, typesDB.IDsTable) // Сначала создаем главную таблицу
	for _, table := range tableNames {                           // Затем добавляем таблицы сервисов
		tablesQuerry += fmt.Sprintf(typesDB.TableServiceSQL, table, typesDB.IDsTable)
	}
	Querry += tablesQuerry
	return Querry
}

// Init инициализирует базу данных, создавая необходимые таблицы
func (instance *DB) Init(tableNames []string) error {
	initSQL := instance.initSQL(tableNames) // Получаем SQL-запрос для создания таблиц

	_, err := instance.Exec(initSQL)
	if err != nil {
		return err
	}
	_, err = instance.Exec("PRAGMA foreign_keys = ON;") // Включаем поддержку внешних ключей
	if err != nil {
		return err
	}

	return nil
}

// checkRecordExists проверяет, существует ли запись в указанной таблице
func (instance *DB) checkRecordExists(article string, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE article = ?)", tableName)
	err := instance.QueryRow(query, article).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// DeleteCodesRecord удаляет запись из таблицы
func (instance *DB) DeleteCodesRecord(id int, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := instance.Exec(query, id)
	return err
}

// GetFormatID формирует строку с идентификатором в формате "Ixxxxx"
func GetFormatID(counter int64) string {
	newNumPart := counter
	newID := fmt.Sprintf("I%05d", newNumPart)
	return newID
}
