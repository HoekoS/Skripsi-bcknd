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
	"skripsi.com/backend-adm/uuid_generator"
)

var err error

var bahanMakananModel *model_struct.BahanMakanan = new(model_struct.BahanMakanan)

var resultbahanMakananModel []model_struct.BahanMakanan

type resultStatus struct {
	Succes     string                      `json:"status"`
	StatusCode string                      `json:"status code"`
	Message    string                      `json:"message"`
	Data       []model_struct.BahanMakanan `json:"data"`
}

type failedStatus struct {
	Succes     string `json:"status"`
	StatusCode string `json:"status code"`
	Message    string `json:"message"`
	Data       string `json:"data"`
}

func queryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.BahanMakanan, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&bahanMakananModel.Id, &bahanMakananModel.Description, &bahanMakananModel.CreateAt, &bahanMakananModel.Quantity, &bahanMakananModel.Satuan)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultbahanMakananModel = append(resultbahanMakananModel, *bahanMakananModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultbahanMakananModel, err
}

func dataToResult(resultModel []model_struct.BahanMakanan, message string) ([]byte, error) {
	log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultbahanMakananModel,
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
	var queryGet, success string
	if r.Method == "GET" {
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		if key != "" {
			log.Print("With Key")
			queryGet = queryForProces + " WHERE ID = '" + key + "'"
			formatted := ""
			success, err = BahanMakananUpdateQuantity(key, formatted)
			fmt.Println(success)
			if err != nil {
				log.Fatal("error tambah makanan", err)
			}
		} else {
			log.Print("No Key")
			queryGet = queryForProces + " ORDER BY description DESC"
		}
		queryProces(dbPool, queryGet)
		fmt.Println(queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResult(resultbahanMakananModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func BahanMakananGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultbahanMakananModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryBahanMakanan
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func BahanMakananPost(w http.ResponseWriter, r *http.Request) {
	resultbahanMakananModel = nil
	var success string
	function.Cors(&w, r)
	if r.Method == "POST" {
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

		// GET UUID
		uuid := uuid_generator.UuidGenerate()

		// GET FORMAT TIME DD-MM-YYYY
		formatted := function.FormatTime()
		// QUERY PROSES
		queryPost := query.QueryBahanMakananInsert + "VALUES ( '" + uuid.String() + "' , '" + bahanMakananModel.Description + "' , '" + formatted + "', '" + bahanMakananModel.Satuan + "') RETURNING * ;"
		success, err = BahanMakananTambah(uuid.String(), formatted, strconv.Itoa(bahanMakananModel.Quantity))
		fmt.Println(success)
		if err != nil {
			log.Fatal("error tambah makanan", err)
		}
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		// fmt.Println(resultbahanMakananModel)
		message := "The data success to save !!"
		result, err := dataToResult(resultbahanMakananModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept POST request", http.StatusBadRequest)
	}
}

func BahanMakananEdit(w http.ResponseWriter, r *http.Request) {
	resultbahanMakananModel = nil
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
		queryUpd := "SET description= '" + bahanMakananModel.Description + "', quantity=" + strconv.Itoa(bahanMakananModel.Quantity) + ", satuan = '" + bahanMakananModel.Satuan + "', create_at = '" + formatted + "' WHERE id = '" + key + "' returning *;"

		queryPost := query.QueryBahanMakananUpdate + queryUpd
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		message := "The data success to save !!"
		result, err := dataToResult(resultbahanMakananModel, message)
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

func BahanMakananDelete(w http.ResponseWriter, r *http.Request) {
	resultbahanMakananModel = nil
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

		queryGet := query.QueryBahanMakananDelete + key + "' Returning *"
		queryProces(dbPool, queryGet)
		if err != nil {
			log.Println("error while executing query")
		}
		fmt.Println(queryGet)

		queryDeleteTambah := query.QueryBahanMakananTambahDeleteWithId + key + "'"
		queryProces(dbPool, queryDeleteTambah)
		if err != nil {
			log.Println("error while executing query")
		}

		queryDeleteKUrang := query.QueryBahanMakananKurangDeleteWithId + key + "'"
		queryProces(dbPool, queryDeleteKUrang)
		if err != nil {
			log.Println("error while executing query")
		}
		// fmt.Println("sini"+queryDeleteTambah, queryDeleteKUrang)

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

func BahanMakananUpdateQuantity(key string, date string) (succes string, err error) {
	dbPool, err := databases.Connect()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	// QUERY PROSES
	queryPost := query.QueryBahanMakananUPdateQuantity(key, date)
	fmt.Println(queryPost)
	// fmt.Println(err)
	queryProces(dbPool, queryPost)
	if err != nil {
		succes = "false"
		return succes, err
	}
	succes = "true"
	return succes, err
}

func BahanMakananDashGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultbahanMakananModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryBahanDash
	GetOrGetIdProses(w, r, key, "", queryForProces)
}
