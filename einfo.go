package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
	"strconv"
)

type commOutput struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type Tx struct {
	TxHash     string `json:"tx_hash"`
	TxContent  string `json:"tx_content"`
	Pending    bool   `json:"pending"`
	CreateTime int64  `json:"create_time"`
}

func formatOutput(ret int, msg string) string {
	output := &commOutput{ret, msg}
	outputJson, _ := json.Marshal(output)
	return string(outputJson)
}

func saveTx(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var (
		txHash    string
		txContent string
	)
	if txHashSlice, ok := r.Form["txHash"]; ok && len(txHashSlice) > 0 {
		txHash = txHashSlice[0]
	}
	if txContentSlice, ok := r.Form["txContent"]; ok && len(txContentSlice) > 0 {
		txContent = txContentSlice[0]
	}

	if txHash == "" || txContent == "" {
		fmt.Fprint(w, formatOutput(1, "empty txHash or empty txContent"))
		return
	}

	nowTime := time.Now().Unix()
	nowTimeStr := strconv.FormatInt(nowTime, 10)
	set(txHash, "0" + nowTimeStr + txContent)

	fmt.Fprint(w, formatOutput(0, ""))
}

func confirmTx(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var txHash string
	if txHashSlice, ok := r.Form["txHash"]; ok && len(txHashSlice) > 0 {
		txHash = txHashSlice[0]
	}

	if txHash == "" {
		fmt.Fprint(w, formatOutput(1, "empty txHash"))
		return
	}

	txContent, err := get(txHash)
	if err != nil {
		fmt.Fprint(w, formatOutput(1, err.Error()))
		return
	}

	txContent[0] = '1'
	set(txHash, string(txContent))

	fmt.Fprint(w, formatOutput(0, ""))
}

func txList(w http.ResponseWriter, r *http.Request) {
	keyList := keys()

	var txList []*Tx
	for _, key := range keyList {
		keyConent, err := get(key)
		if err != nil {
			continue
		}

		tx := new(Tx)
		tx.TxHash = key
		if keyConent[0] == '0' {
			tx.Pending = true
		}

		txTimeStr := keyConent[1:11]
		txTime, _ := strconv.ParseInt(string(txTimeStr), 10, 64)
		tx.CreateTime = txTime

		tx.TxContent = string(keyConent[11:])

		txList = append(txList, tx)
	}

	txListJson, _ := json.Marshal(txList)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(txListJson))
}

func main() {
	f := http.FileServer(http.Dir("./static"))
	http.Handle("/static/css/", http.StripPrefix("/static/", f))
	http.Handle("/static/js/", http.StripPrefix("/static/", f))
	http.Handle("/", f)

	http.HandleFunc("/saveTx", saveTx)
	http.HandleFunc("/confirmTx", confirmTx)
	http.HandleFunc("/txList", txList)

	log.Fatal(http.ListenAndServe(":8089", nil))
}
