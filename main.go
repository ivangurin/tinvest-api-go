package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ivangurin/tinvest-analyser-go"
	"github.com/ivangurin/tinvest-client-go"
	"log"
	"net/http"
	"net/url"
	"time"
)

func returnRoot(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	lvBody := "API for Tinkoff Analyser. See more at https://github.com/ivangurin/tinvest-service-go"

	ioResponse.Write([]byte(lvBody))

}

func returnPositions(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	fmt.Print("\n", time.Now().Format(time.RFC3339), " ", ioRequest.URL.Path, " - ")

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
		fmt.Println(http.StatusText(http.StatusUnauthorized))
		return
	}

	loClient := tinvestclient.Client{}

	loClient.Init(lvBearerToken[7:])

	ltPositions, loError := loClient.GetPositions()

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
		return
	}

	if len(ltPositions) == 0 {
		ltPositions = make([]tinvestclient.Position, 0)
	}

	lvBody, loError := json.Marshal(ltPositions)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=utf-8")

	ioResponse.Write(lvBody)

	fmt.Print("OK")

}

func returnOperations(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	fmt.Print("\n", time.Now().Format(time.RFC3339), " ", ioRequest.URL.Path, "/", ioRequest.URL.RawQuery, " - ")

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
		fmt.Println(http.StatusText(http.StatusUnauthorized))
		return
	}

	ltVars := mux.Vars(ioRequest)

	lvTicker := ltVars["ticker"]

	if lvTicker == ""{
		loError := errors.New("ticker is missing")
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
		return
	}

	loClient := tinvestclient.Client{}

	loClient.Init(lvBearerToken[7:])

	lsInstrument, loError := loClient.GetInstrumentByTicker(lvTicker)

	if loError != nil {
		loError := errors.New("instrument not found")
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
		return
	}

	ltOperations, loError := loClient.GetOperations(lsInstrument.FIGI, time.Now().AddDate(-10, 0, 0), time.Now())

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
		return
	}

	if len(ltOperations) == 0 {
		ltOperations = make([]tinvestclient.Operation, 0)
	}

	lvBody, loError := json.Marshal(ltOperations)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
		return
	}

	ioResponse.Header().Add("Content-Type", "application/json; charset=utf-8")

	ioResponse.Write(lvBody)

	fmt.Println("OK")

}

func returnProfit(ioResponse http.ResponseWriter, ioRequest *http.Request) {

	fmt.Print("\n", time.Now().Format(time.RFC3339), " ", ioRequest.URL.Path, " - ")

	ioResponse.Header().Add("Access-Control-Allow-Origin", "*")
	ioResponse.Header().Add("Access-Control-Allow-Methods", "*")
	ioResponse.Header().Add("Access-Control-Allow-Headers", "*")

	if ioRequest.Method == http.MethodOptions {
		return
	}

	lvBearerToken := ioRequest.Header.Get("Authorization")

	if lvBearerToken == "" {
		ioResponse.WriteHeader(http.StatusUnauthorized)
		fmt.Println(http.StatusText(http.StatusUnauthorized))
		return
	}

	loParameters, loError := url.ParseQuery(ioRequest.URL.RawQuery)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
		return
	}

	lvTicker := ""

	ltTickers, lvExists := loParameters["ticker"]

	if lvExists{
		if len(ltTickers) > 0 {
			lvTicker = ltTickers[0]
		}
	}

	loAnalyzer := tinvestanalyser.Analyser{}

	loAnalyzer.Init(lvBearerToken[7:])

	ltProfit, loError := loAnalyzer.GetProfit(lvTicker, time.Now().AddDate(-10, 0, 0), time.Now())

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
		return
	}

	if len(ltProfit) == 0 {
		ltProfit = make([]tinvestanalyser.Profit, 0)
	}

	lvBody, loError := json.Marshal(ltProfit)

	if loError != nil {
		ioResponse.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(ioResponse, "%v", loError)
		fmt.Println(loError)
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
	loRouter.HandleFunc("/operations/{ticker}", returnOperations).Methods(http.MethodGet, http.MethodOptions)
	loRouter.HandleFunc("/profit", returnProfit).Methods(http.MethodGet, http.MethodOptions)

	log.Fatal(http.ListenAndServe(":8081", loRouter))

}