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
	initSQL :=
		`CREATE TABLE IF NOT EXISTS codes_merlion (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		ms_own_id INTEGER,
		moy_sklad TEXT NOT NULL,
		manufacturer TEXT NOT NULL UNIQUE,
        merlion TEXT NOT NULL,
		manufacturer_name TEXT NOT NULL DEFAULT "",
		try_upload_image INTEGER DEFAULT 0,
		uploaded_image INTEGER DEFAULT 0
    );

	CREATE TABLE IF NOT EXISTS codes_netlab (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		ms_own_id INTEGER,
		moy_sklad TEXT NOT NULL,
		manufacturer TEXT NOT NULL UNIQUE,
        netlab TEXT NOT NULL,
		manufacturer_name TEXT NOT NULL DEFAULT "",
		try_upload_image INTEGER DEFAULT 0,
		uploaded_image INTEGER DEFAULT 0
    );

	CREATE TABLE IF NOT EXISTS codes_ids (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		source TEXT NOT NULL,
		manufacturer_name TEXT NOT NULL,
		ms_own_id INTEGER NOT NULL
    );

	CREATE TRIGGER IF NOT EXISTS trg_insert_merlion 
	AFTER INSERT ON codes_merlion
	BEGIN
		INSERT OR IGNORE INTO codes_ids (source, ms_own_id, manufacturer_name)
		SELECT 'codes_merlion', NEW.ms_own_id, NEW.manufacturer_name
		WHERE NEW.ms_own_id IS NOT NULL AND NEW.ms_own_id <> 0;
	END;

	CREATE TRIGGER IF NOT EXISTS after_update_merlion
	AFTER UPDATE ON codes_merlion
	FOR EACH ROW
	BEGIN
			INSERT OR IGNORE INTO codes_ids (source, ms_own_id, manufacturer_name)
			SELECT 'codes_merlion', NEW.ms_own_id, NEW.manufacturer_name
			WHERE NEW.ms_own_id IS NOT NULL AND NEW.ms_own_id <> 0;
	END;

	CREATE TRIGGER IF NOT EXISTS trg_delete_merlion 
	AFTER DELETE ON codes_merlion
	BEGIN
		DELETE FROM codes_ids WHERE ms_own_id = OLD.id;
	END;

	CREATE TRIGGER IF NOT EXISTS trg_insert_netlab 
	AFTER INSERT ON codes_netlab
	BEGIN
		INSERT OR IGNORE INTO codes_ids (source, ms_own_id, manufacturer_name)
		SELECT 'codes_netlab', NEW.ms_own_id, NEW.manufacturer_name
		WHERE NEW.ms_own_id IS NOT NULL AND NEW.ms_own_id <> 0;
	END;

	CREATE TRIGGER IF NOT EXISTS after_update_netlab
	AFTER UPDATE ON codes_netlab
	FOR EACH ROW
	BEGIN
			INSERT OR IGNORE INTO codes_ids (source, ms_own_id, manufacturer_name)
			SELECT 'codes_netlab', NEW.ms_own_id, NEW.manufacturer_name
			WHERE NEW.ms_own_id IS NOT NULL AND NEW.ms_own_id <> 0;
	END;

	CREATE TRIGGER IF NOT EXISTS trg_delete_netlab 
	AFTER DELETE ON codes_netlab
	BEGIN
		DELETE FROM codes_ids WHERE ms_own_id = OLD.id;
	END;
	`

	_, err := instance.db.Exec(initSQL)
	if err != nil {
		return err
	}

	return nil
}

