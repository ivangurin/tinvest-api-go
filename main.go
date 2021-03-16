package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/ivangurin/tinvest-client-go"
	"github.com/ivangurin/tinvest-analyser-go"
	"strings"
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

	if len(ltPositions) == 0 {
		ltPositions = make([]tinvestclient.Position, 0)
	}

	lvBody, loError := json.Marshal(ltPositions)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=UTF-8")

	ioResponse.Write(lvBody)

}

func returnProfit(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	lvBearerToken := ioRequest.Header.Get("Authorization")

	if lvBearerToken == "" {
		ioResponse.WriteHeader(http.StatusUnauthorized)
		return
	}

	loParameters := mux.Vars(ioRequest)

	ltTickers := strings.Split(loParameters["ticker"], ",")

	lvTicker := ""

	if len(ltTickers) > 0 {
		lvTicker = ltTickers[0]
	}

	loAnalyzer := tinvestanalyser.Analyser{}

	loAnalyzer.Init(lvBearerToken[7:])

	ltProfit, loError := loAnalyzer.GetProfit(lvTicker, time.Now().AddDate(-10, 0, 0), time.Now())

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=UTF-8")

	if len(ltProfit) == 0 {
		ltProfit = make([]tinvestanalyser.Profit, 0)
	}

	lvBody, loError := json.Marshal(ltProfit)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=UTF-8")

	ioResponse.Write(lvBody)

}

func returnSignal(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	lvBearerToken := ioRequest.Header.Get("Authorization")

	if lvBearerToken == "" {
		ioResponse.WriteHeader(http.StatusUnauthorized)
		return
	}

	loParameters := mux.Vars(ioRequest)

	ltTickers := strings.Split(loParameters["ticker"], ",")

	loAnalyzer := tinvestanalyser.Analyser{}

	loAnalyzer.Init(lvBearerToken[7:])

	ltSignals, loError := loAnalyzer.GetSignals(ltTickers)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		return
	}

	if len(ltSignals) == 0 {
		ltSignals = make([]tinvestanalyser.Signal, 0)
	}

	lvBody, loError := json.Marshal(ltSignals)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=UTF-8")

	ioResponse.Write(lvBody)

}

func main() {

	loRouter := mux.NewRouter()

	loRouter.HandleFunc("/", returnRoot).Methods("GET")
	loRouter.HandleFunc("/positions", returnPositions).Methods("GET")
	loRouter.HandleFunc("/profit/{ticker}", returnProfit).Methods("GET")
	loRouter.HandleFunc("/profit", returnProfit).Methods("GET")
	loRouter.HandleFunc("/signal/{ticker}", returnSignal).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", loRouter))

}
