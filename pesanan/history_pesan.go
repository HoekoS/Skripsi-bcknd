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
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/jackc/pgx/v4/pgxpool"
	"skripsi.com/backend-adm/databases"
	"skripsi.com/backend-adm/function"
	"skripsi.com/backend-adm/model_struct"
	"skripsi.com/backend-adm/query"
	"skripsi.com/backend-adm/uuid_generator"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

var historyPesanModel *model_struct.HistoryPesan = new(model_struct.HistoryPesan)

var resultHistory []model_struct.HistoryPesan

var historyViewPesanModel *model_struct.ViewHistory = new(model_struct.ViewHistory)

var resultHistoryView []model_struct.ViewHistory

var historyDetailPesanModel *model_struct.ViewHistoryDetail = new(model_struct.ViewHistoryDetail)

var resultHistoryDetail []model_struct.ViewHistoryDetail

type resulStatusHistory struct {
	Succes     string                      `json:"status"`
	StatusCode string                      `json:"status code"`
	Message    string                      `json:"message"`
	Data       []model_struct.HistoryPesan `json:"data"`
}

type resulStatusHistoryView struct {
	Succes     string                     `json:"status"`
	StatusCode string                     `json:"status code"`
	Message    string                     `json:"message"`
	Data       []model_struct.ViewHistory `json:"data"`
}

type resultStatusHistoryDetail struct {
	Succes     string                           `json:"status"`
	StatusCode string                           `json:"status code"`
	Message    string                           `json:"message"`
	Data       []model_struct.ViewHistoryDetail `json:"data"`
}

type failedStatus struct {
	Succes     string `json:"status"`
	StatusCode string `json:"status code"`
	Message    string `json:"message"`
	Data       string `json:"data"`
}

func historyQueryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.HistoryPesan, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&historyPesanModel.Id, &historyPesanModel.NoNota, &historyPesanModel.TotalPrice, &historyPesanModel.Date)

		if err != nil {
			log.Fatal("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultHistory = append(resultHistory, *historyPesanModel)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultHistory, err
}

func historyDataToResult(resultHistory []model_struct.HistoryPesan, message string) ([]byte, error) {
	if resultHistory != nil {
		jsondat := &resulStatusHistory{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultHistory,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &resulStatusHistory{
			Succes:     "False",
			StatusCode: "404",
			Message:    "No Data",
			Data:       resultHistory,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func viewHistoryQueryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.ViewHistory, error) {
	fmt.Println("cobaaaaaaaaaaaaaaaaaaaaa")
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		log.Println("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(&historyViewPesanModel.NomerNota, &historyViewPesanModel.NamaPelanggan, &historyViewPesanModel.TotalPrice, &historyViewPesanModel.StatusTa, &historyViewPesanModel.Date)

		if err != nil {
			log.Fatal("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultHistoryView = append(resultHistoryView, *historyViewPesanModel)
	}
	fmt.Println("di view history", resultHistoryView, query)

	if err = rows.Err(); err != nil {
		log.Fatal("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	fmt.Println("di view history", resultHistoryView, query)
	return resultHistoryView, err
}

func detailHistoryDataToResult(resultHistory []model_struct.ViewHistory, message string) ([]byte, error) {
	if resultHistory != nil {
		jsondat := &resulStatusHistoryView{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultHistoryView,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &failedStatus{
			Succes:     "False",
			StatusCode: "404",
			Message:    "No Data",
			Data:       "",
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func historyDetailQueryProces(dbPool *pgxpool.Pool, query string) ([]model_struct.ViewHistoryDetail, error) {
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var err = rows.Scan(
			&historyDetailPesanModel.NomerNota,
			&historyDetailPesanModel.Name,
			&historyDetailPesanModel.Price,
			&historyDetailPesanModel.Quantity,
			&historyDetailPesanModel.TotalPrice,
			&historyDetailPesanModel.TotalSemuaPesanan,
			&historyDetailPesanModel.Catatan,
		)

		if err != nil {
			log.Fatal("error di scan kamu ini detailnya : ", err.Error())
			return nil, err
		}

		resultHistoryDetail = append(resultHistoryDetail, *historyDetailPesanModel)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("error di row kamu ini detailnya : ", err.Error())
		return nil, err
	}
	return resultHistoryDetail, err
}

func historyDetailDataToResult(resultHistory []model_struct.ViewHistoryDetail, message string) ([]byte, error) {
	if resultHistory != nil {
		jsondat := &resultStatusHistoryDetail{
			Succes:     "True",
			StatusCode: "200",
			Message:    message,
			Data:       resultHistoryDetail,
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	} else {
		jsondat := &failedStatus{
			Succes:     "False",
			StatusCode: "404",
			Message:    "No Data",
			Data:       "",
		}
		result, err := json.Marshal(jsondat)
		if err != nil {
			log.Fatal("error di json marshal ini detailnya : ", err.Error())
			return nil, err
		}
		return result, err
	}
}

func getHistoryProses(w http.ResponseWriter, r *http.Request, key string, not string, queryForProces string, state string) {
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
			queryGet = queryForProces + " WHERE hp.no_nota = '" + not + "'"
		} else if key != "" {
			log.Print("With Key")
			queryGet = queryForProces + " WHERE ID = '" + key + "'"
		} else {
			log.Print("No Key")
			queryGet = queryForProces
		}
		fmt.Println(queryGet)

		message := "You'r data succes to send !"
		if state == "nodt" {
			viewHistoryQueryProces(dbPool, queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			fmt.Println(resultHistoryView)
			result, err = detailHistoryDataToResult(resultHistoryView, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if state == "dtl" {
			historyDetailQueryProces(dbPool, queryGet)
			if err != nil {
				log.Println("error while executing query")
			}
			result, err = historyDetailDataToResult(resultHistoryDetail, message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

func HistoryPesanGet(w http.ResponseWriter, r *http.Request) {
	var date1, date2 string
	var where string
	resultHistoryView = nil
	function.Cors(&w, r)
	dates1, ok := r.URL.Query()["date1"]
	if !ok || len(dates1[0]) < 1 {
		date1 = ""
		log.Println("Url Param 'date' is missing")
	} else {
		date1 = dates1[0]
	}
	res1 := strings.Split(date1, "-")
	date1 = res1[2] + "-" + res1[1] + "-" + res1[0]
	fmt.Println(date1)
	// date1 =

	dates2, ok := r.URL.Query()["date2"]
	if !ok || len(dates2[0]) < 1 {
		date2 = ""
		log.Println("Url Param 'date' is missing")
	} else {
		date2 = dates2[0]
	}
	res2 := strings.Split(date2, "-")
	date2 = res2[2] + "-" + res2[1] + "-" + res2[0]
	fmt.Println(date2)

	if date1 == date2 {
		where = " WHERE date::DATE = '" + date1 + "'"
	} else if date1 != "" && date2 != "" {
		where = " WHERE date BETWEEN '" + date1 + "' AND '" + date2 + "'"
	} else {
		where = ""
	}

	queryForProces := query.QueryHistoryView + where
	getHistoryProses(w, r, "", "", queryForProces, "nodt")
}

func HistoryPesanDetailGet(w http.ResponseWriter, r *http.Request) {
	var nota string
	resultHistoryDetail = nil
	function.Cors(&w, r)
	not, ok := r.URL.Query()["not"]
	if !ok || len(not[0]) < 1 {
		nota = ""
		log.Println("Url Param 'nota' is missing")
	} else {
		nota = not[0]
	}
	queryForProces := query.QueryHistoryDetail
	getHistoryProses(w, r, "", nota, queryForProces, "dtl")
}

func HistoryPost(nota string, state string) error {
	resultHistory = nil
	var queryPost string
	dbPool, err := databases.Connect()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	uuid := uuid_generator.UuidGenerate()

	formatted := function.FormatTime()

	if state == "Done" {
		queryPost = query.QueryHistoryInsert(uuid.String(), formatted, nota)
	} else {
		queryPost = query.QueryHistoryBatal(uuid.String(), formatted, nota)
	}
	fmt.Println(queryPost)
	historyQueryProces(dbPool, queryPost)
	if err != nil {
		log.Fatal("error while executing query")
	}
	return err
}

//-------------------------not use-------------------------
func HistoryPesanDetail(w http.ResponseWriter, r *http.Request) {
	resultHistory = nil

	if r.Method == "GET" {

		var err error

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

		queryGet := query.QueryHistoryGetOne + "'" + key + "'"
		fmt.Println(queryGet)
		historyQueryProces(dbPool, queryGet)
		if err != nil {
			log.Fatal("error while executing query")
		}

		message := "You'r data succes to send !"
		result, err := historyDataToResult(resultHistory, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else if r.Method == "PUT" {
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&historyPesanModel); err != nil {
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

		queryUpd := " SET no_nota= '" + historyPesanModel.NoNota + "', total_price= " + strconv.Itoa(int(historyPesanModel.TotalPrice)) + "WHERE id = '" + key + "' returning *;"

		queryUpdate := query.QueryHistoryUpdate + queryUpd
		fmt.Println(queryUpdate)
		historyQueryProces(dbPool, queryUpdate)
		if err != nil {
			log.Fatal("error while executing query")
		}

		message := "The data success to save !!"
		result, err := historyDataToResult(resultHistory, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	} else if r.Method == "DELETE" {
		var err error

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

		queryGet := query.QueryHistoryDelete + "'" + key + "'" + " Returning *"
		historyQueryProces(dbPool, queryGet)
		if err != nil {
			log.Fatal("error while executing query")
		}
		// fmt.Println(queryGet)

		jsondat := &resulStatusHistory{
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

//-------------------------not use-------------------------

func HistoryCetak(w http.ResponseWriter, r *http.Request) {
	function.Cors(&w, r)
	if r.Method == "GET" {
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		key := keys[0]

		date1s, ok := r.URL.Query()["date1"]
		if !ok || len(date1s[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		date1 := date1s[0]

		date2s, ok := r.URL.Query()["date2"]
		if !ok || len(date2s[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}
		date2 := date2s[0]

		fmt.Print("sini")

		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			log.Fatal("disini", err)
		}

		var link string = "http://localhost:3737/pesanan/cetak_history.html?date1=" + date1 + "&date2=" + date2

		pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)

		page := wkhtmltopdf.NewPage(link)
		// page.FooterRight.Set("[page]")
		// page.FooterFontSize.Set(10)
		pdfg.AddPage(page)

		err = pdfg.Create()
		if err != nil {
			log.Fatal(err)
		}

		err = pdfg.WriteFile("./file-upload/pdf/history/" + key + ".pdf")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Done")

		// if err := r.ParseForm(); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// path := r.FormValue("path")
		// f, err := os.Open(path)
		// if f != nil {
		// 	defer f.Close()
		// }
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// contentDisposition := fmt.Sprintf("attachment; filename=")
		// w.Header().Set("Content-Disposition", contentDisposition)

		// if _, err := io.Copy(w, f); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
	}
}

//-------------------------not use-------------------------

func HistoryCetakFix(w http.ResponseWriter, r *http.Request) {
	function.Cors(&w, r)
	var date1, date2 string
	if r.Method == "GET" {

		dates1, ok := r.URL.Query()["date1"]
		if !ok || len(dates1[0]) < 1 {
			date1 = ""
			log.Println("Url Param 'date' is missing")
		} else {
			date1 = dates1[0]
		}
		res1 := strings.Split(date1, "-")
		date1 = res1[2] + "-" + res1[1] + "-" + res1[0]

		dates2, ok := r.URL.Query()["date2"]
		if !ok || len(dates2[0]) < 1 {
			date2 = ""
			log.Println("Url Param 'date' is missing")
		} else {
			date2 = dates2[0]
		}
		res2 := strings.Split(date2, "-")
		date2 = res2[2] + "-" + res2[1] + "-" + res2[0]

		m := pdf.NewMaroto(consts.Landscape, consts.A4)
		m.SetPageMargins(10, 10, 10)
		buildHeading(m, date1, date2)
		// data :=
		datas := getDataExport(r, date1, date2)
		// fmt.Println("setelah datas ", datas)
		buildList(m, datas)
		// fmt.Println("setelah buildlist")

		path := "/file-upload/pdf/history/" + date1 + "-" + date2 + ".pdf"
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
		contentDisposition := fmt.Sprintf("attachment; filename= history-" + date1 + "-" + date2 + ".pdf")
		w.Header().Set("Content-Disposition", contentDisposition)

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
		w.Header().Set("Content-Type", "application/pdf")
		w.Write(result)
	}
}

func getDataExport(r *http.Request, date1 string, date2 string) [][]string {
	var where string
	resultHistoryView = nil

	dbPool, err := databases.Connect()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	if date1 == date2 {
		where = " WHERE date::DATE = '" + date1 + "'"
	} else if date1 != "" && date2 != "" {
		where = " WHERE date BETWEEN '" + date1 + "' AND '" + date2 + "'"
	} else {
		where = ""
	}

	queryForProces := query.QueryHistoryView + where
	viewHistoryQueryProces(dbPool, queryForProces)
	fmt.Println("stlh query proces", resultHistoryView)
	if err != nil {
		log.Println("error while executing query")
	}
	fmt.Println("sblm change", resultHistoryView)
	datas := changeType(resultHistoryView)
	// fmt.Println(datas)
	return datas
}

func changeType(data []model_struct.ViewHistory) [][]string {
	var datas [][]string
	// fmt.Println("print data di change", data)
	for _, v := range data {
		var floatFormat string = strconv.FormatFloat(float64(v.TotalPrice), 'f', 0, 32)
		// fmt.Print(floatFormat, reflect.TypeOf(floatFormat))
		var date = strings.Split(v.Date.String(), " ")
		// fmt.Println("format date ", date[1])
		datas = append(datas, []string{
			v.NomerNota,
			v.NamaPelanggan,
			floatFormat,
			v.StatusTa,
			date[0],
		})
	}
	// fmt.Println(datas[0][1])
	return datas
}

func buildHeading(m pdf.Maroto, date1 string, date2 string) {
	// m.RegisterHeader(func() {
	m.Row(30, func() {
		m.Col(2, func() {
			_ = m.FileImage("./file-upload/SAKO_JADI.png", props.Rect{
				Percent: 80,
			})
		})
		m.Col(6, func() {
			m.Text("CAFE SAKO", props.Text{
				Top:   3,
				Size:  20,
				Style: consts.Bold,
			})
			m.Text("Jl. Cemara Raya No.836", props.Text{
				Top:  12,
				Size: 15,
			})
		})
		m.Col(3, func() {
			m.Text("Tanggal", props.Text{
				Top:   5,
				Size:  12,
				Style: consts.Bold,
			})
			m.Text(date1+" sampai "+date2, props.Text{
				Top:  12,
				Size: 12,
			})
		})
	})
	// })
}

func buildList(m pdf.Maroto, datas [][]string) {
	lightPurpleColor := getLightPurpleColor()
	tableHeadings := []string{"Nomor Nota", "Nama Pelanggan", "Total Pesanan", "Status Pesanan", "Tanggal Pembelian"}
	m.TableList(tableHeadings, datas, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      12,
			GridSizes: []uint{2, 2, 3, 3, 2},
			Style:     consts.Bold,
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{2, 2, 3, 3, 2},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})
}

func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
	var date1, date2 string
	var jsondat *failedStatus
	if r.Method == "GET" {
		dates1, ok := r.URL.Query()["date1"]
		if !ok || len(dates1[0]) < 1 {
			date1 = ""
			log.Println("Url Param 'date' is missing")
		} else {
			date1 = dates1[0]
		}
		res1 := strings.Split(date1, "-")
		date1 = res1[2] + "-" + res1[1] + "-" + res1[0]

		dates2, ok := r.URL.Query()["date2"]
		if !ok || len(dates2[0]) < 1 {
			date2 = ""
			log.Println("Url Param 'date' is missing")
		} else {
			date2 = dates2[0]
		}
		res2 := strings.Split(date2, "-")
		date2 = res2[2] + "-" + res2[1] + "-" + res2[0]

		const CONFIG_SMTP_HOST = "smtp.gmail.com"
		const CONFIG_SMTP_PORT = 587
		const CONFIG_SENDER_NAME = "CAFE SAKO - History Pesanan <hoeko10899@gmail.com>"
		const CONFIG_AUTH_EMAIL = "hoeko10899@gmail.com"
		const CONFIG_AUTH_PASSWORD = "vdszltrkhxtfugyn"

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", CONFIG_SENDER_NAME)
		mailer.SetHeader("To", "hoeko10899@gmail.com")
		mailer.SetHeader("Subject", "History Pesanan")
		mailer.SetBody("text/html", "Berikut adalah file history pesanan")
		mailer.Attach("./file-upload/pdf/history/" + date1 + "-" + date2 + ".pdf")

		dialer := gomail.NewDialer(
			CONFIG_SMTP_HOST,
			CONFIG_SMTP_PORT,
			CONFIG_AUTH_EMAIL,
			CONFIG_AUTH_PASSWORD,
		)

		err := dialer.DialAndSend(mailer)
		if err != nil {
			log.Fatal(err.Error())
			jsondat = &failedStatus{
				Succes:     "False",
				StatusCode: "400",
				Message:    "Not succes to send email !",
			}
		} else {
			log.Println("Mail sent!")
			jsondat = &failedStatus{
				Succes:     "True",
				StatusCode: "200",
				Message:    "You'r succes to send email !",
			}
			w.WriteHeader(http.StatusOK)
		}

		result, err := json.Marshal(jsondat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}
