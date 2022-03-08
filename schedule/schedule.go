package schedule

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
	"skripsi.com/backend-adm/uuid_generator"
)

var err error

var scheduleEmpViewModel *model_struct.ViewSchedule = new(model_struct.ViewSchedule)

var resultSchedulEmpVieweModel []model_struct.ViewSchedule

type resulScheduleEmptStatus struct {
	Succes     string                      `json:"status"`
	StatusCode string                      `json:"status code"`
	Message    string                      `json:"message"`
	Data       []model_struct.ViewSchedule `json:"data"`
}

var scheduleModel *model_struct.Schedule = new(model_struct.Schedule)

var resultScheduleModel []model_struct.Schedule

type resultStatus struct {
	Succes     string                  `json:"status"`
	StatusCode string                  `json:"status code"`
	Message    string                  `json:"message"`
	Data       []model_struct.Schedule `json:"data"`
}

type failedStatus struct {
	Succes     string `json:"status"`
	StatusCode string `json:"status code"`
	Message    string `json:"message"`
	Data       string `json:"data"`
}

func queryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.Schedule, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&scheduleModel.Id, &scheduleModel.Day, &scheduleModel.StartAt, &scheduleModel.EndAt, &scheduleModel.Description)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultScheduleModel = append(resultScheduleModel, *scheduleModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultScheduleModel, err
}

func dataToResult(resultModel []model_struct.Schedule, message string) ([]byte, error) {
	// log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultScheduleModel,
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

func queryScheduleEmpProces(dbPool *pgxpool.Pool, query string) ([]model_struct.ViewSchedule, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&scheduleEmpViewModel.FullName, &scheduleEmpViewModel.Day, &scheduleEmpViewModel.Description)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultSchedulEmpVieweModel = append(resultSchedulEmpVieweModel, *scheduleEmpViewModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultSchedulEmpVieweModel, err
}

func dataScheduleEmpToResult(resultModel []model_struct.ViewSchedule, message string) ([]byte, error) {
	// log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resulScheduleEmptStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultSchedulEmpVieweModel,
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

func GetOrGetIdProses(w http.ResponseWriter, r *http.Request, key string, flag int, queryForProces string) {
	var queryGet string
	var result []byte
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
		if flag == 0 {
			queryGet += " ORDER BY id ASC"
			queryProces(dbPool, queryGet)
			fmt.Println(queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			message := "You'r data succes to send !"
			result, err = dataToResult(resultScheduleModel, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			queryScheduleEmpProces(dbPool, queryGet)
			fmt.Println(queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			message := "You'r data succes to send !"
			result, err = dataScheduleEmpToResult(resultSchedulEmpVieweModel, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func ScheduleGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultScheduleModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QuerySchedule
	GetOrGetIdProses(w, r, key, 0, queryForProces)
}

func ScheduleEmployeeView(w http.ResponseWriter, r *http.Request) {
	var key string
	resultSchedulEmpVieweModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryScheduleView
	GetOrGetIdProses(w, r, key, 1, queryForProces)
}

func EmployeeScheduleGetWithId(w http.ResponseWriter, r *http.Request) {
	var key string
	resultScheduleModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryEmployeeFromId + key + "'"
	GetOrGetIdProses(w, r, "", 0, queryForProces)
}

func SchedulePost(w http.ResponseWriter, r *http.Request) {
	resultScheduleModel = nil
	function.Cors(&w, r)
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&scheduleModel); err != nil {
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

		// QUERY PROSES
		queryPost := query.QueryScheduleInsert + "VALUES ( '" + uuid.String() + "' , '" + scheduleModel.Day + "' , '" + scheduleModel.StartAt.String() + "', '" + scheduleModel.EndAt.String() + "', '" + scheduleModel.Description + "') RETURNING * ;"

		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		// fmt.Println(resultScheduleModel)
		message := "The data success to save !!"
		result, err := dataToResult(resultScheduleModel, message)
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

func ScheduleEdit(w http.ResponseWriter, r *http.Request) {
	resultScheduleModel = nil
	function.Cors(&w, r)
	if r.Method == "PUT" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&scheduleModel); err != nil {
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
		queryUpd := "SET description= '" + scheduleModel.Description + "', day='" + scheduleModel.Day + "', start_at = '" + scheduleModel.StartAt.String() + "', end_at = '" + scheduleModel.EndAt.String() + "' WHERE id = '" + key + "' returning *;"

		queryPost := query.QueryScheduleUpdate + queryUpd
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		message := "The data success to save !!"
		result, err := dataToResult(resultScheduleModel, message)
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

func ScheduleDelete(w http.ResponseWriter, r *http.Request) {
	resultScheduleModel = nil
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

		queryGet := query.QueryScheduleDelete + key + "' Returning *"
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept DELETE request", http.StatusBadRequest)
	}
}
