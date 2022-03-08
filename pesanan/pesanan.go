package pesanan

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	qrcode "github.com/skip2/go-qrcode"
	"skripsi.com/backend-adm/databases"
	"skripsi.com/backend-adm/function"
	"skripsi.com/backend-adm/model_struct"
	"skripsi.com/backend-adm/query"
	"skripsi.com/backend-adm/uuid_generator"
)

var pesanModel *model_struct.Pesanan = new(model_struct.Pesanan)

var pesanResultModel []model_struct.Pesanan

var err error

type pesanResultStatus struct {
	Succes     string                 `json:"status"`
	StatusCode string                 `json:"status code"`
	Message    string                 `json:"message"`
	Data       []model_struct.Pesanan `json:"data"`
}

func pesanQueryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.Pesanan, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(
			&pesanModel.Id,
			&pesanModel.FlagTakeAway,
			&pesanModel.State,
			&pesanModel.Date,
			&pesanModel.NomerNota,
			&pesanModel.NamaPelanggan,
			&pesanModel.NomerMeja,
			&pesanModel.Url,
			&pesanModel.DeleteStatus)

		if err != nil {
			log.Fatal("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		pesanResultModel = append(pesanResultModel, *pesanModel)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return pesanResultModel, err
}

func pesanDataToResult(pesanResultModel []model_struct.Pesanan, message string) ([]byte, error) {
	if pesanResultModel != nil {
		jsondat := &pesanResultStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       pesanResultModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &pesanResultStatus{
			Succes:     "False",
			StatusCode: "404",
			Message:    "No Data",
			Data:       pesanResultModel,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func GetOrGetIdProses(w http.ResponseWriter, r *http.Request, key string, not string, queryForProces string, state string) {
	var queryGet string
	var result []byte
	if r.Method == "GET" {
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		if not != "" {
			queryGet = queryForProces + " WHERE no_nota = '" + not + "'"
		} else if key != "" {
			log.Print("With Key")
			queryGet = queryForProces + " AND ID = '" + key + "'"
		} else {
			log.Print("No Key")
			queryGet = queryForProces + " order by no_nota asc"
		}
		// fmt.Println(queryGet)
		message := "You'r data succes to send !"
		if state == "pesan" {
			// fmt.Println("sini")
			pesanQueryProces(dbPool, queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			result, err = pesanDataToResult(pesanResultModel, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if state == "all" {
			// fmt.Println("sini2")
			queryProces(dbPool, queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			result, err = dataToResult(detailPesananResultModel, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if state == "pilih" {
			queryGet += " Order BY id"
			queryProcesView(dbPool, queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			result, err = dataToResultView(detailviewPesananModel, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if state == "total" {
			queryGet += " GROUP BY no_nota"
			totalHargaProcessQUery(dbPool, queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			result, err = totalHargaResult(detailTotalHargaPesananModel, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if state == "detail" {
			queryGet = queryForProces + " WHERE dp.id = " + key
			queryProcesView(dbPool, queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			result, err = dataToResultView(detailviewPesananModel, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if state == "all-detail" {
			queryGet = queryForProces
			queryProcesView(dbPool, queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			result, err = dataToResultView(detailviewPesananModel, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		fmt.Println(queryGet)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func PesananGetAll(w http.ResponseWriter, r *http.Request) {
	var key string
	pesanResultModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryPesananNoDone
	GetOrGetIdProses(w, r, key, "", queryForProces, "pesan")
}

func PesananGetOneByNot(w http.ResponseWriter, r *http.Request) {
	var key string
	pesanResultModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryPesananGetOneNot + "'" + key + "'"
	GetOrGetIdProses(w, r, "", "", queryForProces, "pesan")
}

func Pesanan(w http.ResponseWriter, r *http.Request) {
	pesanResultModel = nil
	function.Cors(&w, r)
	if r.Method == "POST" {
		// DECODE BODY JSON FROM HTTP
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&pesanModel); err != nil {
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
		//path for QR
		path := "file-upload/qr/" + uuid.String() + ".png"
		// QUERY PROSES
		queryPost := query.QueryPesananInsert + "VALUES ( '" + uuid.String() + "' , " + strconv.Itoa(pesanModel.FlagTakeAway) + " , '" + pesanModel.State + "' , '" + formatted + "', '" + pesanModel.NamaPelanggan + "', " + strconv.Itoa(pesanModel.NomerMeja) + ",'" + path + "') RETURNING * ;"
		fmt.Println(queryPost)
		err = qrcode.WriteFile("http://localhost:3636/menu_all.html?nom="+strconv.Itoa(pesanModel.NomerMeja)+"&not="+pesanModel.NomerNota, qrcode.Medium, 256, path)
		// fmt.Println(err)
		pesanQueryProces(dbPool, queryPost)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		message := "The data success to save !!"
		result, err := pesanDataToResult(pesanResultModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept GET or POST request", http.StatusBadRequest)
	}
}

// func PesananUpdate(w http.ResponseWriter, r *http.Request) {
// }

func PesananDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		// CONECT TO DB
		dbPool, err := databases.Connect()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer dbPool.Close()

		// Get Id
		key, err := function.GetId(r)
		if err != nil {
			log.Fatal("error while get Id")
		}

		// DELETE QUERY
		queryDelete := query.QueryPesananDelete + "'" + key + "'" + " Returning *"
		fmt.Println(queryDelete)

		// QUERY PROSES
		pesanQueryProces(dbPool, queryDelete)
		if err != nil {
			log.Fatal("error while executing query")
		}

		// RESPON RESULT STATUS
		jsondat := &pesanResultStatus{
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
		http.Error(w, "Only accept GET,PUT, and DELETE request", http.StatusBadRequest)
	}
}

func DetailPesananUpdateState(w http.ResponseWriter, r *http.Request) {
	detailPesananResultModel = nil
	var err error
	if r.Method == "PUT" {

		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&pesanModel); err != nil {
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
		queryUp := "SET state= '" + pesanModel.State + "' WHERE id = '" + key + "' returning *;"

		queryGet := query.QueryPesananUpdate + queryUp
		fmt.Println(queryGet)
		pesanQueryProces(dbPool, queryGet)
		if err != nil {
			log.Fatal("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := pesanDataToResult(pesanResultModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if pesanModel.State == "Done" || pesanModel.State == "BATAL PESAN" {
			err := HistoryPost(pesanModel.NomerNota, pesanModel.State)
			if err != nil {
				fmt.Println(err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept PUT request", http.StatusBadRequest)
	}
}

func DetailPesananUpdateStateNot(w http.ResponseWriter, r *http.Request) {
	detailPesananResultModel = nil
	var err error
	if r.Method == "PUT" {

		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&pesanModel); err != nil {
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
		queryUp := "SET state= '" + pesanModel.State + "' WHERE no_nota = '" + key + "' returning *;"

		queryGet := query.QueryPesananUpdate + queryUp
		fmt.Println(queryGet)
		pesanQueryProces(dbPool, queryGet)
		if err != nil {
			log.Fatal("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := pesanDataToResult(pesanResultModel, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if pesanModel.State == "Done" || pesanModel.State == "BATAL PESAN" {
			err := HistoryPost(pesanModel.NomerNota, pesanModel.State)
			if err != nil {
				fmt.Println(err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else {
		http.Error(w, "Only accept PUT request", http.StatusBadRequest)
	}
}

func PesananGetForDapur(w http.ResponseWriter, r *http.Request) {
	var key string
	pesanResultModel = nil
	function.Cors(&w, r)
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		key = ""
		log.Println("Url Param 'key' is missing")
	} else {
		key = keys[0]
	}
	queryForProces := query.QueryPesananDapur
	GetOrGetIdProses(w, r, key, "", queryForProces, "pesan")
}

func QrCetakFix(w http.ResponseWriter, r *http.Request) {
	function.Cors(&w, r)
	var key string
	if r.Method == "GET" {

		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			key = ""
			log.Println("Url Param 'date' is missing")
		} else {
			key = keys[0]
		}
		// key = "test"

		m := pdf.NewMarotoCustomSize(consts.Portrait, "C6", "mm", 80, 130)
		m.SetPageMargins(10, 10, 10)
		datas := getDataQrExport(r, key)
		fmt.Println("setelah datas ", datas[0])
		buildHeadingQr(m, datas)
		// fmt.Println("setelah buildlist")

		path := "/file-upload/pdf/qr/" + key + ".pdf"
		err := m.OutputFileAndClose("." + path)
		if err != nil {
			fmt.Println("⚠️  Could not save PDF:", err)
			os.Exit(1)
		}

		fmt.Println("PDF saved successfully")

		f, err := os.Open("." + path)
		if f != nil {
			defer f.Close()
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// contentDisposition := fmt.Sprintf("attachment; filename= qr-" + datas[0][0] + ".pdf")
		// w.Header().Set("Content-Disposition", contentDisposition)

		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsondat := &failedStatus{
			Succes:     "True",
			StatusCode: "200",
			Message:    "You'r data succes to delete !",
			Data:       path,
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
	}
}

func getDataQrExport(r *http.Request, key string) [][]string {
	resultHistoryView = nil

	dbPool, err := databases.Connect()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	queryForProces := query.QueryPesananNoDone + " AND ID = '" + key + "'"
	pesanQueryProces(dbPool, queryForProces)
	if err != nil {
		log.Println("error while executing query")
	}
	datas := changeTypeQr(pesanResultModel)
	// fmt.Println(datas)
	return datas
}

func changeTypeQr(data []model_struct.Pesanan) [][]string {
	var datas [][]string
	var status string

	for _, v := range data {
		if v.FlagTakeAway == 1 {
			status = "Takeaway"
		} else {
			status = "Dine In"
		}
		datas = append(datas, []string{
			v.NomerNota,
			strconv.Itoa(v.NomerMeja),
			v.NamaPelanggan,
			status,
			v.Url,
		})
	}
	// fmt.Println(datas[0][1])
	return datas
}

func buildHeadingQr(m pdf.Maroto, datas [][]string) {
	// m.RegisterHeader(func() {
	m.Row(10, func() {
		m.Col(12, func() {
			_ = m.FileImage("./file-upload/SAKO_JADI.png", props.Rect{
				Percent: 100,
				Center:  true,
			})
		})
	})
	m.Row(4, func() {
		m.Col(12, func() {
			m.Text("CAFE SAKO", props.Text{
				Top:   1,
				Size:  12,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.Row(8, func() {
		m.Col(12, func() {
			m.Text("Jl. Cemara Raya No.836", props.Text{
				Top:   2,
				Size:  9,
				Align: consts.Center,
			})
		})
	})

	m.Row(4, func() {
		m.Col(12, func() {
			m.Text("23 February 2022, 22:03:09", props.Text{
				Top:   1,
				Size:  6,
				Align: consts.Left,
			})
		})
	})
	m.Row(4, func() {
		m.Col(6, func() {
			m.Text("Nomor Nota", props.Text{
				Size:  8,
				Align: consts.Left,
			})
		})
		m.Col(6, func() {
			m.Text(datas[0][0], props.Text{
				Size:  8,
				Align: consts.Right,
			})
		})
	})
	m.Row(4, func() {
		m.Col(6, func() {
			m.Text("Nomor Meja", props.Text{
				Size:  8,
				Align: consts.Left,
			})
		})
		m.Col(6, func() {
			m.Text(datas[0][1], props.Text{
				Size:  8,
				Align: consts.Right,
			})
		})
	})
	m.Row(4, func() {
		m.Col(6, func() {
			m.Text("Nama Pelanggan", props.Text{
				Size:  8,
				Align: consts.Left,
			})
		})
		m.Col(6, func() {
			m.Text(datas[0][2], props.Text{
				Size:  8,
				Align: consts.Right,
			})
		})
	})
	m.Line(2)
	m.Row(4, func() {
		m.Col(12, func() {
			m.Text(datas[0][3], props.Text{
				Size:  9,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
		// m.Line(5)
	})
	m.Line(2)

	m.Row(50, func() {
		m.Col(12, func() {
			_ = m.FileImage("./"+datas[0][4], props.Rect{
				Center:  true,
				Percent: 100,
			})
		})
	})

	m.Line(2)
	// })
}
