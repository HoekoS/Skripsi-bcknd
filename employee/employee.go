package employee

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

var employeeModel *model_struct.Employee = new(model_struct.Employee)

var resultEmployeeModel []model_struct.Employee

type resultStatus struct {
	Succes     string                  `json:"status"`
	StatusCode string                  `json:"status code"`
	Message    string                  `json:"message"`
	Data       []model_struct.Employee `json:"data"`
}

type failedStatus struct {
	Succes     string `json:"status"`
	StatusCode string `json:"status code"`
	Message    string `json:"message"`
	Data       string `json:"data"`
}

func queryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.Employee, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(
			&employeeModel.Id,
			&employeeModel.FullName,
			&employeeModel.Gender,
			&employeeModel.BirthDate,
			&employeeModel.BirthPlace,
			&employeeModel.Address,
			&employeeModel.Email,
			&employeeModel.PhoneNumber,
			&employeeModel.StaffStatus,
			&employeeModel.StaffPosition,
			&employeeModel.StaffActiveStatus,
			&employeeModel.StaffDateJoin,
		)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultEmployeeModel = append(resultEmployeeModel, *employeeModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultEmployeeModel, err
}

func dataToResult(resultModel []model_struct.Employee, message string) ([]byte, error) {
	// log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultEmployeeModel,
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
			queryGet = queryForProces + " WHERE ID = '" + key + "'"
		} else {
			log.Print("No Key")
			queryGet = queryForProces
		}
		queryProces(dbPool, queryGet)
		fmt.Println(queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResult(resultEmployeeModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func EmployeeGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultEmployeeModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryEmployee
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func EmployeePost(w http.ResponseWriter, r *http.Request) {
	resultEmployeeModel = nil
	function.Cors(&w, r)
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&employeeModel); err != nil {
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
		queryPost := query.QueryEmployeeInsert + "VALUES ( '" + uuid.String() + "' , '" + employeeModel.FullName + "' , '" + employeeModel.Gender + "' , '" + employeeModel.BirthDate + "' , '" + employeeModel.BirthPlace + "' , '" + employeeModel.Address + "' , '" + employeeModel.Email + "' , '" + employeeModel.PhoneNumber + "' , '" + employeeModel.StaffStatus + "' , '" + employeeModel.StaffPosition + "' , '" + employeeModel.StaffActiveStatus + "' , '" + employeeModel.StaffDateJoin + `') RETURNING 
		id, full_name, gender, birth_date::Varchar, birth_place, address, email, phone_number, staff_status, staff_position, staff_active_status, staff_date_join::Varchar ;`

		fmt.Println(employeeModel.BirthDate)
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		// fmt.Println(resultEmployeeModel)
		message := "The data success to save !!"
		result, err := dataToResult(resultEmployeeModel, message)
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

func EmployeeEdit(w http.ResponseWriter, r *http.Request) {
	resultEmployeeModel = nil
	function.Cors(&w, r)
	if r.Method == "PUT" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&employeeModel); err != nil {
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
		queryUpd := "SET full_name= '" + employeeModel.FullName + "', gender= '" + employeeModel.Gender + "', birth_date = '" + employeeModel.BirthDate + "', birth_place = '" + employeeModel.BirthPlace + "', address = '" + employeeModel.Address + "', email = '" + employeeModel.Email + "', phone_number = '" + employeeModel.PhoneNumber + "', staff_status = '" + employeeModel.StaffStatus + "', staff_position = '" + employeeModel.StaffPosition + "', staff_active_status = '" + employeeModel.StaffActiveStatus + "', staff_date_join = '" + employeeModel.StaffDateJoin + "' WHERE id = '" + key + `' returning
		id, full_name, gender, birth_date::Varchar, birth_place, address, email, phone_number, staff_status, staff_position, staff_active_status, staff_date_join::Varchar ;`

		queryPost := query.QueryEmployeeUpdate + queryUpd
		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		message := "The data success to save !!"
		result, err := dataToResult(resultEmployeeModel, message)
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

func EmployeeDelete(w http.ResponseWriter, r *http.Request) {
	resultEmployeeModel = nil
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

		queryGet := query.QueryEmployeeDelete + key + "' Returning *"
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
