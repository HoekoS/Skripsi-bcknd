package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"skripsi.com/backend-adm/databases"
	"skripsi.com/backend-adm/function"
	"skripsi.com/backend-adm/model_struct"
	"skripsi.com/backend-adm/query"
)

var jadwalShift *model_struct.JadwalShiftDash = new(model_struct.JadwalShiftDash)

var resultJadwalShift []model_struct.JadwalShiftDash

type resultJadwalShiftStatus struct {
	Succes     string                         `json:"status"`
	StatusCode string                         `json:"status code"`
	Message    string                         `json:"message"`
	Data       []model_struct.JadwalShiftDash `json:"data"`
}

func queryProcesShift(dbPool *pgxpool.Pool, query string) ([]model_struct.JadwalShiftDash, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&jadwalShift.Nama, &jadwalShift.StaffStatus, &jadwalShift.StaffPosition, &jadwalShift.StartAt, &jadwalShift.EndAt, &jadwalShift.Description)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultJadwalShift = append(resultJadwalShift, *jadwalShift)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultJadwalShift, err
}

func dataToresultJadwalShift(resultModel []model_struct.JadwalShiftDash, message string) ([]byte, error) {
	log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resultJadwalShiftStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultJadwalShift,
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

func GetProsesShift(w http.ResponseWriter, r *http.Request, key string, ctg string, queryForProces string) {
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
		queryProcesShift(dbPool, queryGet)
		fmt.Println(queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToresultJadwalShift(resultJadwalShift, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func ShiftDashboardGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultJadwalShift = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	fmt.Print(key)
	day := changeDay()
	queryForProces := query.QueryEmployeeHari + "and s.day = '" + day + "'"
	GetProsesShift(w, r, key, "", queryForProces)
}

func changeDay() string {
	var day string
	day = time.Now().Weekday().String()
	if day == "Monday" {
		day = "Senin"
	} else if day == "Tuesday" {
		day = "Selasa"
	} else if day == "Wednesday" {
		day = "Rabu"
	} else if day == "Thursday" {
		day = "Kamis"
	} else if day == "Friday" {
		day = "Jumat"
	} else if day == "Saturday" {
		day = "Sabtu"
	} else if day == "Sunday" {
		day = "Minggu"
	}

	return day
}
