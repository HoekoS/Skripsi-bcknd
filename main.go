package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"skripsi.com/backend-adm/bahan_makanan"
	"skripsi.com/backend-adm/dashboard"
	"skripsi.com/backend-adm/employee"
	"skripsi.com/backend-adm/login"
	"skripsi.com/backend-adm/menu_makanan"
	"skripsi.com/backend-adm/pesanan"
	"skripsi.com/backend-adm/qr_code"
	"skripsi.com/backend-adm/schedule"
)

// var qrModel *model_struct.QrModel = new(model_struct.QrModel)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/signup", login.UserSignUp)
	mux.HandleFunc("/signin", login.UserCekLogin)
	mux.HandleFunc("/signout", login.UserSignOut)
	mux.HandleFunc("/user", login.UserLoginGet)
	mux.HandleFunc("/user/delete", login.UserLoginDelete)

	mux.HandleFunc("/menu", menu_makanan.MenuAll)
	mux.HandleFunc("/menu/post", menu_makanan.MenuPost)
	mux.HandleFunc("/menu/delete", menu_makanan.MenuDelete)
	mux.HandleFunc("/menu/update", menu_makanan.MenuEdit)
	mux.HandleFunc("/menu/makan", menu_makanan.ViewMenuMakanan)
	mux.HandleFunc("/menu/minum", menu_makanan.ViewMenuMinuman)
	mux.HandleFunc("/menu/kategori", menu_makanan.ViewKategori)

	mux.HandleFunc("/qr", qr_code.QrGetAll)
	mux.HandleFunc("/qr/post", qr_code.QrPost)
	mux.HandleFunc("/qr/delete", qr_code.QrDelete)
	mux.HandleFunc("/qr/update", qr_code.QrUpdate)
	// http.HandleFunc("/qr/detail/", qr_code.QrDetail)

	mux.HandleFunc("/order", pesanan.PesananGetAll)
	mux.HandleFunc("/order/nota", pesanan.PesananGetOneByNot)
	mux.HandleFunc("/order/post", pesanan.Pesanan)
	mux.HandleFunc("/order/qr/cetak", pesanan.QrCetakFix)

	mux.HandleFunc("/order/pilih", pesanan.DetailPesananGet)
	mux.HandleFunc("/order/get", pesanan.DetailPesananGetOne)
	mux.HandleFunc("/order/get/harga", pesanan.TotalHargaGet)
	mux.HandleFunc("/order/pilih/post", pesanan.DetailPesananData)
	mux.HandleFunc("/order/pilih/update", pesanan.DetailPesananUpdate)
	mux.HandleFunc("/order/pilih/delete", pesanan.DetailPesananDelete)
	mux.HandleFunc("/order/pilih/state", pesanan.DetailPesananUpdateState)
	mux.HandleFunc("/order/pilih/state/nota", pesanan.DetailPesananUpdateStateNot)
	mux.HandleFunc("/order/get/dapur", pesanan.PesananGetForDapur)
	mux.HandleFunc("/order/detail/all", pesanan.DetailPesananGetAll)

	mux.HandleFunc("/history", pesanan.HistoryPesanGet)
	mux.HandleFunc("/history/detail", pesanan.HistoryPesanDetailGet)
	mux.HandleFunc("/history/cetak", pesanan.HistoryCetakFix)
	mux.HandleFunc("/history/email", pesanan.SendEmail)

	mux.HandleFunc("/bahan", bahan_makanan.BahanMakananGet)
	mux.HandleFunc("/bahan/input", bahan_makanan.BahanMakananPost)
	mux.HandleFunc("/bahan/update", bahan_makanan.BahanMakananEdit)
	mux.HandleFunc("/bahan/delete", bahan_makanan.BahanMakananDelete)

	mux.HandleFunc("/bahan/kurang", bahan_makanan.BahanMakananKurangGet)
	mux.HandleFunc("/bahan/kurang/input", bahan_makanan.BahanMakananKurangPost)
	mux.HandleFunc("/bahan/kurang/update", bahan_makanan.BahanMakananKurangEdit)
	mux.HandleFunc("/bahan/kurang/delete", bahan_makanan.BahanMakananKurangDelete)

	mux.HandleFunc("/bahan/tambah", bahan_makanan.BahanMakananTambahGet)
	mux.HandleFunc("/bahan/tambah/input", bahan_makanan.BahanMakananTambahPost)
	mux.HandleFunc("/bahan/tambah/update", bahan_makanan.BahanMakananTambahEdit)
	mux.HandleFunc("/bahan/tambah/delete", bahan_makanan.BahanMakananTambahDelete)

	mux.HandleFunc("/schedule", schedule.ScheduleGet)
	mux.HandleFunc("/schedule/view", schedule.ScheduleEmployeeView)
	mux.HandleFunc("/schedule/input", schedule.SchedulePost)
	mux.HandleFunc("/schedule/update", schedule.ScheduleEdit)
	mux.HandleFunc("/schedule/delete", schedule.ScheduleDelete)

	mux.HandleFunc("/employee", employee.EmployeeGet)
	mux.HandleFunc("/employee/input", employee.EmployeePost)
	mux.HandleFunc("/employee/update", employee.EmployeeEdit)
	mux.HandleFunc("/employee/delete", employee.EmployeeDelete)

	mux.HandleFunc("/employee/schedule", employee.EmployeeScheduleGet)
	mux.HandleFunc("/employee/schedule/get", schedule.EmployeeScheduleGetWithId)
	mux.HandleFunc("/employee/schedule/input", employee.EmployeeSchedulePost)
	mux.HandleFunc("/employee/schedule/update", employee.EmployeeScheduleEdit)
	mux.HandleFunc("/employee/schedule/delete", employee.EmployeeScheduleDelete)

	mux.HandleFunc("/dashboard/kosong", dashboard.JmlKosongDashboardGet)
	mux.HandleFunc("/dashboard/selesai", dashboard.JmlSelesaiDashboardGet)
	mux.HandleFunc("/dashboard/batal", dashboard.JmlBatalDashboardGet)
	mux.HandleFunc("/dashboard/habis", dashboard.JmlHabisDashboardGet)
	mux.HandleFunc("/dashboard/shift", dashboard.ShiftDashboardGet)
	mux.HandleFunc("/dashboard/menu", menu_makanan.MenuDashboard)
	mux.HandleFunc("/dashboard/bahan", bahan_makanan.BahanMakananDashGet)

	var handler http.Handler = mux
	// handler = MiddlewareAuth(handler)
	// handler = MiddlewareAllowOnlyGet(handler)
	handler = cors.Default().Handler(mux)

	log.Println("Listening...")

	http.ListenAndServe(":9000", handler)
}
