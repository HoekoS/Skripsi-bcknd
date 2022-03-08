package menu_makanan

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
	"skripsi.com/backend-adm/uuid_generator"
)

func ConDb() {
	var err error

	dbPool, err := pgxpool.Connect(context.Background(), os.Getenv("GO_DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Berhasil")
	}

	// to close DB pool
	defer dbPool.Close()
	rows, err := dbPool.Query(context.Background(), "select * from test")
	if err != nil {
		log.Println("error while executing query")
	}

	// iterate through the rows
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Println("error while iterating dataset")
		}

		// convert DB types to Go types
		id := values[0].(int32)
		desc := values[1].(string)

		log.Println("[id:", id, ", desc:", desc, "]")
	}
}

var menuModel *model_struct.Menu = new(model_struct.Menu)

var resultModel []model_struct.Menu

type resultStatus struct {
	Succes     string              `json:"status"`
	StatusCode string              `json:"status code"`
	Message    string              `json:"message"`
	Data       []model_struct.Menu `json:"data"`
}

func queryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.Menu, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&menuModel.Id, &menuModel.Name, &menuModel.Description, &menuModel.Price, &menuModel.Pic, &menuModel.CreateAt, &menuModel.Kategori, &menuModel.SubKategori)

		if err != nil {
			log.Println("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultModel = append(resultModel, *menuModel)
	}

	if err = rows.Err(); err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultModel, err
}

func dataToResult(resultModel []model_struct.Menu, message string) ([]byte, error) {
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

func GetOrGetIdProses(w http.ResponseWriter, r *http.Request, key string, ctg string, queryForProces string) {
	var queryGet string
	if r.Method == "GET" {
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		if ctg != "" {
			log.Print("With Category")
			queryGet = queryForProces + " WHERE sub_kategori = '" + ctg + "'"
		} else if key != "" {
			log.Print("With Key")
			queryGet = queryForProces + " WHERE ID = '" + key + "'"
		} else {
			log.Print("No Key")
			queryGet = queryForProces
		}
		queryProces(dbPool, queryGet)
		// fmt.Println(queryGet)
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
	}
}

func ViewMenuMakanan(w http.ResponseWriter, r *http.Request) {
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
	queryForProces := query.QueryMenuMakanan
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func ViewMenuMinuman(w http.ResponseWriter, r *http.Request) {
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
	queryForProces := query.QueryMenuMinuman
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func ViewKategori(w http.ResponseWriter, r *http.Request) {
	var ctg string
	resultModel = nil
	function.Cors(&w, r)
	ctgy, ok := r.URL.Query()["ctg"]
	if !ok || len(ctgy[0]) < 1 {
		ctg = ""
		log.Println("Url Param 'ctg' is missing")
	} else {
		ctg = ctgy[0]
	}
	queryForProces := query.QueryMenu
	GetOrGetIdProses(w, r, "", ctg, queryForProces)
}

func MenuAll(w http.ResponseWriter, r *http.Request) {
	var key string
	resultModel = nil
	function.Cors(&w, r)

	// tokenCookie, err := r.Cookie("token")
	// if err != nil {
	// 	log.Fatalf("Error occured while reading cookie")
	// }
	// fmt.Println("\nPrinting cookie with name as token")
	// fmt.Println(tokenCookie)

	// fmt.Println("\nPrinting all cookies")
	// for _, c := range r.Cookies() {
	// 	fmt.Println(c)
	// }

	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryMenu
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func MenuDetail(w http.ResponseWriter, r *http.Request) {
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
	queryForProces := query.QueryMenu
	GetOrGetIdProses(w, r, key, "", queryForProces)
}

func MenuPost(w http.ResponseWriter, r *http.Request) {
	resultModel = nil
	function.Cors(&w, r)
	if r.Method == "POST" {
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")
		description := r.FormValue("description")
		price := r.FormValue("price")
		kategori := r.FormValue("kategori")
		subKategori := r.FormValue("subkategori")
		if name == "" || description == "" || price == "" {
			log.Println("name description not null", err)
			http.Error(w, "name description null", http.StatusInternalServerError)

		}
		uploadedFile, handler, err := r.FormFile("file")
		if err != nil {
			log.Println("error di uplodedfile", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()

		namefile := handler.Filename
		path := "file-upload/menu/"
		filename, err := function.UploadImage(&uploadedFile, namefile, path, w)
		if err != nil {
			log.Println("error upload :", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// fmt.Println(filename, name, description)

		uuid := uuid_generator.UuidGenerate()
		var now = time.Now()
		formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		queryPost := query.QueryMenuInsert + "VALUES ( '" + uuid.String() + "' , '" + name + "' , '" + description + "' , " + price + " , '" + filename + "' , '" + formatted + "' , '" + kategori + "' , '" + subKategori + "') RETURNING * ;"
		fmt.Println(queryPost)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Println("error while executing query")
		}

		message := "The data success to save !!"
		result, err := dataToResult(resultModel, message)
		// fmt.Print(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

func MenuEdit(w http.ResponseWriter, r *http.Request) {
	resultModel = nil
	queryFile := ""
	function.Cors(&w, r)
	if r.Method == "PUT" {

		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")
		description := r.FormValue("description")
		price := r.FormValue("price")
		kategori := r.FormValue("kategori")
		subkategori := r.FormValue("subkategori")
		if name == "" || description == "" || price == "" {
			log.Println("name description not null", err)
			http.Error(w, "name description null", http.StatusInternalServerError)

		}
		uploadedFile, handler, err := r.FormFile("file")
		// fmt.Println(uploadedFile, handler)
		if err == nil {
			namefile := handler.Filename
			path := "file-upload/menu/"
			filename, err := function.UploadImage(&uploadedFile, namefile, path, w)
			if err != nil {
				log.Println("error upload :", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Println(filename, name, description)
			queryFile = ", pic='" + filename + "'"
			uploadedFile.Close()
		} else {
			log.Println("error upload file", err)
		}
		fmt.Println(queryFile)

		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		key := keys[0]
		queryUpd := "SET name= '" + name + "', description= '" + description + "', price=" + price + queryFile + ", kategori = '" + kategori + "', sub_kategori = '" + subkategori + "' WHERE id = '" + key + "' returning *;"

		queryGet := query.QueryMenuUpdate + queryUpd
		fmt.Println(queryGet)
		queryProces(dbPool, queryGet)
		if err != nil {
			log.Println("error while executing query", err)
		}

		message := "You'r data succes to send !"
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

func MenuDelete(w http.ResponseWriter, r *http.Request) {
	resultModel = nil
	// function.Cors(&w, r)
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

		queryGet := query.QueryMenuDelete + key + "' Returning *"
		queryProces(dbPool, queryGet)
		if err != nil {
			log.Println("error while executing query")
		}
		// fmt.Println(queryGet)

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

func MenuDashboard(w http.ResponseWriter, r *http.Request) {
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
	queryForProces := query.QueryMenuDash
	GetOrGetIdProses(w, r, key, "", queryForProces)
}
