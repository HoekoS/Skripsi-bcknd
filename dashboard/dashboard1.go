package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"skripsi.com/backend-adm/databases"
	"skripsi.com/backend-adm/function"
	"skripsi.com/backend-adm/model_struct"
	"skripsi.com/backend-adm/query"
)

// var err error

var jumlah *model_struct.JumlahDashboard = new(model_struct.JumlahDashboard)

var resultJumlah []model_struct.JumlahDashboard

type resultJumlahStatus struct {
	Succes     string                         `json:"status"`
	StatusCode string                         `json:"status code"`
	Message    string                         `json:"message"`
	Data       []model_struct.JumlahDashboard `json:"data"`
}

type failedStatus struct {
	Succes     string `json:"status"`
	StatusCode string `json:"status code"`
	Message    string `json:"message"`
	Data       string `json:"data"`
}

func queryProcesJumlah(dbPool *pgxpool.Pool, query string) ([]model_struct.JumlahDashboard, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&jumlah.Jumlah)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultJumlah = append(resultJumlah, *jumlah)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultJumlah, err
}

func dataToResultJumlah(resultModel []model_struct.JumlahDashboard, message string) ([]byte, error) {
	// log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resultJumlahStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultJumlah,
		}
		result, err := json.Marshal(jsondat)
		// fmt.Println(jsondat)
		if err != nil {
			log.Println("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &failedStatus{
			Succes:     "False",
			StatusCode: "404",
			Message:    "Your Insert Data is Not Complete",
			Data:       "",
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Println("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func GetOrGetIdProses(w http.ResponseWriter, r *http.Request, key string, ctg string, queryForProces string) {
	var queryGet string
	if r.Method == "GET" {
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		if key != "" {
			log.Print("With Key")
			queryGet = queryForProces + " WHERE id = '" + key + "'"
		} else {
			log.Print("No Key")
			queryGet = queryForProces
		}
		queryProcesJumlah(dbPool, queryGet)
		fmt.Println(queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResultJumlah(resultJumlah, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func JmlKosongDashboardGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultJumlah = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryDashboardKosong
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func JmlSelesaiDashboardGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultJumlah = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryDashboardSelesai
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func JmlBatalDashboardGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultJumlah = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryDashboardBatal
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func JmlHabisDashboardGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultJumlah = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryJmlhBahanHabis
	GetOrGetIdProses(w, r, key, "", queryForProces)
}
