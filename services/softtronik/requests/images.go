package requests

import (
	"MerlionScript/keeper"
	softtronikTypesRest "MerlionScript/types/restTypes/softtronik"
	"MerlionScript/utils/rest"
	"encoding/json"
	"fmt"
)

func getItemsImages(itemCode string) (softtronikTypesRest.ImageItem, error) {
	url := fmt.Sprintf(softtronikTypesRest.ImageUrl, keeper.K.GetSofttronikContractor(), itemCode)

	response, err := rest.CreateRequest("GET", url, nil, "")
	if err != nil {
		return softtronikTypesRest.ImageItem{}, err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return softtronikTypesRest.ImageItem{}, err
	}

	var fileList softtronikTypesRest.ImageItem
	if err := json.Unmarshal(response.Body, &fileList); err != nil {
		return softtronikTypesRest.ImageItem{}, fmt.Errorf("ошибка при декодировании item (getitemsimages): %s", err.Error())
	}

	return fileList, nil
}

func GetImagesByItemIdFormatted(id string) (map[string]string, error) {
	rawImages, err := getItemsImages(id)
	if err != nil {
		return nil, err
	}

	var images map[string]string = make(map[string]string, 5)
	for i := range rawImages.DataFiles {
		if rawImages.DataFiles[i].ExtensionFile == "png" || rawImages.DataFiles[i].ExtensionFile == "jpg" || rawImages.DataFiles[i].ExtensionFile == "jpeg" {
			images[rawImages.DataFiles[i].NameFile] = rawImages.DataFiles[i].LinkToDataFile
		}
	}

	return images, nil
}
