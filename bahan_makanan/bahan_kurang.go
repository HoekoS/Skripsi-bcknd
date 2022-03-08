package bahan_makanan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"skripsi.com/backend-adm/databases"
	"skripsi.com/backend-adm/function"
	"skripsi.com/backend-adm/model_struct"
	"skripsi.com/backend-adm/query"
)

var bahanMakananKurangModel *model_struct.BahanMakananKurang = new(model_struct.BahanMakananKurang)

var resultBahanMakananKurangModel []model_struct.BahanMakananKurang

type resultBahanMakananKurangStatus struct {
	Succes     string                            `json:"status"`
	StatusCode string                            `json:"status code"`
	Message    string                            `json:"message"`
	Data       []model_struct.BahanMakananKurang `json:"data"`
}

func queryProcesBahanMakananKurang(dbPool *pgxpool.Pool, query string) ([]model_struct.BahanMakananKurang, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&bahanMakananKurangModel.Id, &bahanMakananKurangModel.IdBahan, &bahanMakananKurangModel.Quantity, &bahanMakananKurangModel.CreateAt)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultBahanMakananKurangModel = append(resultBahanMakananKurangModel, *bahanMakananKurangModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultBahanMakananKurangModel, err
}

func dataToResultBahanKurang(resultModel []model_struct.BahanMakananKurang, message string) ([]byte, error) {
	// log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resultBahanMakananKurangStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultBahanMakananKurangModel,
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

func GetOrGetIdProsesBahanKurang(w http.ResponseWriter, r *http.Request, key string, ctg string, queryForProces string) {
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
			queryGet = queryForProces + " WHERE id_bahan_makanan = '" + key + "'"
		} else {
			log.Print("No Key")
			queryGet = queryForProces
		}
		queryGet += " ORDER BY create_at ASC"
		queryProcesBahanMakananKurang(dbPool, queryGet)
		fmt.Println(queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResultBahanKurang(resultBahanMakananKurangModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func BahanMakananKurangGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultBahanMakananKurangModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryBahanMakananKurang
	GetOrGetIdProsesBahanKurang(w, r, key, "", queryForProces)
}

func BahanMakananKurangPost(w http.ResponseWriter, r *http.Request) {
	resultBahanMakananKurangModel = nil
	var success string
	function.Cors(&w, r)
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&bahanMakananKurangModel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Connect to DB
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		key := keys[0]
		// GET FORMAT TIME DD-MM-YYYY
		formatted := function.FormatTime()
		// QUERY PROSES
		queryPost := query.QueryBahanMakananKurangInsert + "VALUES ( '" + key + "' , " + strconv.Itoa(bahanMakananKurangModel.Quantity) + " , '" + formatted + "') RETURNING * ;"
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProcesBahanMakananKurang(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		// fmt.Println(resultBahanMakananKurangModel)
		message := "The data success to save !!"
		result, err := dataToResultBahanKurang(resultBahanMakananKurangModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		success, err = BahanMakananUpdateQuantity(key, formatted)
		fmt.Println(success)
		if err != nil {
			log.Fatal("error tambah makanan", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept POST request", http.StatusBadRequest)
	}
}

func BahanMakananKurangEdit(w http.ResponseWriter, r *http.Request) {
	resultBahanMakananKurangModel = nil
	function.Cors(&w, r)
	if r.Method == "PUT" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&bahanMakananModel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Connect to DB
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		key := keys[0]

		// GET FORMAT TIME DD-MM-YYYY
		formatted := function.FormatTime()
		// QUERY PROSES
		queryUpd := "SET quantity= " + strconv.Itoa(bahanMakananModel.Quantity) + ", create_at='" + formatted + "' WHERE id = " + key + " returning *;"

		queryPost := query.QueryBahanMakananKurangUpdate + queryUpd
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProcesBahanMakananKurang(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		message := "The data success to save !!"
		result, err := dataToResultBahanKurang(resultBahanMakananKurangModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept PUT request", http.StatusBadRequest)
	}
}

func BahanMakananKurangDelete(w http.ResponseWriter, r *http.Request) {
	resultBahanMakananKurangModel = nil
	function.Cors(&w, r)
	if r.Method == "DELETE" {

		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		key := keys[0]

		queryGet := query.QueryBahanMakananKurangDelete + key + "' Returning *"
		queryProces(dbPool, queryGet)
		if err != nil {
			log.Println("error while executing query")
		}
		fmt.Println(queryGet)

		jsondat := &resultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    "You'r data succes to delete !",
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", "*")
		// w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept DELETE request", http.StatusBadRequest)
	}
}
