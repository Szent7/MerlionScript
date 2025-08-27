package rest

import (
	"MerlionScript/keeper"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// Response - структура HTTP ответа
type Response struct {
	Body       []byte
	StatusCode int
}

// Создает HTTP-запрос на основе заданных параметров и выполняет его
func CreateRequest(reqType string, url string, bodyRequest io.Reader, bearerToken string) (*Response, error) {
	// Создание HTTP-запроса
	req, err := http.NewRequest(reqType, url, bodyRequest)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil, err
	}

	// Устанавливает заголовки запроса
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	if bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bearerToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req) // Выполняет HTTP запрос
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()

	/*if resp.StatusCode != 200 {
		fmt.Println("Неверный статус код:", resp.Status)
		//return nil, err
	}*/
	var body []byte
	// Обработка ответа с сжатием и без
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("ошибка при создании gzip.reader:", err)
			return nil, err
		}
		defer gz.Close()
		body, err = io.ReadAll(gz)
		if err != nil {
			fmt.Println("ошибка при чтении декомпрессированного тела ответа:", err)
			return nil, err
		}
	} else {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ошибка при чтении тела ответа:", err)
			return nil, err
		}
	}

	return &Response{Body: body, StatusCode: resp.StatusCode}, nil
}

// Создает HTTP-запрос на основе заданных параметров и выполняет его. Только для интеграции ЭМ
func CreateRequestElektronmir(reqType string, url string, bodyRequest io.Reader, Token string) (*Response, error) {
	// Создание HTTP-запроса
	req, err := http.NewRequest(reqType, url, bodyRequest)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil, err
	}

	// Устанавливает заголовки запроса
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	if Token != "" {
		req.Header.Add("Authorization", "token "+Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req) // Выполняет HTTP запрос
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()

	/*if resp.StatusCode != 200 {
		fmt.Println("Неверный статус код:", resp.Status)
		//return nil, err
	}*/
	var body []byte
	// Обработка ответа с сжатием и без
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("ошибка при создании gzip.reader:", err)
			return nil, err
		}
		defer gz.Close()
		body, err = io.ReadAll(gz)
		if err != nil {
			fmt.Println("ошибка при чтении декомпрессированного тела ответа:", err)
			return nil, err
		}
	} else {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ошибка при чтении тела ответа:", err)
			return nil, err
		}
	}

	return &Response{Body: body, StatusCode: resp.StatusCode}, nil
}

// Создает HTTP-запрос на основе заданных параметров и выполняет его. Только для запросов на получение изображений. Возвращает Content-Type
func CreateRequestImageHeader(reqType string, url string, bodyRequest io.Reader, bearerToken string) (*Response, string, error) {
	// Создание HTTP-запроса
	req, err := http.NewRequest(reqType, url, bodyRequest)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil, "", err
	}

	// Устанавливает заголовки запроса
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	if bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bearerToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req) // Выполняет HTTP запрос
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, "", err
	}
	defer resp.Body.Close()

	/*if resp.StatusCode != 200 {
		fmt.Println("Неверный статус код:", resp.Status)
		//return nil, err
	}*/
	var body []byte
	// Обработка ответа с сжатием и без
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("ошибка при создании gzip.reader:", err)
			return nil, "", err
		}
		defer gz.Close()
		body, err = io.ReadAll(gz)
		if err != nil {
			fmt.Println("ошибка при чтении декомпрессированного тела ответа:", err)
			return nil, "", err
		}
	} else {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ошибка при чтении тела ответа:", err)
			return nil, "", err
		}
	}

	return &Response{Body: body, StatusCode: resp.StatusCode}, resp.Header.Get("Content-Type"), nil
}

// Обертка над HTTP-запросом. Только для ERP системы. Использует данные для входа ERP системы
func CreateRequestMS(reqType string, url string, bodyRequest io.Reader) (*Response, error) {
	return CreateRequest(reqType, url, bodyRequest, keeper.GetMSCredentials())
}

// Обертка над HTTP-запросом. Только если ответ приходит в виде XML.
func CreateRequestXML(reqType string, url string, bodyRequest io.Reader) (*xml.Decoder, error) {
	resp, err := CreateRequest(reqType, url, bodyRequest, "")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ошибка при выполнении запроса (createrequestxml) statuscode: %d", resp.StatusCode)
	}

	return xml.NewDecoder(bytes.NewReader(resp.Body)), nil
}
