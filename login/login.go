package login

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"skripsi.com/backend-adm/databases"
	"skripsi.com/backend-adm/function"
	"skripsi.com/backend-adm/model_struct"
	"skripsi.com/backend-adm/query"
)

var err error

var userLoginViewModel *model_struct.UserLogin = new(model_struct.UserLogin)

var resultUserlLoginVieweModel []model_struct.UserLogin

type resulUserLogintStatus struct {
	Succes     string                   `json:"status"`
	StatusCode string                   `json:"status code"`
	Message    string                   `json:"message"`
	Data       []model_struct.UserLogin `json:"data"`
}

type failedStatus struct {
	Succes     string `json:"status"`
	StatusCode string `json:"status code"`
	Message    string `json:"message"`
	Data       string `json:"data"`
}

func queryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.UserLogin, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&userLoginViewModel.Id, &userLoginViewModel.Username, &userLoginViewModel.Password, &userLoginViewModel.Email, &userLoginViewModel.RoleCode, &userLoginViewModel.CreatedOn)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultUserlLoginVieweModel = append(resultUserlLoginVieweModel, *userLoginViewModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultUserlLoginVieweModel, err
}

func dataToResult(resultModel []model_struct.UserLogin, message string) ([]byte, error) {
	// log.Println(resultModel)
	if resultModel != nil {
		// fmt.Println("true")
		jsondat := &resulUserLogintStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultUserlLoginVieweModel,
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
			Message:    "No User Found",
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

func getHash(pwd []byte) string {
	hashPasw, err := bcrypt.GenerateFromPassword([]byte(userLoginViewModel.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hashPasw)
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
			queryGet = queryForProces + " ORDER BY id ASC"
		}
		queryProces(dbPool, queryGet)
		fmt.Println(queryGet)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResult(resultUserlLoginVieweModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func UserLoginGet(w http.ResponseWriter, r *http.Request) {
	var key string
	resultUserlLoginVieweModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryUserLogin
	GetOrGetIdProses(w, r, key, queryForProces)
}

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	resultUserlLoginVieweModel = nil
	function.Cors(&w, r)
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&userLoginViewModel); err != nil {
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

		hashPasw := getHash([]byte(userLoginViewModel.Password))

		formatted := function.FormatTime()

		// QUERY PROSES
		queryPost := query.QueryUserLoginInsert + "VALUES ( '" + userLoginViewModel.Username + "' , '" + hashPasw + "', '" + userLoginViewModel.Email + "', '" + userLoginViewModel.RoleCode + "', '" + formatted + "') RETURNING * ;"

		fmt.Println(queryPost)
		// fmt.Println(err)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		// fmt.Println(resultUserlLoginVieweModel)
		message := "The data success to save !!"
		result, err := dataToResult(resultUserlLoginVieweModel, message)
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

func UserSearch(email string) (succes string, err error) {
	dbPool, err := databases.Connect()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	// QUERY PROSES
	queryPost := query.QueryUserSearch + "'" + email + "'"
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

func UserCekLogin(w http.ResponseWriter, r *http.Request) {
	resultUserlLoginVieweModel = nil
	function.Cors(&w, r)
	var result []byte
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&userLoginViewModel); err != nil {
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

		userPass := userLoginViewModel.Password
		getUser, err := UserSearch(userLoginViewModel.Email)
		if err != nil {
			log.Println("Error Search User")
		}
		fmt.Println(getUser)
		message := "Succes Login to aplication !!"
		result, err = dataToResult(resultUserlLoginVieweModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		cekPassErr := bcrypt.CompareHashAndPassword([]byte(resultUserlLoginVieweModel[0].Password), []byte(userPass))
		if cekPassErr != nil {
			http.Error(w, cekPassErr.Error(), http.StatusInternalServerError)
			result = []byte("null")
		}
		jwtToken, err := function.GenerateJWT()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			result = []byte("null")
		}
		cookieName := "Token"
		c := &http.Cookie{}
		if storedCookie, _ := r.Cookie(cookieName); storedCookie != nil {
			c = storedCookie
		}
		if c.Value == "" {
			c = &http.Cookie{}
			c.Name = cookieName
			c.Value = jwtToken
			c.Expires = time.Now().Add(5 * time.Minute)
			http.SetCookie(w, c)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept POST request", http.StatusBadRequest)
	}
}

// func UserLoginEdit(w http.ResponseWriter, r *http.Request) {
// 	resultUserlLoginVieweModel = nil
// 	function.Cors(&w, r)
// 	if r.Method == "PUT" {
// 		// DECODE BODY JSON FROM HTTP
// 		decoder := json.NewDecoder(r.Body)
// 		if err = decoder.Decode(&userLoginViewModel); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		// Connect to DB
// 		dbPool, err := databases.Connect()
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			os.Exit(1)
// 		}
// 		defer dbPool.Close()

// 		keys, ok := r.URL.Query()["key"]
// 		if !ok || len(keys[0]) < 1 {
// 			log.Println("Url Param 'key' is missing")
// 			return
// 		}
// 		key := keys[0]

// 		// QUERY PROSES
// 		queryUpd := "SET description= '" + userLoginViewModel.Description + "', day='" + userLoginViewModel.Day + "', start_at = '" + userLoginViewModel.StartAt.String() + "', end_at = '" + userLoginViewModel.EndAt.String() + "' WHERE id = '" + key + "' returning *;"

// 		queryPost := query.QueryUserLoginUpdate + queryUpd
// 		fmt.Println(queryPost)
// 		// fmt.Println(err)
// 		queryProces(dbPool, queryPost)
// 		if err != nil {
// 			log.Fatal("error while executing query")
// 		}

// 		// RESPON RESULT STATUS
// 		message := "The data success to save !!"
// 		result, err := dataToResult(resultUserlLoginVieweModel, message)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(result)
// 	} else {
// 		http.Error(w, "Only accept PUT request", http.StatusBadRequest)
// 	}
// }

func UserLoginDelete(w http.ResponseWriter, r *http.Request) {
	resultUserlLoginVieweModel = nil
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

		queryGet := query.QueryUserLoginDelete + key + "' Returning *"
		queryProces(dbPool, queryGet)
		if err != nil {
			log.Println("error while executing query")
		}
		fmt.Println(queryGet)

		jsondat := &resulUserLogintStatus{
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

func UserSignOut(w http.ResponseWriter, r *http.Request) {
	function.Cors(&w, r)
	if r.Method == "GET" {
		c := &http.Cookie{}
		c.Name = "Token"
		c.Expires = time.Unix(0, 0)
		c.MaxAge = -1
		http.SetCookie(w, c)

		jsondat := &resulUserLogintStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    "Success Logout !",
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}