func (instance *DBInstance) AddCodeRecord(record *typesDB.Codes, tableName string) (bool, error) {
	exists, err := instance.checkRecordExists(record.Manufacturer, tableName)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	service := getServiceByTableName(tableName)
	if service == "" {
		return false, fmt.Errorf("неправильное название таблицы")
	}
	//insertCodeRecordSQL := `INSERT INTO ? (ms_own_id, moy_sklad, manufacturer, ?) VALUES (?, ?, ?, ?) ON CONFLICT(manufacturer) DO NOTHING`
	query := fmt.Sprintf("INSERT INTO %s (ms_own_id, moy_sklad, manufacturer, manufacturer_name, %s) VALUES (?, ?, ?, ?, ?) ON CONFLICT(manufacturer) DO NOTHING",
		tableName, service)
	statement, err := instance.db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer statement.Close()

	res, err := statement.Exec(record.MsOwnId, record.MoySklad, record.Manufacturer, record.ManufacturerName, record.Service)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (instance *DBInstance) GetCodeRecordByManufacturerAll(manufacturer string) (typesDB.Codes, bool, error) {
	tableNames := []string{typesDB.MerlionTable, typesDB.NetlabTable}
	for _, tn := range tableNames {
		record, exists, err := instance.GetCodeRecordByManufacturer(manufacturer, tn)
		if err != nil {
			return typesDB.Codes{}, false, err
		}
		if exists {
			return record, true, nil
		}
	}
	return typesDB.Codes{}, false, nil
}

func (instance *DBInstance) checkRecordExists(manufacturer string, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE manufacturer = ?)", tableName)
	err := instance.db.QueryRow(query, manufacturer).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (instance *DBInstance) DeleteCodeRecord(id int, tableName string) error {
	//deleteCodeRecordSQL := `DELETE FROM ? WHERE id = ?`
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := instance.db.Exec(query, id)
	return err
}

func (instance *DBInstance) EditCodeRecord(record *typesDB.Codes, tableName string) error {
	//editCodeRecordSQL := `UPDATE ? SET ms_own_id = ?, moy_sklad = ?, manufacturer = ?, ? = ? WHERE id = ?`
	service := getServiceByTableName(tableName)
	if service == "" {
		return fmt.Errorf("неправильное название таблицы")
	}
	query := fmt.Sprintf("UPDATE %s SET ms_own_id = ?, moy_sklad = ?, manufacturer = ?, manufacturer_name = ?, try_upload_image = ?, uploaded_image = ?,  %s = ? WHERE id = ?",
		tableName, service)
	statement, err := instance.db.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(record.MsOwnId, record.MoySklad, record.Manufacturer, record.ManufacturerName, record.TryLoadImage, record.LoadedImage, record.Service, record.Id)
	return err
}

func (instance *DBInstance) GetCodeRecordByManufacturer(manufacturer string, tableName string) (typesDB.Codes, bool, error) {
	//GetCodeRecordSQL := `SELECT id, moy_sklad, ? FROM ? WHERE manufacturer = ?`
	service := getServiceByTableName(tableName)
	if service == "" {
		return typesDB.Codes{}, false, fmt.Errorf("неправильное название таблицы")
	}
	query := fmt.Sprintf("SELECT id, moy_sklad, %s FROM %s WHERE manufacturer = ?", service, tableName)
	record := typesDB.Codes{Manufacturer: manufacturer}
	err := instance.db.QueryRow(query, manufacturer).Scan(&record.Id, &record.MoySklad, &record.Service)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.Codes{}, false, nil
		}
		return typesDB.Codes{}, false, err
	}
	return record, true, nil
}

func (instance *DBInstance) GetCodeRecordByMS(msCode string, tableName string) (typesDB.Codes, bool, error) {
	//GetCodeRecordSQL := `SELECT id, manufacturer, ? FROM ? WHERE moy_sklad = ?`
	service := getServiceByTableName(tableName)
	if service == "" {
		return typesDB.Codes{}, false, fmt.Errorf("неправильное название таблицы")
	}
	query := fmt.Sprintf("SELECT id, manufacturer, %s FROM %s WHERE moy_sklad = ?", service, tableName)
	record := typesDB.Codes{MoySklad: msCode}
	err := instance.db.QueryRow(query, msCode).Scan(&record.Id, &record.Manufacturer, &record.Service)
	if err != nil {
		if err == sql.ErrNoRows {
			return typesDB.Codes{}, false, nil
		}
		return typesDB.Codes{}, false, err
	}
	return record, true, nil
}

