package moneybot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

const monobankAPIURL = "https://api.monobank.ua"

type webhook struct {
	WebHookURL string `json:"webHookUrl"`
}

type statementItem struct {
	Mcc         int    `json:"mcc"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

func (s statementItem) getCategory() string {
	if s.Mcc >= 3000 && s.Mcc < 4000 || s.Mcc == 4011 || s.Mcc == 4111 || s.Mcc == 4112 || s.Mcc == 4131 || s.Mcc == 4304 || s.Mcc == 4411 || s.Mcc == 4415 || s.Mcc == 4418 || s.Mcc == 4457 || s.Mcc == 4468 || s.Mcc == 4511 || s.Mcc == 4582 || s.Mcc == 4722 || s.Mcc == 4784 || s.Mcc == 4789 || s.Mcc == 5962 || s.Mcc == 6513 || s.Mcc == 7011 || s.Mcc == 7032 || s.Mcc == 7033 || s.Mcc == 7512 || s.Mcc == 7513 || s.Mcc == 7519 {
		return "подорожі"
	}

	if s.Mcc == 4119 || s.Mcc == 5047 || s.Mcc == 5122 || s.Mcc == 5292 || s.Mcc == 5295 || s.Mcc == 5912 || s.Mcc == 5975 || s.Mcc == 5976 || s.Mcc == 5977 || s.Mcc == 7230 || s.Mcc == 7297 || s.Mcc == 7298 || s.Mcc == 8011 || s.Mcc == 8021 || s.Mcc == 8031 || s.Mcc == 8041 || s.Mcc == 8042 || s.Mcc == 8043 || s.Mcc == 8049 || s.Mcc == 8050 || s.Mcc == 8062 || s.Mcc == 8071 || s.Mcc == 8099 {
		return "ліки"
	}

	if s.Mcc >= 7911 && s.Mcc < 7923 || s.Mcc == 5733 || s.Mcc == 5735 || s.Mcc == 5815 || s.Mcc == 5816 || s.Mcc == 5817 || s.Mcc == 5818 || s.Mcc == 5941 || s.Mcc == 5945 || s.Mcc == 5946 || s.Mcc == 5947 || s.Mcc == 5970 || s.Mcc == 5971 || s.Mcc == 5972 || s.Mcc == 5973 || s.Mcc == 7221 || s.Mcc == 7333 || s.Mcc == 7395 || s.Mcc == 7929 || s.Mcc == 7932 || s.Mcc == 7933 || s.Mcc == 7941 || s.Mcc == 7991 || s.Mcc == 7992 || s.Mcc == 7993 || s.Mcc == 7994 || s.Mcc == 7996 || s.Mcc == 7997 || s.Mcc == 7998 || s.Mcc == 7999 || s.Mcc == 8664 {
		return "розваги"
	}

	if s.Mcc >= 5811 && s.Mcc < 5815 {
		return "харчування"
	}

	if s.Mcc == 5297 || s.Mcc == 5298 || s.Mcc == 5300 || s.Mcc == 5311 || s.Mcc == 5331 || s.Mcc == 5399 || s.Mcc == 5411 || s.Mcc == 5412 || s.Mcc == 5422 || s.Mcc == 5441 || s.Mcc == 5451 || s.Mcc == 5462 || s.Mcc == 5499 || s.Mcc == 5715 || s.Mcc == 5921 {
		return "продукти"
	}
	if s.Mcc == 7829 || s.Mcc == 7832 || s.Mcc == 7841 {
		// actually this is кіно
		return "розваги"
	}

	if s.Mcc == 5172 || s.Mcc == 5511 || s.Mcc == 5532 || s.Mcc == 5533 || s.Mcc == 5541 || s.Mcc == 5542 || s.Mcc == 5983 || s.Mcc == 7511 || s.Mcc == 7523 || s.Mcc == 7531 || s.Mcc == 7534 || s.Mcc == 7535 || s.Mcc == 7538 || s.Mcc == 7542 || s.Mcc == 7549 {
		return "авто"
	}

	if s.Mcc == 5131 || s.Mcc == 5137 || s.Mcc == 5139 || s.Mcc == 5611 || s.Mcc == 5621 || s.Mcc == 5631 || s.Mcc == 5641 || s.Mcc == 5651 || s.Mcc == 5655 || s.Mcc == 5661 || s.Mcc == 5681 || s.Mcc == 5691 || s.Mcc == 5697 || s.Mcc == 5698 || s.Mcc == 5699 || s.Mcc == 5931 || s.Mcc == 5948 || s.Mcc == 5949 || s.Mcc == 7251 || s.Mcc == 7296 {
		return "одяг"
	}
	if s.Mcc == 4121 {
		return "транспорт"
	}
	if s.Mcc == 0742 || s.Mcc == 5995 {
		// actually this is тварини
		return ""
	}
	if s.Mcc == 2741 || s.Mcc == 5111 || s.Mcc == 5192 || s.Mcc == 5942 || s.Mcc == 5994 {
		// actually this is книги
		return ""
	}
	if s.Mcc == 5992 || s.Mcc == 5193 {
		// actually this is квіти
		return ""
	}

	return ""
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

// SetWebhook sets webhook for monobank
func SetWebhook(token string, url string, port int) error {
	data, _ := json.Marshal(webhook{WebHookURL: url})
	r, _ := http.NewRequest("POST", fmt.Sprintf("%s/personal/webhook", monobankAPIURL), bytes.NewReader(data))
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
		return fmt.Errorf("Status code: %d, Resp: %s", resp.StatusCode, string(body))
	}
	return nil
}

// ListenWebhook runs listener for webhook
func ListenWebhook(port int, monobankEvents chan Item) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)

		if req.Method == http.MethodGet {
			return
		}

		body, err := ioutil.ReadAll(req.Body)

		logrus.Infof("webhook data: %s", string(body))

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
			Name:     event.Data.StatementItem.Description,
			Amount:   event.Data.StatementItem.getNormalizedAmount(),
			Category: event.Data.StatementItem.getCategory(),
		}

		if item.Name == "" {
			if item.Category != "" {
				item.Name = item.Category
			} else {
				item.Name = "transaction"
			}
		}

		monobankEvents <- item
	})

	logrus.Info(fmt.Sprintf("listen webhook on: %d", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logrus.Error(err)
	}
}
