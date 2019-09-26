package moneybot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const monobankApiUrl = "https://api.monobank.ua"

type webhook struct {
	WebHookUrl string `json:"webHookUrl"`
}

type statementItem struct {
	Mcc    int `json:"mcc"`
	Amount int `json:"amount"`
}

func (s statementItem) getCategory() string {
	// TODO: get right category based on mcc code
	return "продукти"
}

func (s statementItem) getNormalizedAmount() float64 {
	return float64(s.Amount) / 100 * -1
}

type webhookEvent struct {
	Type string `json:"type"`
	Data struct {
		StatementItem statementItem `json:"statementItem"`
	} `json:"data"`
}

func SetWebhook(token string, url string) error {
	data, _ := json.Marshal(webhook{WebHookUrl: url})
	r, _ := http.NewRequest("POST", fmt.Sprintf("%s/personal/webhook", monobankApiUrl), bytes.NewReader(data))
	r.Header.Set("X-Token", token)
	r.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("Status code: %s, Resp: %s", resp.StatusCode, string(body)))
	}
	return nil
}

func ListenWebhook(port int, monobankEvents chan Item) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		fmt.Println(string(body))
		if err != nil {
			logrus.Error(err)
			return
		}

		event := webhookEvent{}

		err = json.Unmarshal(body, &event)
		if err != nil {
			logrus.Error(err)
			return
		}
		item := Item{
			Name:     event.Data.StatementItem.getCategory(),
			Amount:   event.Data.StatementItem.getNormalizedAmount(),
			Category: event.Data.StatementItem.getCategory(),
		}

		monobankEvents <- item

		w.WriteHeader(200)
	})
	logrus.Info(fmt.Sprintf("listen webhook on: %d", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logrus.Error(err)
	}
}
