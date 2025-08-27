package backup

import (
	"MerlionScript/utils/dir"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Глобальная переменная для хранения текущего экземпляра резервного копирования
var bck *BackupObj

// BackupObj - структура объекта резервного копирования
type BackupObj struct {
	SrcPath      string
	BackupDir    string
	BackupNumber int
}

// InitBackup инициализирует объект резервного копирования
func InitBackup(newBackup BackupObj) {
	bck = &newBackup
}

// Публичная обертка над createDefaultBackup()
func CreateDefaultBackup() error { return bck.createDefaultBackup() }

// createDefaultBackup проверяет, существует ли уже резервная копия для сегодняшнего дня.
// Если нет, создает новую резервную копию и удаляет старые, если их больше чем BackupNumber.
func (bk *BackupObj) createDefaultBackup() error {
	// Проверка существования директории бэкапов
	exists, err := dir.CheckDirectoryExists(bck.BackupDir)
	if err != nil {
		return err
	}
	if !exists {
		// Создание директории бэкапов
		if err = dir.CreateDirectoryDefault(bck.BackupDir); err != nil {
			return err
		}
	}

	// Проверка существования бэкапа сегодняшнего дня
	if !isBackupAlreadyDoneToday(bk.BackupDir) {
		// Создание бэкапа
		if err := createBackup(bk.SrcPath, bk.BackupDir); err != nil {
			return err
		}
	}

	// Удаление старых бэкапов
	return removeOldBackups(bk.BackupDir, bk.BackupNumber)
}

// createBackup создает резервную копию файла
func createBackup(srcPath, backupDir string) error {
	// Формирование названия бэкапа
	now := time.Now()
	timestamp := now.Format("2006-01-02")
	backupName := fmt.Sprintf("backup_%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupName) // Путь до нового бэкапа

	// Открытие файла БД
	srcFile, err := os.Open(filepath.Join(srcPath, "database.db"))
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Создание копии файла БД
	dstFile, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Копирование файла в копию
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// removeOldBackups удаляет старые резервные копии, если их больше чем backupNumber
func removeOldBackups(backupDir string, backupNumber int) error {
	// Получение списка файлов и директорий
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	// Формирование списка существующих бэкапов
	var backups []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "backup_") && strings.HasSuffix(file.Name(), ".db") {
			backups = append(backups, file)
		}
	}

	// Проверка кол-ва бэкапов
	if len(backups) > backupNumber {
		// Сортировка по времени создания
		sort.Slice(backups, func(i, j int) bool {
			infoI, errI := backups[i].Info()
			infoJ, errJ := backups[j].Info()
			if errI != nil || errJ != nil {
				return false
			}
			return infoI.ModTime().Before(infoJ.ModTime())
		})

		// Удаление старых бэкапов
		for i := 0; i < len(backups)-backupNumber; i++ {
			file := backups[i]
			filePath := filepath.Join(backupDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

// isBackupAlreadyDoneToday проверяет, была ли уже выполнена резервная копия для сегодняшнего дня.

func isBackupAlreadyDoneToday(backupDir string) bool {
	// Получение списка файлов и директорий
	files, err := os.ReadDir(backupDir)
	if err != nil {
		fmt.Println("Ошибка при сканировании директории:", err)
		return false
	}

	// Получение сегодняшней даты
	today := time.Now().Format("2006-01-02")

	// Цикл по списку файлов и директорий
	for _, file := range files {
		// Пропускаем директории
		if file.IsDir() {
			continue
		}

		// Пропускаем файлы, которые не являются бэкапами
		if !strings.HasPrefix(file.Name(), "backup_") {
			continue
		}

		// Получаем информацию о файле
		fileInfo, err := file.Info()
		if err != nil {
			continue // или обрабатывайте ошибку
		}

		// Получаем время модификации
		modTime := fileInfo.ModTime()
		modTimeDay := modTime.Format("2006-01-02")
		if modTimeDay == today {
			return true
		}
	}

	return false
}
