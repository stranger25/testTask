package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testTask/internal/httphandler"
	"testTask/internal/logerr"
	"testTask/internal/model"
	"testTask/internal/repository"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var Repo *repository.Repository
var HttpHandler *httphandler.HTTPClient
var Cfg model.Config

// -------------------------------------------------------------------------------------------
func init() {
	println("[ INF ]", "start init params...")

	file, err := os.Open("../helm/config.json")
	if err != nil {
		log.Fatal("[ ERR ]", "open config file", err)
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Cfg)
	if err != nil {
		log.Fatal("[ ERR ]", "decode config", err)
	}

	eLog := logerr.NewLogerr("../log/error.log")

	db, err := repository.InitDataBase(Cfg.Database.Connection)
	if err != nil {
		log.Fatal("[ ERR ]", err)
	}
	Repo = repository.NewRepository(db, eLog)

	client := httphandler.InitHTTPClient()
	HttpHandler = httphandler.NewClient("https://nationalbank.kz/rss/get_rates.cfm?fdate=", client)

	println("[ INF ]", "start init params...sucsess")
}

// -------------------------------------------------------------------------------------------
func main() {
	route := mux.NewRouter()
	route.Use(commonMiddleware)

	route.HandleFunc("/currency/save/{date}", downloadCurr).Methods("GET")
	s := route.PathPrefix("/currency").Subrouter()
	s.HandleFunc("/{date}", uploadCurr).Methods("GET")
	s.HandleFunc("/{date}/{code}", uploadCurr).Methods("GET")
	http.Handle("/", route)

	println("[ INF ]", "start server on port:", Cfg.Server.Port)

	err := http.ListenAndServe(":"+Cfg.Server.Port, route)
	if err != nil {
		log.Fatal("[ ERR ]", err.Error())
	}

}

// -------------------------------------------------------------------------------------------
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept", "*/*")
		w.Header().Add("Content-Type", "application/json;charset=UTF-8")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		w.Header().Set("Access-Control-Request-Headers", "content-type")
		next.ServeHTTP(w, r)
	})
}

// -------------------------------------------------------------------------------------------
func downloadCurr(w http.ResponseWriter, r *http.Request) {
	var ResOK model.Success
	var ResErr model.UnSuccess

	ResOK.Success = true

	date := mux.Vars(r)["date"]
	println("[ INF ]", "method downloadCurr is run", date)

	_, err := time.Parse("02.01.2006", date)
	if err != nil {
		println("[ ERR ] convert date" + err.Error())
		ResErr.Errmsg = "Ошибка даты в запросе"
		json.NewEncoder(w).Encode(&ResErr)
		return
	}
	rates, err := HttpHandler.GetData(date)
	if err != nil {
		println("[ ERR ]", "prepare request", err.Error())
		ResErr.Errmsg = "Ошибка запроса к сервису нац.банка"
		json.NewEncoder(w).Encode(&ResErr)
		return
	}

	go Repo.AddRate(rates)
	println("[ INF ]", "send answer to user")

	json.NewEncoder(w).Encode(&ResOK)

}

// -------------------------------------------------------------------------------------------
func uploadCurr(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]
	code := mux.Vars(r)["code"]
	var ResErr model.UnSuccess

	a_date, err := time.Parse("02.01.2006", date)
	if err != nil {
		println("[ ERR ] convert date" + err.Error())
		ResErr.Errmsg = "Ошибка даты в запросе"
		json.NewEncoder(w).Encode(&ResErr)
		return
	}

	res, err := Repo.GetRate(a_date, code)
	if err != nil {
		ResErr.Errmsg = "Ошибка чтения данных"
		json.NewEncoder(w).Encode(&ResErr)
		return
	}

	json.NewEncoder(w).Encode(res)
}
