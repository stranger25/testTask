package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"testTask/structs"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var Db *sql.DB
var Cfg structs.Config

//-------------------------------------------------------------------------------------------
func init() {
	println("[ INF ]", "start init params...")

	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal("[ ERR ]", "open config file", err.Error())
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Cfg)
	if err != nil {
		log.Fatal("[ ERR ]", "decode config", err.Error())
	}

	Db, err = sql.Open("mysql", Cfg.Database.Connection)
	if err != nil {
		log.Fatal("[ ERR ]", "open DB connection", err.Error())
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal("[ ERR ]", "ping database", err.Error())
	}

	println("[ INF ]", "start init params...sucsess")
}

//-------------------------------------------------------------------------------------------
func main() {
	route := mux.NewRouter()
	route.Use(commonMiddleware)

	route.HandleFunc("/currency/save/{date}", downloadCurr).Methods("GET")
	route.HandleFunc("/currency/{date}/{*code}", uploadCurr).Methods("GET")
	http.Handle("/", route)

	println("[ INF ]", "start server on port:", Cfg.Server.Port)

	err := http.ListenAndServe(":"+Cfg.Server.Port, route)
	if err != nil {
		log.Fatal("[ ERR ]", err.Error())
	}

}

//-------------------------------------------------------------------------------------------
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

//-------------------------------------------------------------------------------------------
func downloadCurr(w http.ResponseWriter, r *http.Request) {
	var ResOK structs.Success
	var ResErr structs.UnSuccess

	ResOK.Success = true

	date := mux.Vars(r)["date"]
	println("[ INF ]", "method downloadCurr is run", date)

	req, err := http.NewRequest("GET", "https://nationalbank.kz/rss/get_rates.cfm?fdate="+date, nil)
	if err != nil {
		println("[ ERR ]", "prepare request", err.Error())
		ResErr.Errmsg = "Ошибка формирования запроса к платформе нац.банка"
		json.NewEncoder(w).Encode(&ResErr)
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		println("[ ERR ]", "do request", err.Error())
		ResErr.Errmsg = "Ошибка отправки запроса к платформе нац.банка"
		json.NewEncoder(w).Encode(&ResErr)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		println("[ ERR ]", "read response", err.Error())
		ResErr.Errmsg = "Ошибка чтения ответа от платформы нац.банка"
		json.NewEncoder(w).Encode(&ResErr)
		return
	}

	var rates structs.Rates

	err = xml.Unmarshal(body, &rates)
	if err != nil {
		println("[ ERR ]", "Unmarshal xml", err.Error())
		ResErr.Errmsg = "Ошибка разбора документа от платформы нац.банка"
		json.NewEncoder(w).Encode(&ResErr)
		return
	}
	go writeDB(&rates)
	println("[ INF ]", "send answer to user")

	json.NewEncoder(w).Encode(&ResOK)

}

//-------------------------------------------------------------------------------------------
func uploadCurr(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]
	code := mux.Vars(r)["code"]
	var ResErr structs.UnSuccess

	a_date, err := time.Parse("02.01.2006", date)
	if err != nil {
		ErrLog("[ ERR ] convert date" + err.Error())
		ResErr.Errmsg = "Ошибка преобразования даты в запросе"
		return
	}

	var rows *sql.Rows

	if len(code) != 0 {
		rows, err = Db.Query("select title, code, value, a_date from R_CURRENCY where a_date = ? and code = ?", a_date, code)
		if err != nil {
			ErrLog("[ ERR ] convert date" + err.Error())
			ResErr.Errmsg = "Ошибка преобразования даты в запросе"
			json.NewEncoder(w).Encode(&ResErr)
			return
		}
	} else {
		rows, err = Db.Query("select title, code, value, a_date from R_CURRENCY where a_date = ?", a_date)
		if err != nil {
			ErrLog("[ ERR ] convert date" + err.Error())
			ResErr.Errmsg = "Ошибка преобразования даты в запросе"
			json.NewEncoder(w).Encode(&ResErr)
			return
		}
	}
	println("[ INF ]", "method uploadCurr is run", date, code)

	var RatesOut []structs.RatesOut

	for rows.Next() {
		var Rates structs.RatesOut
		rows.Scan(&Rates.Title, &Rates.Code, &Rates.Value, &Rates.Date)
		if err != nil {
			ErrLog("[ ERR ] read result from DS" + err.Error())
			ResErr.Errmsg = "Ошибка чтения жанных запроса"
			json.NewEncoder(w).Encode(&ResErr)
			break
		}
		RatesOut = append(RatesOut, Rates)
	}
	json.NewEncoder(w).Encode(&RatesOut)
}

//-------------------------------------------------------------------------------------------
func writeDB(data *structs.Rates) {

	a_date, err := time.Parse("02.01.2006", data.Date)
	if err != nil {
		ErrLog("[ ERR ] convert date" + err.Error())
		return
	}
	for i, item := range data.Item {

		value, err := strconv.ParseFloat(item.Description, 2)
		if err != nil {
			ErrLog("[ ERR ] convert to float" + err.Error())
		} else {
			println("[ INF ]", "write item", i, "to database")
			_, err := Db.Exec("insert into R_CURRENCY(title, code, value, a_date) values (?, ?, ?, ?)", item.Fullname, item.Title, value, a_date)
			if err != nil {
				ErrLog("[ ERR ] insert in database :" + " data item.Title, value: " + item.Title + " " + err.Error())
			}
		}
	}
}

//-------------------------------------------------------------------------------------------
func Exist(file_path string) bool {
	_, err := os.Stat(file_path)
	return !os.IsNotExist(err)
}

//-------------------------------------------------------------------------------------------
func ErrLog(s string) {
	var f *os.File
	var err error
	var file_path = "./log/error.log"

	if Exist(file_path) {
		f, err = os.OpenFile(file_path, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			println("[ ERR ]", "open log file :", err.Error())
			return
		}
	} else {
		f, err = os.OpenFile(file_path, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			println("[ ERR ]", "create log file :", err.Error())
			return
		}
	}
	_, err = io.WriteString(f, s+"\n")
	if err != nil {
		println("[ ERR ]", "write to log file :", err.Error())
		return
	}
}
