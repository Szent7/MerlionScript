package db

import (
	"MerlionScript/utils/db/typesDB"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

var instance *DB

func createDirectory(path string) error {
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

func createDBConnection(pathToDB string) (*sql.DB, error) {
	err := createDirectory(pathToDB)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(pathToDB, "database.db")
	db, err := sql.Open("sqlite3", fullPath)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err != nil {
		return nil, err
	}
	return db, nil
}

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

func CloseDB() {
	if instance != nil {
		err := instance.Close()
		if err != nil {
			panic(err)
		}
		instance = nil
	}
}

func (instance *DB) initSQL(tableNames []string) string {
	tablesQuerry := ""
	/*Querry := `CREATE TABLE IF NOT EXISTS "codes_ids" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"ms_own_id" INTEGER NOT NULL,
		"moy_sklad" TEXT NOT NULL,
		"article" TEXT NOT NULL UNIQUE,
		"manufacturer" TEXT NOT NULL
	);`*/
	Querry := fmt.Sprintf(typesDB.TableIDsSQL, typesDB.IDsTable)
	for _, table := range tableNames {
		/*tablesQuerry += fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"article" TEXT NOT NULL UNIQUE,
			"service" TEXT NOT NULL,
			"try_upload_image" INTEGER NOT NULL,
			FOREIGN KEY ("article") REFERENCES "codes_ids"("article")
			ON UPDATE CASCADE ON DELETE CASCADE
		);`, table)*/
		tablesQuerry += fmt.Sprintf(typesDB.TableServiceSQL, table, typesDB.IDsTable)
	}
	Querry += tablesQuerry
	return Querry
}

func (instance *DB) Init(tableNames []string) error {
	initSQL := instance.initSQL(tableNames)

	_, err := instance.Exec(initSQL)
	if err != nil {
		return err
	}
	_, err = instance.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	return nil
}

func (instance *DB) checkRecordExists(article string, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE article = ?)", tableName)
	err := instance.QueryRow(query, article).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (instance *DB) DeleteCodesRecord(id int, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := instance.Exec(query, id)
	return err
}

func GetFormatID(counter int64) string {
	newNumPart := counter
	newID := fmt.Sprintf("I%05d", newNumPart)
	return newID
}
