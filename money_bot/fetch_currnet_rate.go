package money_bot

//
//import (
//	"encoding/json"
//	"io/ioutil"
//	"log"
//	"net/http"
//)
//
//type rate struct {
//	Buy      string `json:"buy"`
//	Sale     string `json:"sale"`
//	Currency string `json:"ccy"`
//}
//
//const currentyEUR = 0
//const currentyUSD = 2
//
//func fetchRate(currency int) rate {
//	var arr []rate
//
//	r, err := http.Get("https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=3")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer r.Body.Close()
//
//	dataJSON, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = json.Unmarshal([]byte(dataJSON), &arr)
//
//	return arr[currency]
//}
