package qr_code

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

var qrModel *model_struct.QrModel = new(model_struct.QrModel)

var qrNpesanan *model_struct.QrNPesanan = new(model_struct.QrNPesanan)

var resultModel []model_struct.QrModel

type resultStatus struct {
	Succes     string                 `json:"status"`
	StatusCode string                 `json:"status code"`
	Message    string                 `json:"message"`
	Data       []model_struct.QrModel `json:"data"`
}

func queryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.QrModel, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&qrModel.Id, &qrModel.NomerNota, &qrModel.NomerMeja)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultModel = append(resultModel, *qrModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultModel, err
}

func dataToResult(resultModel []model_struct.QrModel, message string) ([]byte, error) {
	// log.Println("test")
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultModel,
		}
		result, err := json.Marshal(jsondat)
		// fmt.Println(jsondat)
		if err != nil {
			log.Println("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &resultStatus{
			Succes:     "False",
			StatusCode: "404",
			Message:    "Your Insert Data is Not Complete",
			Data:       resultModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Println("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func GetOrGetIdProses(w http.ResponseWriter, r *http.Request, key string, queryForProces string) {
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
			queryGet = queryForProces + " WHERE ID = '" + key + "'"
		} else {
			log.Print("No Key")
			queryGet = queryForProces
		}
		queryProces(dbPool, queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResult(resultModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func QrGetAll(w http.ResponseWriter, r *http.Request) {
	var key string
	resultModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryQr
	GetOrGetIdProses(w, r, key, queryForProces)
}

func QrgetOne(w http.ResponseWriter, r *http.Request) {
	var key string
	resultModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryQrGetOne
	GetOrGetIdProses(w, r, key, queryForProces)
}

func QrPost(w http.ResponseWriter, r *http.Request) {
	resultModel = nil
	var err error
	function.Cors(&w, r)

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&qrNpesanan); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("error disini")
			return
		}

		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		uuid := uuid_generator.UuidGenerate()
		formatted := function.FormatTime()

		passQueryPost := `Values ('` + uuid.String() + `','Pilih Menu','` + qrNpesanan.NamaPelanggan + `','` + formatted + `',` + strconv.Itoa(qrNpesanan.FlagTakeAway) + `)`

		queryPost := query.QrNPesananInsert(passQueryPost, uuid.String(), strconv.Itoa(qrNpesanan.NomerMeja))

		// rows, err := dbPool.Query(context.Background(), query)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	fmt.Println("error disini")
		// }

		// queryPost := query.QueryQrInsert + "VALUES ( '" + qrNpesanan.NomerNota + "' , " + strconv.Itoa(qrNpesanan.NomerMeja) + ") RETURNING * ;"

		fmt.Println(queryPost)
		// queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		message := "The data success to save !!"
		result, err := dataToResult(resultModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept POST request", http.StatusBadRequest)
	}
}

func QrUpdate(w http.ResponseWriter, r *http.Request) {
	resultModel = nil
	var err error

	if r.Method == "PUT" {
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&qrModel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()
		key, err := function.GetId(r)
		if err != nil {
			log.Fatal("error while get Id")
		}

		queryUpd := "SET no_nota= '" + qrModel.NomerNota + "', no_meja= " + strconv.Itoa(qrModel.NomerMeja) + "WHERE id = " + key + " returning *;"

		queryUpdate := query.QueryQrUpdate + queryUpd
		fmt.Println(queryUpdate)
		queryProces(dbPool, queryUpdate)
		if err != nil {
			log.Fatal("error while executing query")
		}

		message := "The data success to save !!"
		result, err := dataToResult(resultModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept PUT request", http.StatusBadRequest)
	}
}
func QrDelete(w http.ResponseWriter, r *http.Request) {
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

		queryGet := query.QueryQrDelete + key + " Returning *"
		queryProces(dbPool, queryGet)
		if err != nil {
			log.Fatal("error while executing query")
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept DELETE request", http.StatusBadRequest)
	}
}
