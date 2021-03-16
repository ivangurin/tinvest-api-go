package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/ivangurin/tinvest-client-go"
	"github.com/ivangurin/tinvest-analyser-go"
	"time"
)

func returnRoot(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	ioResponse.WriteHeader(http.StatusBadRequest)

}

func returnPositions(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	lvBearerToken := ioRequest.Header.Get("Authorization")

	if lvBearerToken == "" {
		ioResponse.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(ioResponse, http.StatusText(http.StatusUnauthorized))
		return
	}

	loClient := tinvestclient.Client{}

	loClient.Init(lvBearerToken[7:])

	ltPositions, loError := loClient.GetPositions()

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=UTF-8")

	var lvJsonBytes []byte

	if len(ltPositions) > 0{
		lvJsonBytes, _ = json.Marshal(ltPositions)
	}else {
		lvJsonBytes = []byte("[]")
	}

	ioResponse.Write(lvJsonBytes)

}

func returnProfit(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	lvBearerToken := ioRequest.Header.Get("Authorization")

	if lvBearerToken == "" {
		ioResponse.WriteHeader(http.StatusUnauthorized)
		return
	}

	lvTicker := ioRequest.URL.Query().Get("ticker")

	fmt.Println(lvTicker)

	loAnalyzer := tinvestanalyser.Analyser{}

	loAnalyzer.Init(lvBearerToken[7:])

	ltProfit, loError := loAnalyzer.GetProfit(lvTicker, time.Now().AddDate(-10, 0, 0), time.Now())

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=UTF-8")

	var lvJsonBytes []byte

	if len(ltProfit) > 0{
		lvJsonBytes, _ = json.Marshal(ltProfit)
	}else {
		lvJsonBytes = []byte("[]")
	}

	ioResponse.Write(lvJsonBytes)

}

func main() {

	loRouter := mux.NewRouter()

	loRouter.HandleFunc("/", returnRoot).Methods("GET")
	loRouter.HandleFunc("/positions", returnPositions).Methods("GET")
	loRouter.HandleFunc("/profit", returnProfit).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", loRouter))

}
