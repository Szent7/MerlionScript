package soap

import (
	"MerlionScript/keeper"
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type soapRQ struct {
	XMLName   xml.Name `xml:"soap:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soap,attr"`
	XMLNsXSI  string   `xml:"xmlns:xsi,attr"`
	XMLNsXSD  string   `xml:"xmlns:xsd,attr"`
	Body      soapBody
}

type soapBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	Payload interface{}
}

// func SoapCallHandleResponse(ws string, action string, payloadInterface interface{}, result interface{}) error {
func SoapCallHandleResponse(ws string, action string, payloadInterface interface{}) (*xml.Decoder, error) {
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

func soapCall(ws string, action string, payloadInterface interface{}) ([]byte, error) {
	v := soapRQ{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsXSD:  "http://www.w3.org/2001/XMLSchema",
		XMLNsXSI:  "http://www.w3.org/2001/XMLSchema-instance",
		Body: soapBody{
			Payload: payloadInterface,
		},
	}
	payload, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", ws, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Authorization", "Basic "+keeper.GetMerlionCredentials())

	/*dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Dump: %q\n", dump)
	*/
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("Body Bytes:%s\n", string(bodyBytes))
	defer response.Body.Close()
	return bodyBytes, nil
}
