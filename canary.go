package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

type Token struct {
	Token    string `json:"Token"`
	Hostname string `json:"Hostname"`
	Url      string `json:"Url"`
}

func generateToken(uri string, params map[string]string) (*http.Request, error) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}

func main() {
	tokenType := flag.String("type", "dns", "Type of token to create (dns or web)")
	email := flag.String("email", "", "Email to configure token on")
	memo := flag.String("memo", "", "Description of token to remember it by")
	flag.Parse()

	if *email == "" {
		flag.PrintDefaults()
		return
	}

	params := map[string]string{
		"type":                     *tokenType,
		"email":                    *email,
		"webhook":                  "",
		"fmt":                      "",
		"memo":                     *memo,
		"clonedsite":               "",
		"sql_server_table_name":    "TABLE1",
		"sql_server_view_name":     "VIEW1",
		"sql_server_function_name": "FUNCTION1",
		"sql_server_trigger_name":  "TRIGGER1",
		"redirect_url":             "",
	}
	request, err := generateToken("https://canarytokens.org/generate", params)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Origin", "https://canarytokens.org")
	request.Header.Set("Referer", "https://canarytokens.org/generate")
	request.Header.Set("Connection", "close")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	var token Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		log.Fatal(err)
	}
	if *tokenType == "dns" {
		fmt.Printf("%s\n", token.Hostname)
	} else if *tokenType == "web" {
		fmt.Printf("%s\n", token.Url)
	} else {
		data, err := json.Marshal(token)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", data)
	}
}
