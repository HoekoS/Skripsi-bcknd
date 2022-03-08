package employee

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

var employeeScheduleModel *model_struct.EmployeeSchedule = new(model_struct.EmployeeSchedule)

var resultEmployeeScheduleModel []model_struct.EmployeeSchedule

type employeeScheduleResultStatus struct {
	Succes     string                          `json:"status"`
	StatusCode string                          `json:"status code"`
	Message    string                          `json:"message"`
	Data       []model_struct.EmployeeSchedule `json:"data"`
}

func employeeScheduleQueryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.EmployeeSchedule, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(
			&employeeScheduleModel.Id,
			&employeeScheduleModel.EmployeeId,
			&employeeScheduleModel.ScheduleId,
		)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultEmployeeScheduleModel = append(resultEmployeeScheduleModel, *employeeScheduleModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultEmployeeScheduleModel, err
}

func dataToResultEmployeeSchedule(resultModel []model_struct.EmployeeSchedule, message string) ([]byte, error) {
	log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &employeeScheduleResultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultEmployeeScheduleModel,
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

func employeeScheduleGetOrGetIdProses(w http.ResponseWriter, r *http.Request, key string, ctg string, queryForProces string) {
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
		employeeScheduleQueryProces(dbPool, queryGet)
		fmt.Println(queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResultEmployeeSchedule(resultEmployeeScheduleModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func EmployeeScheduleGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultEmployeeScheduleModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryEmployeeSchedule
	employeeScheduleGetOrGetIdProses(w, r, key, "", queryForProces)
}

func EmployeeSchedulePost(w http.ResponseWriter, r *http.Request) {
	resultEmployeeScheduleModel = nil
	function.Cors(&w, r)
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&employeeScheduleModel); err != nil {
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

		// QUERY PROSES
		queryPost := query.QueryEmployeeScheduleInsert + "VALUES ( '" + employeeScheduleModel.EmployeeId + "' , " + strconv.Itoa(employeeScheduleModel.ScheduleId) + " ) RETURNING * ;"

		fmt.Println(queryPost)
		// fmt.Println(err)
		employeeScheduleQueryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		// fmt.Println(resultEmployeeScheduleModel)
		message := "The data success to save !!"
		result, err := dataToResultEmployeeSchedule(resultEmployeeScheduleModel, message)
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

func EmployeeScheduleEdit(w http.ResponseWriter, r *http.Request) {
	resultEmployeeScheduleModel = nil
	function.Cors(&w, r)
	if r.Method == "PUT" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&employeeScheduleModel); err != nil {
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

		// QUERY PROSES
		queryUpd := "SET employee_id= '" + employeeScheduleModel.EmployeeId + "', schedule_id=" + strconv.Itoa(employeeScheduleModel.ScheduleId) + "' WHERE id = '" + key + "' returning *;"

		queryPost := query.QueryEmployeeScheduleUpdate + queryUpd
		fmt.Println(queryPost)
		// fmt.Println(err)
		employeeScheduleQueryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		message := "The data success to save !!"
		result, err := dataToResultEmployeeSchedule(resultEmployeeScheduleModel, message)
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

func EmployeeScheduleDelete(w http.ResponseWriter, r *http.Request) {
	resultEmployeeScheduleModel = nil
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

		queryGet := query.QueryEmployeeScheduleDelete + key + " Returning *"
		employeeScheduleQueryProces(dbPool, queryGet)
		if err != nil {
			log.Println("error while executing query")
		}
		fmt.Println(queryGet)

		jsondat := &employeeScheduleResultStatus{
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
