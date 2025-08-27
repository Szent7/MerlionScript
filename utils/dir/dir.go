package dir

import "os"

// Проверка существования директории
func CheckDirectoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// Проверка существования директории
func CheckFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Создание директории
func CreateDirectoryDefault(path string) error {
	return os.MkdirAll(path, 0660)
}

// Удаляет директорию/файл
func RemovePath(path string) error {
	return os.Remove(path)
}
