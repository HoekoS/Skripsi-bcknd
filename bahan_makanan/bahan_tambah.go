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

var bahanMakananTambahModel *model_struct.BahanMakananTambah = new(model_struct.BahanMakananTambah)

var resultBahanMakananTambahModel []model_struct.BahanMakananTambah

type resultBahanMakananTambahStatus struct {
	Succes     string                            `json:"status"`
	StatusCode string                            `json:"status code"`
	Message    string                            `json:"message"`
	Data       []model_struct.BahanMakananTambah `json:"data"`
}

func queryProcesBahanMakananTambah(dbPool *pgxpool.Pool, query string) ([]model_struct.BahanMakananTambah, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&bahanMakananTambahModel.Id, &bahanMakananTambahModel.IdBahan, &bahanMakananTambahModel.Quantity, &bahanMakananTambahModel.CreateAt)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultBahanMakananTambahModel = append(resultBahanMakananTambahModel, *bahanMakananTambahModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultBahanMakananTambahModel, err
}

func dataToResultBahanTambah(resultModel []model_struct.BahanMakananTambah, message string) ([]byte, error) {
	// log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resultBahanMakananTambahStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultBahanMakananTambahModel,
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

func GetOrGetIdProsesBahanTambah(w http.ResponseWriter, r *http.Request, key string, ctg string, queryForProces string) {
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
		queryProcesBahanMakananTambah(dbPool, queryGet)
		fmt.Println(queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResultBahanTambah(resultBahanMakananTambahModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func BahanMakananTambahGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultBahanMakananTambahModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryBahanMakananTambah
	GetOrGetIdProsesBahanTambah(w, r, key, "", queryForProces)
}

func BahanMakananTambahPost(w http.ResponseWriter, r *http.Request) {
	resultBahanMakananTambahModel = nil
	var success string
	function.Cors(&w, r)
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&bahanMakananTambahModel); err != nil {
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
		queryPost := query.QueryBahanMakananTambahInsert + "VALUES ( '" + key + "' , " + strconv.Itoa(bahanMakananTambahModel.Quantity) + " , '" + formatted + "') RETURNING * ;"
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProcesBahanMakananTambah(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		// fmt.Println(resultBahanMakananTambahModel)
		message := "The data success to save !!"
		result, err := dataToResultBahanTambah(resultBahanMakananTambahModel, message)
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

func BahanMakananTambahEdit(w http.ResponseWriter, r *http.Request) {
	resultBahanMakananTambahModel = nil
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

		queryPost := query.QueryBahanMakananTambahUpdate + queryUpd
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProcesBahanMakananTambah(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		message := "The data success to save !!"
		result, err := dataToResultBahanTambah(resultBahanMakananTambahModel, message)
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

func BahanMakananTambahDelete(w http.ResponseWriter, r *http.Request) {
	resultBahanMakananTambahModel = nil
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

		queryGet := query.QueryBahanMakananTambahDelete + key + "' Returning *"
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

func BahanMakananTambah(key string, formatted string, quantity string) (succes string, err error) {
	// Connect to DB
	dbPool, err := databases.Connect()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	// QUERY PROSES
	queryPost := query.QueryBahanMakananTambahInsert + "VALUES ( '" + key + "' , " + quantity + " , '" + formatted + "') RETURNING * ;"
	fmt.Println(queryPost)
	// fmt.Println(err)
	queryProcesBahanMakananTambah(dbPool, queryPost)
	if err != nil {
		succes = "false"
		return succes, err
	}
	succes = "true"
	return succes, err
}