func (instance *DBInstance) GetCodeRecords(tableName string) (*[]typesDB.Codes, error) {
	var records []typesDB.Codes
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := instance.db.Query(query)
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
		rows.Scan(&record.Id, &MsOwnId, &record.MoySklad, &record.Manufacturer, &record.Service, &record.ManufacturerName, &record.TryLoadImage, &record.LoadedImage)
		if MsOwnId.Valid {
			record.MsOwnId = MsOwnId.Int64
		}
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) GetCodeRecordsNoMS(tableName string) (*[]typesDB.Codes, error) {
	var records []typesDB.Codes
	query := fmt.Sprintf("SELECT * FROM %s WHERE moy_sklad = ''", tableName)
	rows, err := instance.db.Query(query)
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
		rows.Scan(&record.Id, &MsOwnId, &record.MoySklad, &record.Manufacturer, &record.Service, &record.ManufacturerName, &record.TryLoadImage, &record.LoadedImage)
		if MsOwnId.Valid {
			record.MsOwnId = MsOwnId.Int64
		}
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) GetCodeRecordsFilledMS(tableName string) (*[]typesDB.Codes, error) {
	var records = make([]typesDB.Codes, 0, 500)
	service := getServiceByTableName(tableName)
	if service == "" {
		return nil, fmt.Errorf("неправильное название таблицы")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE moy_sklad != '' AND manufacturer != '' AND %s != ''", tableName, service)
	rows, err := instance.db.Query(query)
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
		rows.Scan(&record.Id, &MsOwnId, &record.MoySklad, &record.Manufacturer, &record.Service, &record.ManufacturerName, &record.TryLoadImage, &record.LoadedImage)
		if MsOwnId.Valid {
			record.MsOwnId = MsOwnId.Int64
		}
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) GetCodeRecordsFilledMSWithNoImage(tableName string) (*[]typesDB.Codes, error) {
	var records = make([]typesDB.Codes, 0, 500)
	service := getServiceByTableName(tableName)
	if service == "" {
		return nil, fmt.Errorf("неправильное название таблицы")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE moy_sklad != '' AND manufacturer != '' AND %s != '' AND try_upload_image == 0", tableName, service)
	rows, err := instance.db.Query(query)
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
		rows.Scan(&record.Id, &MsOwnId, &record.MoySklad, &record.Manufacturer, &record.Service, &record.ManufacturerName, &record.TryLoadImage, &record.LoadedImage)
		if MsOwnId.Valid {
			record.MsOwnId = MsOwnId.Int64
		}
		records = append(records, record)
	}
	return &records, nil
}

func (instance *DBInstance) CheckIfExistRecord(manufacturer string, tableName string) (bool, error) {
	//existsSQL := `SELECT 1 FROM ? WHERE manufacturer = ?`
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE manufacturer = ?", tableName)
	var exists int = 0
	err := instance.db.QueryRow(query, manufacturer).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (instance *DBInstance) GetLastOwnIdMS(tableName string) (int64, error) {
	var maxValue sql.NullInt64
	query := fmt.Sprintf("SELECT MAX(ms_own_id) FROM %s", tableName)

	err := instance.db.QueryRow(query).Scan(&maxValue)
	if err != nil {
		return -1, err
	}

	// Если maxValue не имеет значения, устанавливаем его равным 1
	if !maxValue.Valid {
		maxValue.Int64 = 1
	}

	return maxValue.Int64, nil
}

func getServiceByTableName(tableName string) string {
	switch tableName {
	case typesDB.MerlionTable:
		return typesDB.MerlionService
	case typesDB.NetlabTable:
		return typesDB.NetlabService
	default:
		return ""
	}
}
