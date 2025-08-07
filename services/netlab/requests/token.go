package requests

import (
	"MerlionScript/keeper"
	netlabTypes "MerlionScript/services/netlab/types"
	"MerlionScript/utils/rest"
	"encoding/json"
	"fmt"
)

func GetNewToken() (string, error) {
	netlabLogin, netlabPassword := keeper.K.GetCredentialsNetlab()
	url := fmt.Sprintf(netlabTypes.TokenUrl, netlabLogin, netlabPassword)

	response, err := rest.CreateRequest("GET", url, nil, "")
	if err != nil || response.StatusCode != 200 {
		return "", err
	}

	jsonData := string(response.Body)
	cleanedData := jsonData[6:]

	tokenStruct := netlabTypes.TokenMessage{}
	if err := json.Unmarshal([]byte(cleanedData), &tokenStruct); err != nil {
		return "", fmt.Errorf("ошибка при декодировании tokenstruct (gettemptoken): %s", err.Error())
	}
	if tokenStruct.TokenResponse.Status.Code != "200" {
		return "", fmt.Errorf("ошибка при получении токена (gettemptoken): %s", tokenStruct.TokenResponse.Status.Message)
	}

	return tokenStruct.TokenResponse.Data.Token, nil
}
