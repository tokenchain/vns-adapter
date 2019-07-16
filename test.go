package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func sendRequest() {
	// Request (POST http://192.168.123.162:8545/?Content-Type=application%2Fjso)

	json := []byte(`{"jsonrpc": "2.0","method": "vns_getBalance","id": 1,"params": ["0x142b097b802b5224f9d4bfc93db189b4f4621df2","latest"]}`)
	body := bytes.NewBuffer(json)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://192.168.31.142:2000/?Content-Type=application%2Fjso", body)

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}

func main() {
	sendRequest()
}
