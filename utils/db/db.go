package db

import (
	"MerlionScript/utils/db/typesDB"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DBInstance struct {
	db *sql.DB
}

var instance *DBInstance

func CheckDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0660)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func createDBConnection() (*sql.DB, error) {
	err := CheckDirectory("data")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "./data/database.db")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDBInstance() (*DBInstance, error) {
	if instance == nil {
		db, err := createDBConnection()
		if err != nil {
			return nil, err
		}
		instance = &DBInstance{db: db}
	}
	return instance, nil
}

func CloseDB() {
	if instance != nil {
		err := instance.db.Close()
		if err != nil {
			panic(err)
		}
		instance = nil
	}
}

func (instance *DBInstance) Init() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS codes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		moy_sklad TEXT NOT NULL,
		manufacturer TEXT NOT NULL,
        merlion TEXT NOT NULL
    );`

	_, err := instance.db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func (instance *DBInstance) AddCodeRecord(record *typesDB.Codes) error {
	insertCodeRecordSQL := `INSERT INTO codes (moy_sklad, manufacturer, merlion) VALUES (?, ?, ?)`
	statement, err := instance.db.Prepare(insertCodeRecordSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(record.MoySklad, record.Manufacturer, record.Merlion)
	return err
}

func (instance *DBInstance) DeleteCodeRecord(id int) error {
	deleteCodeRecordSQL := `DELETE FROM codes WHERE id = ?`
	_, err := instance.db.Exec(deleteCodeRecordSQL, id)
	return err
}

func (instance *DBInstance) EditCodeRecord(record *typesDB.Codes) error {
	editCodeRecordSQL := `UPDATE codes SET moy_sklad = ?, manufacturer = ?, merlion = ? WHERE id = ?`
	statement, err := instance.db.Prepare(editCodeRecordSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(record.MoySklad, record.Manufacturer, record.Merlion, record.Id)
	return err
}

func (instance *DBInstance) GetCodeRecordByManufacturer(manufacturer string) (typesDB.Codes, bool, error) {
	GetCodeRecordSQL := `SELECT id, moy_sklad, merlion FROM codes WHERE manufacturer = ?`

	record := typesDB.Codes{Manufacturer: manufacturer}
	err := instance.db.QueryRow(GetCodeRecordSQL, manufacturer).Scan(&record.Id, &record.MoySklad, &record.Merlion)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.Codes{}, false, nil
		}
		return typesDB.Codes{}, false, err
	}
	return record, true, nil
}

func (instance *DBInstance) GetCodeRecords() (*[]typesDB.Codes, error) {
	var records []typesDB.Codes
	rows, err := instance.db.Query("SELECT * FROM codes")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record typesDB.Codes
		rows.Scan(&record.Id, &record.MoySklad, &record.Manufacturer, &record.Merlion)
		records = append(records, record)
	}
	return &records, nil
}
