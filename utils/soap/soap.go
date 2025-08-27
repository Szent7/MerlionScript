package soap

import (
	"MerlionScript/keeper"
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// soapRequest - структура SOAP запроса
type soapRequest struct {
	XMLName   xml.Name `xml:"soap:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soap,attr"`
	XMLNsXSI  string   `xml:"xmlns:xsi,attr"`
	XMLNsXSD  string   `xml:"xmlns:xsd,attr"`
	Body      soapBody
}

// soapBody - структура тела SOAP запроса
type soapBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	Payload any
}

//----------------------------------------------------------------------------------
//! На данный момент SOAP запросы адаптированы под использование интеграцией Мерлион
//----------------------------------------------------------------------------------

// SoapCallHandleResponse выполняет SOAP запрос и обрабатывает ответ
func SoapCallHandleResponse(ws string, action string, payloadInterface any) (*xml.Decoder, error) {
	body, err := soapCall(ws, action, payloadInterface)
	if err != nil {
		return nil, err
	}

	result := xml.NewDecoder(bytes.NewReader(body))
	/*err = xml.Unmarshal(body, &result)
	if err != nil {
		return err
	}*/

	return result, nil
}

// SoapCallHandleResponse выполняет SOAP запрос и возвращает ответ в бинарном виде
func soapCall(ws string, action string, payloadInterface any) ([]byte, error) {
	v := soapRequest{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsXSD:  "http://www.w3.org/2001/XMLSchema",
		XMLNsXSI:  "http://www.w3.org/2001/XMLSchema-instance",
		Body: soapBody{
			Payload: payloadInterface,
		},
	}
	payload, err := xml.MarshalIndent(v, "", "  ") // Сериализует структуру в XML
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(30 * time.Second) // Устанавливает таймаут для HTTP клиента
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", ws, bytes.NewBuffer(payload)) // Создает новый HTTP запрос
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Authorization", "Basic "+keeper.GetMerlionCredentials())

	// дамп запроса для отладки
	/*dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Dump: %q\n", dump)
	*/

	response, err := client.Do(req) // Выполняет HTTP запрос
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body) // Читает тело ответа
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return bodyBytes, nil
}
