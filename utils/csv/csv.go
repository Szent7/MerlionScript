package csv

import (
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/typesDB"
	"encoding/csv"
	"os"
)

type CSVFile struct {
	file   *os.File
	writer *csv.Writer
	reader *csv.Reader
}

var instance *CSVFile

func GetCSVInstance() (*CSVFile, error) {
	if instance == nil {
		var err error
		instance, err = сreateCSV()
		if err != nil {
			return nil, err
		}
		/*if instance.initCSV() != nil {
			fmt.Printf("Error initCSV")
		}*/
	}
	return instance, nil
}

func CloseCSV() {
	instance.file.Close()
	if instance != nil {
		instance = nil
	}
}

func сreateCSV() (*CSVFile, error) {
	err := db.CheckDirectory("data")
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile("./data/codes.csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &CSVFile{file: file, reader: csv.NewReader(file), writer: csv.NewWriter(file)}, nil
}

func (cf *CSVFile) initCSV() error {
	return cf.WriteRecord([]string{"Moy Sklad", "Manufacturer code", "Merlion"})
}

func (cf *CSVFile) WriteRecord(record []string) error {
	defer cf.writer.Flush()
	return cf.writer.Write(record)
}

func (cf *CSVFile) ReadAllRecords() ([][]string, error) {
	cf.file.Seek(0, 0)
	return cf.reader.ReadAll()
}

func (cf *CSVFile) ReadSpecificRows(startRow int, numRows int) ([][]string, error) {
	records, err := cf.ReadAllRecords()
	if err != nil {
		return nil, err
	}
	if startRow+numRows > len(records) {
		numRows = len(records) - startRow
	}
	return records[startRow : startRow+numRows], nil
}

func (cf *CSVFile) ImportCodes() error {
	records, err := cf.ReadAllRecords()
	if err != nil {
		return err
	}
	dbInstance, err := db.GetDBInstance()
	if err != nil {
		return err
	}
	dbInstance.Init()
	for i := 1; i < len(records); i++ {
		err = dbInstance.AddCodeRecord(&typesDB.Codes{
			MoySklad:     records[i][0],
			Manufacturer: records[i][1],
			Merlion:      records[i][2],
		})
		if err != nil {
			return err
		}
	}
	return nil
}
