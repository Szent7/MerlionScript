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

var bck *BackupObj

type BackupObj struct {
	SrcPath      string
	BackupDir    string
	BackupNumber int
}

func init() {
	bck = new(BackupObj)
}

func InitBackup(newBackup BackupObj) {
	bck = &newBackup
}

func CreateDefaultBackup() error { return bck.createDefaultBackup() }
func (bk *BackupObj) createDefaultBackup() error {
	exists, err := dir.CheckDirectoryExists(bck.BackupDir)
	if err != nil {
		return err
	}
	if !exists {
		if err = dir.CreateDirectoryDefault(bck.BackupDir); err != nil {
			return err
		}
	}

	if !isBackupAlreadyDoneToday(bk.BackupDir) {
		if err := createBackup(bk.SrcPath, bk.BackupDir); err != nil {
			return err
		}
	}

	return removeOldBackups(bk.BackupDir)
}

func createBackup(srcPath, backupDir string) error {
	now := time.Now()
	timestamp := now.Format("2006-01-02")
	backupName := fmt.Sprintf("backup_%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	srcFile, err := os.Open(filepath.Join(srcPath, "database.db"))
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

func removeOldBackups(backupDir string) error {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	var backups []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "backup_") && strings.HasSuffix(file.Name(), ".db") {
			backups = append(backups, file)
		}
	}

	if len(backups) > 7 {
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
		for i := 0; i < len(backups)-7; i++ {
			file := backups[i]
			filePath := filepath.Join(backupDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func isBackupAlreadyDoneToday(backupDir string) bool {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		fmt.Println("Ошибка при сканировании директории:", err)
		return false
	}

	today := time.Now().Format("2006-01-02")

	for _, file := range files {
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
