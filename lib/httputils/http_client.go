package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/netsec-ethz/conn-tester/lib/common"
)

type ClientWrapper struct {
	Client *http.Client
}

func (c *ClientWrapper) SendCommand(host string, method string, payload interface{}) (bool, error) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)

	req, err := http.NewRequest(method, host, buf)
	if err != nil {
		return false, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var commandResponse ServerResponseMessage
	jsonParser := json.NewDecoder(resp.Body)
	jsonParser.Decode(&commandResponse)
	if commandResponse.Success {
		return commandResponse.Success, nil
	} else {
		return commandResponse.Success, fmt.Errorf(commandResponse.Message)
	}

}

func (c *ClientWrapper) SendRequestCompareResponse(host string, method string,
	request []byte, expectedResponse []byte) (bool, error) {

	req, err := http.NewRequest(method, host, bytes.NewReader(request))
	if err != nil {
		return false, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if equal, err := common.CompareBytesWithStream(expectedResponse, resp.Body); err != nil {
		return false, err
	} else {
		return equal, nil
	}
}

func CreateNewClientWrapper(timeout time.Duration) (*ClientWrapper, error) {
	client := &http.Client{Timeout: timeout}
	wrapper := &ClientWrapper{Client: client}

	return wrapper, nil
}
