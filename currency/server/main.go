package main

import (
	"database/sql"
	"desafio-padawan-go/currency/converter"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Conversor struct {
	ID              int64     `json:"id"`
	FromCurrency    string    `json:"fromCurrency"`
	ToCurrency      string    `json:"toCurrency"`
	Rate            float64   `json:"rate"`
	Amount          float64   `json:"amount"`
	ConvertedAmount float64   `json:"convertedAmount"`
	CreationDate    time.Time `json:"createdAt"`
}

var db *sql.DB

func main() {
	var err error
	db, err := sql.Open("mysql", "guilherme:secretkey@tcp(localhost:3306)/converter_db")
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	println("Database connection successfully established!")

	r := mux.NewRouter()
	r.HandleFunc("/exchange/{amount}/{from}/{to}/{rate}", handleExchangeRequest).Methods("CREATE")
	r.HandleFunc("/exchange/{id}", getConversion).Methods("GET")
	println("Running server on 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func handleExchangeRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	amount, _ := strconv.ParseFloat(params["amount"], 64)
	from := params["from"]
	to := params["to"]
	rate, _ := strconv.ParseFloat(params["rate"], 64)
	response := converter.Convert(amount, from, to, rate)

	var conversor Conversor

	ins, err := db.Prepare("INSERT INTO converter(amount, from, to, rate, response) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer ins.Close()

	result, err := ins.Exec(conversor.FromCurrency, conversor.ToCurrency, conversor.Rate, conversor.Amount, conversor.ConvertedAmount, conversor.CreationDate)
	if err != nil {
		panic(err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	println(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getConversion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, from_currency, to_currency, rate, amount, converted_amount, creation_date FROM conversor WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var conversor Conversor
	if result.Next() {
		err := result.Scan(&conversor.ID, &conversor.FromCurrency, &conversor.ToCurrency, &conversor.Rate, &conversor.Amount, &conversor.ConvertedAmount, &conversor)
		if err != nil {
			panic(err.Error())
		}
	} else {
		http.Error(w, "Conversion not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversor)
}
