package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ivangurin/tinvest-analyser-go"
	"github.com/ivangurin/tinvest-client-go"
	"log"
	"net/http"
	"strings"
	"time"
)

func returnRoot(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	ioResponse.WriteHeader(http.StatusBadRequest)

}

func returnPositions(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	fmt.Print("\n", time.Now().Format(time.RFC3339), ioRequest.Method, " Positions were requested - ")

	ioResponse.Header().Add("Access-Control-Allow-Origin", "*")
	ioResponse.Header().Add("Access-Control-Allow-Methods", "*")
	ioResponse.Header().Add("Access-Control-Allow-Headers", "*")

	if ioRequest.Method == http.MethodOptions {
		return
	}

	lvBearerToken := ioRequest.Header.Get("Authorization")

	if lvBearerToken == "" {
		ioResponse.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(ioResponse, http.StatusText(http.StatusUnauthorized))
		fmt.Print(http.StatusText(http.StatusUnauthorized))
		return
	}

	loClient := tinvestclient.Client{}

	loClient.Init(lvBearerToken[7:])

	ltPositions, loError := loClient.GetPositions()

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Print(loError)
		return
	}

	if len(ltPositions) == 0 {
		ltPositions = make([]tinvestclient.Position, 0)
	}

	lvBody, loError := json.Marshal(ltPositions)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Print(loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=utf-8")

	ioResponse.Write(lvBody)

	fmt.Print("OK")

}

func returnProfit(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	fmt.Print("\n", time.Now().Format(time.RFC3339), ioRequest.Method, " Profit was requested - ")

	ioResponse.Header().Add("Access-Control-Allow-Origin", "*")
	ioResponse.Header().Add("Access-Control-Allow-Methods", "*")
	ioResponse.Header().Add("Access-Control-Allow-Headers", "*")

	if ioRequest.Method == http.MethodOptions {
		return
	}

	lvBearerToken := ioRequest.Header.Get("Authorization")

	if lvBearerToken == "" {
		ioResponse.WriteHeader(http.StatusUnauthorized)
		fmt.Print(http.StatusText(http.StatusUnauthorized))
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
		fmt.Print(loError)
		return
	}

	ioResponse.Header().Add("Access-Control-Allow-Origin", "true")

	if len(ltProfit) == 0 {
		ltProfit = make([]tinvestanalyser.Profit, 0)
	}

	lvBody, loError := json.Marshal(ltProfit)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Print(loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=utf-8")

	ioResponse.Write(lvBody)

	fmt.Print("OK")

}

func returnSignal(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	fmt.Print("\n", time.Now().Format(time.RFC3339), ioRequest.Method, " Signal was requested - ")

	ioResponse.Header().Add("Access-Control-Allow-Origin", "*")
	ioResponse.Header().Add("Access-Control-Allow-Methods", "*")
	ioResponse.Header().Add("Access-Control-Allow-Headers", "*")

	if ioRequest.Method == http.MethodOptions {
		return
	}

	lvBearerToken := ioRequest.Header.Get("Authorization")

	if lvBearerToken == "" {
		ioResponse.WriteHeader(http.StatusUnauthorized)
		fmt.Print(http.StatusText(http.StatusUnauthorized))
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
		fmt.Print(loError)
		return
	}

	if len(ltSignals) == 0 {
		ltSignals = make([]tinvestanalyser.Signal, 0)
	}

	lvBody, loError := json.Marshal(ltSignals)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Print(loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=utf-8")

	ioResponse.Write(lvBody)

	fmt.Print("OK")

}

func main() {

	loRouter := mux.NewRouter()

	loRouter.HandleFunc("/", returnRoot).Methods(http.MethodGet)
	loRouter.HandleFunc("/positions", returnPositions).Methods(http.MethodGet, http.MethodOptions)
	loRouter.HandleFunc("/profit", returnProfit).Methods(http.MethodGet, http.MethodOptions)
	loRouter.HandleFunc("/signal/{ticker}", returnSignal).Methods(http.MethodGet, http.MethodOptions)

	log.Fatal(http.ListenAndServe(":8081", loRouter))

}
