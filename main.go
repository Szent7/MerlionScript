package main

import (
	"MerlionScript/controller"
	"MerlionScript/keeper"
	csvInstance "MerlionScript/utils/csv"
	"MerlionScript/utils/db"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	dbInstance, err := db.GetDBInstance()
	defer db.CloseDB()
	if err != nil {
		log.Fatalf("Ошибка при создании экземпляра базы данных: %s", err)
	}
	err = dbInstance.Init()
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %s", err)
	}
	importENV()
	importCSV()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	fmt.Println("Боже, Царя храни!")
	wg.Add(1)
	go controller.StartController(ctx, &wg)

	select {
	case <-sigCh:
		fmt.Println("Получен сигнал завершения приложения")
	case <-ctx.Done():
		fmt.Println("Контекст был отменен")
	}

	// Отмечаем, что контекст будет отменен и ждем завершения всех горутин
	cancel()
	wg.Wait()

	fmt.Println("Приложение корректно завершило работу")
}

func importENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка при парсинге .env файла: %s", err)
	}

	//SkladCredentials := os.Getenv("MOY_SKLAD_CREDENTIALS")
	SkladToken := os.Getenv("MOY_SKLAD_TOKEN")
	MerlionCredentials := os.Getenv("MERLION_CREDENTIALS")
	SkladName := os.Getenv("SKLAD")
	OrgName := os.Getenv("ORGANIZATION")
	CatName := os.Getenv("CATALOG")

	if SkladToken == "" || MerlionCredentials == "" || SkladName == "" || OrgName == "" {
		log.Fatalf("Данные для входа не обнаружены")
	}
	keeper.K.SetData(SkladToken, MerlionCredentials, SkladName, OrgName, CatName)
}

func importCSV() {
	csv, err := csvInstance.GetCSVInstance()
	if err != nil {
		log.Printf("Ошибка при открытии codes.csv: %s", err)
		return
	}
	if csv == nil {
		fmt.Println("Файл для импорта не обнаружен (codes.csv)")
		return
	}
	err = csv.ImportCodes()
	if err != nil {
		log.Printf("Ошибка при импорте codes.csv: %s", err)
	}
	csvInstance.CloseCSV()
}
