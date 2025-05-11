package requests

import (
	skladTypes "MerlionScript/types/restTypes/sklad"
	"MerlionScript/utils/rest"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

func UploadImage(id string, image skladTypes.UploadImage) (rest.Response, error) {
	url := fmt.Sprintf(skladTypes.ImageUrl, id)

	jsonBody, err := json.MarshalIndent(image, "", "  ")
	if err != nil {
		log.Println("ошибка при преобразовании структуры в JSON:", err)
		return rest.Response{}, err
	}

	body, err := rest.CreateRequestMS("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return rest.Response{}, err
	}

	return *body, nil
}

func GetItemsImagesData(id string) (rest.Response, error) {
	url := fmt.Sprintf(skladTypes.ImageUrl, id)
	body, err := rest.CreateRequestMS("GET", url, nil)
	if err != nil {
		return rest.Response{}, err
	}
	return *body, nil
}
