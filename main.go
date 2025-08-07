package main

import (
	"MerlionScript/controller"
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	"MerlionScript/services/common/initializer"
	"MerlionScript/utils/cache"
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/typesDB"
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
	ctx, cancel := context.WithCancel(context.Background())

	importENV()

	initializer.InitServices()

	dbInstance, err := db.GetDB(typesDB.PathDB)
	defer db.CloseDB()
	if err != nil {
		log.Fatalf("Ошибка при создании экземпляра базы данных: %s", err)
	}
	err = dbInstance.Init(common.GetTableNames())
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %s", err)
	}

	err = cache.InitCache(ctx)
	if err != nil {
		log.Fatalf("Ошибка при инициализации кэша: %s", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	fmt.Println("Боже, Царя храни!")
	wg.Add(1)
	go controller.StartController(ctx, &wg, dbInstance)

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
	var data = make(map[string]string, 13)
	data[keeper.MerlionCredentialsEnv] = os.Getenv(keeper.MerlionCredentialsEnv)
	data[keeper.MerlionOrgEnv] = os.Getenv(keeper.MerlionOrgEnv)
	data[keeper.MerlionSkladEnv] = os.Getenv(keeper.MerlionSkladEnv)

	data[keeper.NetlabLoginEnv] = os.Getenv(keeper.NetlabLoginEnv)
	data[keeper.NetlabPasswordEnv] = os.Getenv(keeper.NetlabPasswordEnv)
	data[keeper.NetlabOrgEnv] = os.Getenv(keeper.NetlabOrgEnv)
	data[keeper.NetlabSkladEnv] = os.Getenv(keeper.NetlabSkladEnv)

	data[keeper.SofttronikContractorKeyEnv] = os.Getenv(keeper.SofttronikContractorKeyEnv)
	data[keeper.SofttronikContractKeyEnv] = os.Getenv(keeper.SofttronikContractKeyEnv)
	data[keeper.SofttronikOrgEnv] = os.Getenv(keeper.SofttronikOrgEnv)
	data[keeper.SofttronikSkladEnv] = os.Getenv(keeper.SofttronikSkladEnv)

	data[keeper.ElektronmirIDEnv] = os.Getenv(keeper.ElektronmirIDEnv)
	data[keeper.ElektronmirSecretEnv] = os.Getenv(keeper.ElektronmirSecretEnv)
	data[keeper.ElektronmirOrgEnv] = os.Getenv(keeper.ElektronmirOrgEnv)
	data[keeper.ElektronmirSkladOneEnv] = os.Getenv(keeper.ElektronmirSkladOneEnv)

	data[keeper.SkladTokenEnv] = os.Getenv(keeper.SkladTokenEnv)
	data[keeper.CatSkladNameEnv] = os.Getenv(keeper.CatSkladNameEnv)

	for _, v := range data {
		if v == "" {
			log.Fatalf("Данные для входа не обнаружены")
		}
	}

	keeper.K.SetData(data)
}

// func importCSV() {
// 	csv, err := csvInstance.GetCSVInstance()
// 	if err != nil {
// 		log.Printf("Ошибка при открытии codes.csv: %s", err)
// 		return
// 	}
// 	if csv == nil {
// 		fmt.Println("Файл для импорта не обнаружен (codes.csv)")
// 		return
// 	}
// 	err = csv.ImportCodes()
// 	if err != nil {
// 		log.Printf("Ошибка при импорте codes.csv: %s", err)
// 	}
// 	csvInstance.CloseCSV()
// }
