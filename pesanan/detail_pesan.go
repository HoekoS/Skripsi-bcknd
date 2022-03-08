package pesanan

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

var pesananModel *model_struct.DetailPesanan = new(model_struct.DetailPesanan)

var detailPesananResultModel []model_struct.DetailPesanan

var viewPesananModel *model_struct.ViewPesan = new(model_struct.ViewPesan)

var detailviewPesananModel []model_struct.ViewPesan

var totalHargaPesananModel *model_struct.TotalHargaPesanan = new(model_struct.TotalHargaPesanan)

var detailTotalHargaPesananModel []model_struct.TotalHargaPesanan

type resultStatus struct {
	Succes     string                       `json:"status"`
	StatusCode string                       `json:"status code"`
	Message    string                       `json:"message"`
	Data       []model_struct.DetailPesanan `json:"data"`
}

type viewResultStatus struct {
	Succes     string                   `json:"status"`
	StatusCode string                   `json:"status code"`
	Message    string                   `json:"message"`
	Data       []model_struct.ViewPesan `json:"data"`
}

type totalHargaResultStatus struct {
	Succes     string                           `json:"status"`
	StatusCode string                           `json:"status code"`
	Message    string                           `json:"message"`
	Data       []model_struct.TotalHargaPesanan `json:"data"`
}

func queryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.DetailPesanan, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&pesananModel.Id, &pesananModel.NoNota, &pesananModel.IdMenu, &pesananModel.Quantity, &pesananModel.Catatan)

		if err != nil {
			log.Fatal("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		detailPesananResultModel = append(detailPesananResultModel, *pesananModel)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return detailPesananResultModel, err
}

func dataToResult(resultModel []model_struct.DetailPesanan, message string) ([]byte, error) {
	if resultModel != nil {
		jsondat := &resultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &resultStatus{
			Succes:     "False",
			StatusCode: "404",
			Message:    "No Data",
			Data:       resultModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func queryProcesView(dbPool *pgxpool.Pool, query string) ([]model_struct.ViewPesan, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&viewPesananModel.Id, &viewPesananModel.NomerNota, &viewPesananModel.NamaPelanggan, &viewPesananModel.Quantity, &viewPesananModel.Catatan, &viewPesananModel.Price, &viewPesananModel.TotalHarga, &viewPesananModel.Url)

		if err != nil {
			log.Fatal("error di scan proces view kamu ini detailnya : ", err.Error())
			return nil, err
		}

		detailviewPesananModel = append(detailviewPesananModel, *viewPesananModel)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return detailviewPesananModel, err
}

func dataToResultView(viewResultModel []model_struct.ViewPesan, message string) ([]byte, error) {
	if viewResultModel != nil {
		jsondat := &viewResultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       viewResultModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &viewResultStatus{
			Succes:     "False",
			StatusCode: "404",
			Message:    "No Data",
			Data:       viewResultModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func totalHargaProcessQUery(dbPool *pgxpool.Pool, query string) ([]model_struct.TotalHargaPesanan, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&totalHargaPesananModel.NomerNota, &totalHargaPesananModel.TotalHarga)

		if err != nil {
			log.Fatal("error di scan proces view kamu ini detailnya : ", err.Error())
			return nil, err
		}

		detailTotalHargaPesananModel = append(detailTotalHargaPesananModel, *totalHargaPesananModel)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return detailTotalHargaPesananModel, err
}

func totalHargaResult(totalHargaPesananModel []model_struct.TotalHargaPesanan, message string) ([]byte, error) {
	if totalHargaPesananModel != nil {
		jsondat := &totalHargaResultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       totalHargaPesananModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &totalHargaResultStatus{
			Succes:     "False",
			StatusCode: "404",
			Message:    "No Data",
			Data:       totalHargaPesananModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func DetailPesananGet(w http.ResponseWriter, r *http.Request) {
	var key string
	var not string
	var state string
	var queryForProces string
	detailPesananResultModel = nil
	detailviewPesananModel = nil
	function.Cors(&w, r)

	nota, err := r.URL.Query()["not"]
	if !err || len(nota[0]) < 1 {
		not = ""
		log.Println("Url Param 'key' is missing")
	} else {
		fmt.Println("sini")
		state = "pilih"
		not = nota[0]
	}

	keys, ok2 := r.URL.Query()["key"]
	if !ok2 || len(keys[0]) < 1 {
		key = ""
		state = "all"
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
		state = "all"
	}

	if nota != nil {
		state = "pilih"
		queryForProces = query.QueryDetailPesananNota
	} else {
		queryForProces = query.QueryDetailPesanan
	}
	fmt.Println(state)
	GetOrGetIdProses(w, r, key, not, queryForProces, state)
}

func DetailPesananGetOne(w http.ResponseWriter, r *http.Request) {
	var id, not string
	detailviewPesananModel = nil
	function.Cors(&w, r)

	nota, err := r.URL.Query()["not"]
	if !err || len(nota[0]) < 1 {
		not = ""
		log.Println("Url Param 'key' is missing")
	} else {
		fmt.Println("sini")
		not = nota[0]
	}

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		id = ""
		log.Println("Url Param 'id' is missing")
	} else {
		id = ids[0]
	}
	queryForProces := query.QueryDetailPesananNota
	GetOrGetIdProses(w, r, id, not, queryForProces, "detail")
}

func TotalHargaGet(w http.ResponseWriter, r *http.Request) {
	var nota string
	detailTotalHargaPesananModel = nil
	function.Cors(&w, r)
	not, ok := r.URL.Query()["not"]
	if !ok || len(not[0]) < 1 {
		nota = ""
		log.Println("Url Param 'ctg' is missing")
	} else {
		nota = not[0]
	}
	queryForProces := query.QueryTotalHarga
	GetOrGetIdProses(w, r, "", nota, queryForProces, "total")
}

func DetailPesananData(w http.ResponseWriter, r *http.Request) {
	detailPesananResultModel = nil
	var err error
	function.Cors(&w, r)
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&pesananModel); err != nil {
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
		queryPost := query.QueryDetailPesananInsert + "VALUES ( " + "'" + pesananModel.NoNota + "' , '" + pesananModel.IdMenu + "' , " + strconv.Itoa(pesananModel.Quantity) + ", '" + pesananModel.Catatan + "') RETURNING * ;"
		fmt.Println("detail " + queryPost)
		queryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		message := "The data success to save !!"
		result, err := dataToResult(detailPesananResultModel, message)
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

func DetailPesananUpdate(w http.ResponseWriter, r *http.Request) {
	detailPesananResultModel = nil
	var err error
	if r.Method == "PUT" {

		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&pesananModel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
		queryUp := "SET quantity= " + strconv.Itoa(pesananModel.Quantity) + ",catatan='" + pesananModel.Catatan + "'WHERE id = " + key + " returning *;"

		queryGet := query.QueryDetailPesananUpdate + queryUp
		fmt.Println(queryGet)
		queryProces(dbPool, queryGet)
		if err != nil {
			log.Fatal("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := dataToResult(detailPesananResultModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept PUT request", http.StatusBadRequest)
	}
}

func DetailPesananDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		// CONECT TO DB
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		// Get Id
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		key := keys[0]

		// DELETE QUERY
		queryGet := query.QueryDetailPesananDelete + key + " Returning *"
		queryProces(dbPool, queryGet)

		// QUERY PROSES
		if err != nil {
			log.Fatal("error while executing query")
		}
		fmt.Println(queryGet)

		// RESPON RESULT STATUS
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

func DetailPesananGetAll(w http.ResponseWriter, r *http.Request) {
	var id, not string
	detailviewPesananModel = nil
	function.Cors(&w, r)
	queryForProces := query.QueryDetailPesananNota
	GetOrGetIdProses(w, r, id, not, queryForProces, "all-detail")
}
