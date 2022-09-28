package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"pvg/entity"
)

func CallCheckUser(input entity.CheckUserExist) (result []byte, err error) {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/checkuserexist", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	respServer, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer respServer.Body.Close()
	result, err = io.ReadAll(respServer.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
