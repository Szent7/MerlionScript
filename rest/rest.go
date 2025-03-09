package rest

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

func CreateRequest(reqType string, url string, authHeader string, bodyRequest io.Reader) ([]byte, error) {
	req, err := http.NewRequest(reqType, url, bodyRequest)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Неверный статус код:", resp.Status)
		//return nil, err
	}
	var body []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("Ошибка при создании gzip.reader:", err)
			return nil, err
		}
		defer gz.Close()
		body, err = io.ReadAll(gz)
		if err != nil {
			fmt.Println("Ошибка при чтении декомпрессированного тела ответа:", err)
			return nil, err
		}
	} else {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Ошибка при чтении тела ответа:", err)
			return nil, err
		}
	}
	return body, nil
}
