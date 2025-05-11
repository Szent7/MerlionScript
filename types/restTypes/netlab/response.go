package netlab

type TokenMessage struct {
	TokenResponse TokenResponse `json:"tokenResponse"`
}

type TokenResponse struct {
	Status status    `json:"status"`
	Data   tokenData `json:"data"`
}

type status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type tokenData struct {
	Token     string `json:"token"`
	ExpiredIn string `json:"expiredIn"`
}
