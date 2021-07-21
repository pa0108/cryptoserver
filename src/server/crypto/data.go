package crypto

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

//Result ... outdata struct
type Result struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
}

const (
	API_BASE  = "https://api.hitbtc.com/api/2"
	WAIT_TIME = 10
)

var symbols = []string{"BTCUSD", "ETHBTC"}
var outData = make(map[string]interface{})

// GetData ... call to crypto server for latest values on regular intervals
func (c *Client) GetData() {
	for {
		for _, symbol := range symbols {
			var result Result
			c.callAPI("public/symbol/"+symbol, &result)
			c.callAPI("public/ticker/"+symbol, &result)
			//c.callAPI("public/currency/"+result.ID, &result)
			outData[symbol] = result
		}
		time.Sleep(WAIT_TIME * time.Second)
	}
}

func (c *Client) callAPI(api string, res *Result) {
	//Get data from ticker API
	r, err := c.do("GET", api, nil, false)
	if err != nil {
		log.Printf("Http request failed. API - %s, Error - %s", api, err)
	} else {
		err = json.Unmarshal(r, &res)
		if err != nil {
			log.Printf("Unable to unmarshall the response - %s", err)
		}
	}
}

//NewCurrencyHandler ... outdata handler
func NewCurrencyHandler() *Result {
	return &Result{}
}

//ServeHTTP ... func to servce http request
func (c *Result) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	symbol := r.RequestURI[strings.LastIndex(r.RequestURI, "/")+1:]

	var resp interface{}
	if strings.ToLower(symbol) == "all" {
		var allVal []interface{}
		for _, v := range outData {
			allVal = append(allVal, v)
		}
		resp = map[string][]interface{}{"currencies": allVal}

	} else {
		resp = outData[strings.ToUpper(symbol)]
	}

	if resp == nil {
		http.Error(rw, fmt.Sprintf("No currency with symbol - %s", symbol), http.StatusNotFound)
		return
	}

	d, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	} else {
		rw.Write(d)
		return
	}
}
