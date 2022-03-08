package model_struct

import "time"

type QrModel struct {
	Id        int    `json:"id"`
	NomerNota string `json:"no_nota"`
	NomerMeja int    `json:"no_meja"`
}

type Menu struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Pic         string    `json:"pic"`
	CreateAt    time.Time `json:"create_at"`
	Kategori    string    `json:"kategori"`
	SubKategori string    `json:"sub_kategori"`
}

type Pesanan struct {
	Id            string    `json:"id"`
	NomerNota     string    `json:"no_nota"`
	NomerMeja     int       `json:"no_meja"`
	NamaPelanggan string    `json:"nama_pelanggan"`
	FlagTakeAway  int       `json:"flag_ta"`
	State         string    `json:"State"`
	Url           string    `json:"Url"`
	DeleteStatus  int       `json:"delete_status"`
	Date          time.Time `json:"date"`
}

type DetailPesanan struct {
	Id       int    `json:"id"`
	NoNota   string `json:"no_nota"`
	IdMenu   string `json:"id_menu"`
	Quantity int    `json:"quantity"`
	Catatan  string `json:"catatan"`
}

type HistoryPesan struct {
	Id         string    `json:"id"`
	NoNota     string    `json:"no_nota"`
	TotalPrice float64   `json:"total_price"`
	Date       time.Time `json:"date"`
}

type QrNPesanan struct {
	Id            string    `json:"id"`
	NomerNota     string    `json:"no_nota"`
	NomerMeja     int       `json:"no_meja"`
	NamaPelanggan string    `json:"nama_pelanggan"`
	FlagTakeAway  int       `json:"flag_ta"`
	State         string    `json:"State"`
	Url           string    `json:"Url"`
	Date          time.Time `json:"date"`
}

type ViewPesan struct {
	Id            int     `json:"id"`
	NomerNota     string  `json:"no_nota"`
	NamaPelanggan string  `json:"nama_pelanggan"`
	Quantity      int     `json:"quantity"`
	Catatan       string  `json:"catatan"`
	Price         float64 `json:"price"`
	TotalHarga    float64 `json:"total_harga"`
	Url           string  `json:"url"`
}

type TotalHargaPesanan struct {
	NomerNota  string  `json:"no_nota"`
	TotalHarga float64 `json:"total_harga"`
}

type ViewHistory struct {
	NomerNota     string    `json:"no_nota"`
	NamaPelanggan string    `json:"nama_pelanggan"`
	TotalPrice    float64   `json:"total_price"`
	StatusTa      string    `json:"status_ta"`
	Date          time.Time `json:"date"`
}

type ViewHistoryDetail struct {
	NomerNota         string  `json:"no_nota"`
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	TotalPrice        float64 `json:"total_price"`
	TotalSemuaPesanan float64 `json:"total_price_pesanan"`
	Catatan           string  `json:"catatan"`
}

type BahanMakanan struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Satuan      string    `json:"satuan"`
	CreateAt    time.Time `json:"create_at"`
}

type BahanMakananTambah struct {
	Id       int       `json:"id"`
	IdBahan  string    `json:"id_bahan_makanan"`
	Quantity int       `json:"quantity"`
	CreateAt time.Time `json:"create_at"`
}

type BahanMakananKurang struct {
	Id       int       `json:"id"`
	IdBahan  string    `json:"id_bahan_makanan"`
	Quantity int       `json:"quantity"`
	CreateAt time.Time `json:"create_at"`
}

type UserLogin struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	RoleCode  string    `json:"role_code"`
	CreatedOn time.Time `json:"created_on"`
}

type Employee struct {
	Id                string `json:"id"`
	FullName          string `json:"full_name"`
	Gender            string `json:"gender"`
	BirthDate         string `json:"birth_date"`
	BirthPlace        string `json:"birth_place"`
	Address           string `json:"address"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
	StaffStatus       string `json:"staff_status"`
	StaffPosition     string `json:"staff_position"`
	StaffActiveStatus string `json:"staff_active"`
	StaffDateJoin     string `json:"staff_date_join"`
}

type Schedule struct {
	Id          int       `json:"id"`
	Day         string    `json:"day"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Description string    `json:"description"`
}

type EmployeeSchedule struct {
	Id         int    `json:"id"`
	EmployeeId string `json:"employee_id"`
	ScheduleId int    `json:"schedule_id"`
}

type ViewSchedule struct {
	FullName    string `json:"full_name"`
	Day         string `json:"day"`
	Description string `json:"description"`
}

type JumlahDashboard struct {
	Jumlah int `json:"jumlah"`
}

type JadwalShiftDash struct {
	Nama          string `json:"nama"`
	StaffStatus   string `json:"staff_status"`
	StaffPosition string `json:"staff_position"`
	StartAt       string `json:"start_at"`
	EndAt         string `json:"end_at"`
	Description   string `json:"description"`
}
