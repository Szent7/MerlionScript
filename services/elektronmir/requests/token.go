package requests

import (
	"MerlionScript/keeper"
	elektronmirTypes "MerlionScript/services/elektronmir/types"
	"MerlionScript/utils/rest"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

func GetNewToken() (string, error) {
	request := elektronmirTypes.TokenRequest{
		ClientID:     keeper.GetElektronmirID(),
		ClientSecret: keeper.GetElektronmirSecret(),
	}
	jsonBody, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		log.Println("ошибка при преобразовании структуры в JSON:", err)
		return "", err
	}

	response, err := rest.CreateRequestElektronmir("POST", elektronmirTypes.TokenUrl, bytes.NewBuffer(jsonBody), "")
	if err != nil || response.StatusCode != 200 {
		return "", err
	}

	tokenStruct := elektronmirTypes.TokenResponse{}
	if err := json.Unmarshal(response.Body, &tokenStruct); err != nil {
		return "", fmt.Errorf("ошибка при декодировании tokenstruct (gettemptoken): %s", err.Error())
	}
	if tokenStruct.Status != "ok" {
		return "", fmt.Errorf("ошибка при получении токена (gettemptoken): %s", tokenStruct.Status)
	}

	return tokenStruct.Token, nil
}
