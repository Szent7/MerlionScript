package main

import (
	"MerlionScript/controller"
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	"MerlionScript/services/common/initializer"
	"MerlionScript/utils/backup"
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

	// Инициализация хранилища
	keeper.InitKeeper()

	// Инициализация(регистрация) сервисов
	initializer.InitServices()

	// Инициализация системы бэкапов, если включено в настройках
	if keeper.GetBackupToggle() {
		backup.InitBackup(backup.BackupObj{
			SrcPath:      "./data",
			BackupDir:    "./data/backup",
			BackupNumber: keeper.GetBackupNumber(),
		})
	}

	// Получаем инициализированный экземпляр базы данных
	dbInstance, err := db.GetDB(typesDB.PathDB)
	defer db.CloseDB()
	if err != nil {
		log.Fatalf("Ошибка при создании экземпляра базы данных: %s", err)
	}
	// Инициализация таблицы в базе данных
	err = dbInstance.Init(common.GetTableNames())
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %s", err)
	}

	// Инициализация кэша
	err = cache.InitCache(ctx)
	if err != nil {
		log.Fatalf("Ошибка при инициализации кэша: %s", err)
	}

	// Создание канала для сигналов завершения (SIGINT, SIGTERM)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	fmt.Println("MScript v2.7.4")
	fmt.Println("Боже, Царя храни!")
	// Запуск контроллера в отдельной горутине
	wg.Add(1)
	go controller.StartController(ctx, &wg, dbInstance)

	// Ожидание сигнала завершения или отмены контекста
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
