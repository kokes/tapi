package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type IData struct {
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	FormData map[string]string `json:"form_data"`
}

// response data
type RData struct {
	StatusCode int               `json:"status_code"`
	Protocol   string            `json:"protocol"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Timing     int64             `json:"timing"`
}

func main() {
	fmt.Println("Serving at http://localhost:7678")
	http.HandleFunc("/", route)
	log.Fatal(http.ListenAndServe(":7678", nil))
}

func route(w http.ResponseWriter, r *http.Request) {
	pth := strings.TrimRight(r.URL.Path, "/")

	switch pth {
	case "":
		// fmt.Fprintf(w, "index this, biatch")
		w.Header().Set("Content-Type", "text/html")

		dt, _ := ioutil.ReadFile("index.html")
		w.Write(dt)

	case "/ping":
		fmt.Fprintf(w, "pong")
	case "/query":
		dt, err := handleRequest(w, r)
		if err != nil {
			panic(err) // TODO: return JSON with the error instead
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(dt)

	default:
		fmt.Fprintf(w, "ne hablo")

	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var idta IData
	bd, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(bd, &idta)
	if err != nil {
		panic(err)
	}

	method := strings.ToUpper(idta.Method)

	rd, err := makeRequest(method, idta)
	if err != nil {
		return []byte{}, err
	}

	b, err := json.Marshal(rd)
	if err != nil {
		return []byte{}, err
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	// ioutil.WriteFile("res-b64.txt", []byte(rd.Body), 0666)

	return b, nil

}

func makeRequest(method string, bd IData) (RData, error) {
	client := &http.Client{} // TODO: handle redirects?

	req, err := http.NewRequest(method, bd.URL, nil)
	if err != nil {
		return RData{}, err
	}

	t := time.Now()
	resp, err := client.Do(req)
	tm := time.Since(t)
	if err != nil {
		return RData{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// ioutil.WriteFile("res.txt", body, 0666)
	if err != nil {
		return RData{}, err
	}
	b64body := base64.StdEncoding.EncodeToString(body)

	// the header is map[string][]string, we need to flatten it
	hd := make(map[string]string)
	for k, v := range resp.Header {
		hd[k] = strings.Join(v, "; ")
	}

	rd := RData{
		StatusCode: resp.StatusCode,
		Protocol:   resp.Proto,
		Headers:    hd,
		Body:       b64body,
		Timing:     int64(tm) / (1000 * 1000),
	}

	return rd, nil
}
