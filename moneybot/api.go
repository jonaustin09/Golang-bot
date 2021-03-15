package moneybot

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a Application) startApiServer(){
	http.HandleFunc("/api/export", func(w http.ResponseWriter, r *http.Request) {

		items, err := a.LogItemRepository.GetRecords()
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}

		logrus.Infof("Fetch items count %v", len(items))

		serializedItems, err := json.Marshal(items)
		w.Header().Set("Content-Type", "application/json")
		w.Write(serializedItems)
	})

	logrus.Info(fmt.Sprintf("listen api on: %d", a.Config.APIServer))
	err := http.ListenAndServe(fmt.Sprintf(":%d", a.Config.APIServer), nil)
	if err != nil {
		logrus.Error(err)
	}
}