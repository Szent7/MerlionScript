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
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	keeper.InitKeeper()

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

	fmt.Println("MScript v2.0")
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
