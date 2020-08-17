package main

import (
	"bufio"
	"bytes"
	"encoding/base32"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
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

func SplitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

//https://docs.canarytokens.org/guide/dns-token.html#encoding-additional-information-in-your-token
func assembleDNS(data string, token string) (string, error) {
	b32Data := base32.StdEncoding.EncodeToString([]byte(data))
	b32Data = strings.ReplaceAll(b32Data, "=", "")
	dnsname := strings.Join(SplitSubN(b32Data, 63), ".")
	dnsname += ".G" + fmt.Sprintf("%02d", rand.Intn(99)) + "." + token
	if len(dnsname) > 253 {
		return "", errors.New("DNS name exceeds 253 bytes. Less data needs to be prepended")
	}
	return dnsname, nil
}

func prependData(token string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		results, err := assembleDNS(scanner.Text(), token)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(results)
	}
}

func main() {
	tokenType := flag.String("type", "dns", "Type of token to create (dns or web)")
	email := flag.String("email", "", "Email to configure token on")
	memo := flag.String("memo", "", "Description of token to remember it by")
	previousToken := flag.String("token", "", "DNS token to prepend generic data to. Pipe data from stdin")
	flag.Parse()

	if *previousToken != "" {
		prependData(*previousToken)
		return
	}

	if *email == "" && *previousToken != "" {
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
