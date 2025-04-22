package db

import (
	"MerlionScript/utils/db/typesDB"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DBInstance struct {
	db *sql.DB
}

var instance *DBInstance

func CreateDirectory(path string) error {
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
	err := CreateDirectory("data")
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

func (instance *DBInstance) Init() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS codes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		ms_own_id INTEGER,
		moy_sklad TEXT NOT NULL,
		manufacturer TEXT NOT NULL UNIQUE,
        merlion TEXT NOT NULL
    );`

	_, err := instance.db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (instance *DBInstance) AddCodeRecord(record *typesDB.Codes) (bool, error) {
	exists, err := instance.checkRecordExists(record.Manufacturer)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	insertCodeRecordSQL := `INSERT INTO codes (ms_own_id, moy_sklad, manufacturer, merlion) VALUES (?, ?, ?, ?) ON CONFLICT(manufacturer) DO NOTHING`
	statement, err := instance.db.Prepare(insertCodeRecordSQL)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	res, err := statement.Exec(record.MsOwnId, record.MoySklad, record.Manufacturer, record.Merlion)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (instance *DBInstance) checkRecordExists(manufacturer string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM codes WHERE manufacturer = ?)"
	err := instance.db.QueryRow(query, manufacturer).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (instance *DBInstance) DeleteCodeRecord(id int) error {
	deleteCodeRecordSQL := `DELETE FROM codes WHERE id = ?`
	_, err := instance.db.Exec(deleteCodeRecordSQL, id)
	return err
}

func (instance *DBInstance) EditCodeRecord(record *typesDB.Codes) error {
	editCodeRecordSQL := `UPDATE codes SET ms_own_id = ?, moy_sklad = ?, manufacturer = ?, merlion = ? WHERE id = ?`
	statement, err := instance.db.Prepare(editCodeRecordSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(record.MsOwnId, record.MoySklad, record.Manufacturer, record.Merlion, record.Id)
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

func (instance *DBInstance) GetCodeRecordByMS(msCode string) (typesDB.Codes, bool, error) {
	GetCodeRecordSQL := `SELECT id, manufacturer, merlion FROM codes WHERE moy_sklad = ?`

	record := typesDB.Codes{MoySklad: msCode}
	err := instance.db.QueryRow(GetCodeRecordSQL, msCode).Scan(&record.Id, &record.Manufacturer, &record.Merlion)
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
		var MsOwnId sql.NullInt64
		rows.Scan(&record.Id, &MsOwnId, &record.MoySklad, &record.Manufacturer, &record.Merlion)
		if MsOwnId.Valid {
			record.MsOwnId = MsOwnId.Int64
		}
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) GetCodeRecordsNoMS() (*[]typesDB.Codes, error) {
	var records []typesDB.Codes
	rows, err := instance.db.Query("SELECT * FROM codes WHERE moy_sklad = ''")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record typesDB.Codes
		var MsOwnId sql.NullInt64
		rows.Scan(&record.Id, &MsOwnId, &record.MoySklad, &record.Manufacturer, &record.Merlion)
		if MsOwnId.Valid {
			record.MsOwnId = MsOwnId.Int64
		}
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) GetCodeRecordsFilledMS() (*[]typesDB.Codes, error) {
	var records []typesDB.Codes
	rows, err := instance.db.Query("SELECT * FROM codes WHERE moy_sklad != '' AND manufacturer != '' AND merlion != ''")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record typesDB.Codes
		var MsOwnId sql.NullInt64
		rows.Scan(&record.Id, &MsOwnId, &record.MoySklad, &record.Manufacturer, &record.Merlion)
		if MsOwnId.Valid {
			record.MsOwnId = MsOwnId.Int64
		}
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) CheckIfExistRecord(manufacturer string) (bool, error) {
	existsSQL := `SELECT 1 FROM codes WHERE manufacturer = ?`
	var exists int = 0
	err := instance.db.QueryRow(existsSQL, manufacturer).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (instance *DBInstance) GetLastOwnIdMS() (int64, error) {
	var maxValue sql.NullInt64

	err := instance.db.QueryRow("SELECT MAX(ms_own_id) FROM codes").Scan(&maxValue)
	if err != nil {
		return -1, err
	}

	// Если maxValue не имеет значения, устанавливаем его равным 1
	if !maxValue.Valid {
		maxValue.Int64 = 1
	}

	return maxValue.Int64, nil
}
