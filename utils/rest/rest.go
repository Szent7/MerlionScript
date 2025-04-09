package rest

import (
	"MerlionScript/keeper"
	"MerlionScript/types/restTypes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

func CreateRequest(reqType string, url string, bodyRequest io.Reader) (*restTypes.Response, error) {
	req, err := http.NewRequest(reqType, url, bodyRequest)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	//if authHeader != "" {
	req.Header.Add("Authorization", "Bearer "+keeper.K.GetMSCredentials())
	//}

	client := &http.Client{}
	resp, err := client.Do(req)
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
	return &restTypes.Response{Body: body, StatusCode: resp.StatusCode}, nil
}
