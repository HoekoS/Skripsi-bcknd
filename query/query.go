package query

var QueryQr string = "SELECT id,no_nota,no_meja FROM qr_build"
var QueryQrGetOne string = "SELECT id,no_nota,no_meja FROM qr_build where id = "
var QueryQrInsert string = "INSERT INTO public.qr_build( no_nota, no_meja) "
var QueryQrUpdate string = "UPDATE public.qr_build "
var QueryQrDelete string = "DELETE FROM public.qr_build WHERE id= "

var QueryMenu string = "SELECT id, name, description, price, pic, created_at,kategori,sub_kategori FROM menu"
var QueryMenuGetOne string = "SELECT id, name, description, price, pic, created_at, kategori,sub_kategori FROM menu where id = "
var QueryMenuInsert string = "INSERT INTO menu(id, name, description, price, pic, created_at,kategori,sub_kategori)"
var QueryMenuUpdate string = "UPDATE public.menu "
var QueryMenuDelete string = "DELETE FROM menu WHERE id='"
var QueryMenuMakanan string = "SELECT id, name, description, price, pic, created_at,kategori FROM menu WHERE kategori='Makanan'"
var QueryMenuMinuman string = "SELECT id, name, description, price, pic, created_at,kategori FROM menu WHERE kategori='Minuman'"

var QueryDetailPesanan string = "SELECT id, no_nota, id_menu, quantity, catatan FROM detail_pesanan"
var QueryDetailPesananGetOne string = "SELECT id, no_nota, id_menu, quantity FROM detail_pesanan where id = "
var QueryDetailPesananInsert string = "INSERT INTO public.detail_pesanan(no_nota, id_menu, quantity,catatan)"
var QueryDetailPesananUpdate string = "UPDATE public.detail_pesanan "
var QueryDetailPesananDelete string = "DELETE FROM public.detail_pesanan WHERE id= "
var QueryDetailPesananNota string = `SELECT dp.id,dp.no_nota,m.name,dp.quantity,dp.catatan,m.price,(dp.quantity*m.price)as total_harga, m.pic FROM detail_pesanan as dp left join menu as m on m.id = dp.id_menu`
var QueryTotalHarga string = `
SELECT no_nota,sum(total_harga) as total_semua
FROM(
SELECT dp.id,dp.no_nota,m.name,dp.quantity,dp.catatan,m.price,(dp.quantity*m.price)as total_harga 
FROM detail_pesanan as dp left join menu as m on m.id = dp.id_menu
) as fb`

var QueryHistory string = "SELECT * FROM history_pesanan"
var QueryHistoryGetOne string = "SELECT * FROM history_pesanan where id = "

// var QueryHistoryInsert string = "INSERT INTO public.history_pesanan(id, no_nota, total_price, date)"
var QueryHistoryUpdate string = "UPDATE public.history_pesanan"
var QueryHistoryDelete string = "DELETE FROM history_pesanan WHERE id= "
var QueryHistoryView = `
SELECT *
FROM(
	SELECT hp.no_nota,p.nama_pelanggan,hp.total_price,
        Case 
                WHEN p.flag_take_away = 1 THEN 'Take Away'
                ELSE 'Dine In'
		END AS status_ta,
		hp.date::DATE
	FROM history_pesanan AS hp
	JOIN pesanan AS p ON p.no_nota = hp.no_nota  
) AS fx `
var QueryHistoryDetail = `
select hp.no_nota,m.name,m.price,dp.quantity,(m.price*dp.quantity) as total_price,hp.total_price as total_price_pesanan,dp.catatan
from history_pesanan as hp
join detail_pesanan as dp on hp.no_nota = dp.no_nota
join menu as m on m.id = dp.id_menu`

func QueryHistoryInsert(id string, date string, nota string) string {
	var query string = `
		INSERT INTO public.history_pesanan(id, no_nota, total_price, date)
			SELECT '` + id + `',no_nota,sum(total_harga) as total_semua,'` + date + `'
			FROM(
			SELECT dp.id,dp.no_nota,m.name,dp.quantity,dp.catatan,m.price,(dp.quantity*m.price)as total_harga
			FROM detail_pesanan as dp left join menu as m on m.id = dp.id_menu
			) as fb
			Where no_nota ='` + nota + `'
			group by no_nota`
	return query
}
func QueryHistoryBatal(id string, date string, nota string) string {
	var query string = `
		INSERT INTO public.history_pesanan(id, no_nota, total_price, date)
			SELECT '` + id + `','` + nota + `',0 as total_semua,'` + date + `'`
	return query
}

var QueryPesanan string = "SELECT * FROM public.pesanan"
var QueryPesananGetOne string = "SELECT id, flag_take_away, state, date,no_nota, nama_pelanggan, no_meja FROM public.pesanan where id = "
var QueryPesananGetOneNot string = "SELECT * FROM public.pesanan where no_nota = "
var QueryPesananInsert string = "INSERT INTO public.pesanan(id, flag_take_away, state, date, nama_pelanggan, no_meja,url) "
var QueryPesananUpdate string = "UPDATE public.pesanan "
var QueryPesananDelete string = "DELETE FROM public.pesanan WHERE id= "
var QueryPesananNoDone string = "SELECT * FROM public.pesanan where state != 'Done'"
var QueryPesananDapur string = "SELECT * FROM pesanan WHERE state = 'VERIFIKASI STAFF KASIR' OR state = 'SEDANG DIMASAK' "

func QrNPesananInsert(values string, uuid string, no_meja string) string {
	query := `
	Begin;
		INSERT INTO public.pesanan(id,state,nama_pelanggan,date,flag_take_away) ` + values + `;
		INSERT INTO public.qr_build(no_nota, no_meja)
			Select no_nota,` + no_meja + ` as no_meja from pesanan where id='` + uuid + `' RETURNING *;
	COMMIT;`
	return query
}

var QueryBahanMakanan string = `
SELECT 
	id,
	description,
	create_at,
	quantity_tambah-quantity_kurang as quantity,
	satuan
FROM(
	SELECT 
	b.id,	
	description,
		quantity,
		CASE
			WHEN bt.quantity_tambah is NULL 
				THEN 0
			ELSE bt.quantity_tambah
		END AS quantity_tambah,
		CASE
			WHEN bk.quantity_kurang is NULL 
				THEN 0
			ELSE bk.quantity_kurang
		END AS quantity_kurang,
		b.satuan,
		b.create_at
	FROM public.bahan_makanan as b
	LEFT JOIN(
		SELECT id_bahan_makanan,sum(quantity) AS quantity_tambah
		FROM public.bahan_makanan_history_penambahan
		GROUP BY id_bahan_makanan
	) as bt on b.id = bt.id_bahan_makanan
	LEFT JOIN(
		SELECT id_bahan_makanan,sum(quantity) AS quantity_kurang
		FROM public.bahan_makanan_history_pengurangan
		GROUP BY id_bahan_makanan
	)as bk on b.id = bk.id_bahan_makanan
) as fx `
var QueryBahanMakananInsert string = "INSERT INTO public.bahan_makanan(id, description, create_at,satuan) "
var QueryBahanMakananUpdate string = "UPDATE public.bahan_makanan "
var QueryBahanMakananDelete string = "DELETE FROM public.bahan_makanan WHERE id= '"

func QueryBahanMakananUPdateQuantity(id string, date string) string {
	if date == "" {
		date = ""
	} else {
		date = `,create_at='` + date + `'`
	}
	output := `
	UPDATE public.bahan_makanan
	SET quantity=fixQuery.quantity` + date + `
	FROM(
		SELECT quantity_tambah-quantity_kurang as quantity
		FROM(
				SELECT
					b.id,
					quantity,
					CASE
						WHEN bt.quantity_tambah is NULL 
							THEN 0
						ELSE bt.quantity_tambah
					END AS quantity_tambah,
					CASE
						WHEN bk.quantity_kurang is NULL 
							THEN 0
						ELSE bk.quantity_kurang
					END AS quantity_kurang
				FROM public.bahan_makanan as b
				LEFT JOIN(
						SELECT id_bahan_makanan,sum(quantity) AS quantity_tambah
						FROM public.bahan_makanan_history_penambahan
						GROUP BY id_bahan_makanan
				) as bt on b.id = bt.id_bahan_makanan
				LEFT JOIN(
						SELECT id_bahan_makanan,sum(quantity) AS quantity_kurang
						FROM public.bahan_makanan_history_pengurangan
						GROUP BY id_bahan_makanan
				)as bk on b.id = bk.id_bahan_makanan
		) as fx 
		WHERE id='` + id + `'
	) AS fixQuery
	WHERE id='` + id + `';`
	return output
}

var QueryBahanMakananKurang string = "SELECT id,id_bahan_makanan,quantity,create_at FROM public.bahan_makanan_history_pengurangan"
var QueryBahanMakananKurangInsert string = "INSERT INTO public.bahan_makanan_history_pengurangan(id_bahan_makanan, quantity, create_at) "
var QueryBahanMakananKurangUpdate string = "UPDATE public.bahan_makanan_history_pengurangan "
var QueryBahanMakananKurangDelete string = "DELETE FROM public.bahan_makanan_history_pengurangan WHERE id= '"
var QueryBahanMakananKurangDeleteWithId string = "DELETE FROM public.bahan_makanan_history_pengurangan WHERE id_bahan_makanan= '"

var QueryBahanMakananTambah string = "SELECT id,id_bahan_makanan,quantity,create_at FROM public.bahan_makanan_history_penambahan"
var QueryBahanMakananTambahInsert string = "INSERT INTO public.bahan_makanan_history_penambahan(id_bahan_makanan, quantity, create_at) "
var QueryBahanMakananTambahUpdate string = "UPDATE public.bahan_makanan_history_penambahan "
var QueryBahanMakananTambahDelete string = "DELETE FROM public.bahan_makanan_history_penambahan WHERE id= '"
var QueryBahanMakananTambahDeleteWithId string = "DELETE FROM public.bahan_makanan_history_penambahan WHERE id_bahan_makanan= '"

var QuerySchedule string = "SELECT * FROM public.schedule"
var QueryScheduleInsert string = "INSERT INTO public.schedule(id, day, start_at, end_at, description)"
var QueryScheduleUpdate string = "UPDATE public.schedule "
var QueryScheduleDelete string = "DELETE FROM public.schedule WHERE id= '"
var QueryScheduleView string = `
select e.full_name, s.day, s.description
from employee as e
inner join employee_schedule as es on es.employee_id = e.id
inner join schedule as s on es.schedule_id = s.id`

var QueryEmployee string = "SELECT id, full_name, gender, birth_date::Varchar, birth_place, address, email, phone_number, staff_status, staff_position, staff_active_status, staff_date_join::Varchar FROM public.employee "
var QueryEmployeeInsert string = "INSERT INTO public.employee(id, full_name, gender, birth_date, birth_place, address, email, phone_number, staff_status, staff_position, staff_active_status, staff_date_join)"
var QueryEmployeeUpdate string = "UPDATE public.employee "
var QueryEmployeeDelete string = "DELETE FROM public.employee WHERE id= '"

var QueryEmployeeSchedule string = "SELECT * FROM public.employee_schedule"
var QueryEmployeeScheduleInsert string = "INSERT INTO public.employee_schedule(employee_id, schedule_id)"
var QueryEmployeeScheduleUpdate string = "UPDATE public.employee_schedule "
var QueryEmployeeScheduleDelete string = "DELETE FROM public.employee_schedule WHERE id= "
var QueryEmployeeFromId string = `
select es.id,s.day,s.start_at,s.end_at,s.description
from employee as e
left join employee_schedule as es on es.employee_id = e.id
left join schedule as s on es.schedule_id = s.id 
WHERE e.id = '`

var QueryUserLogin string = "SELECT id, username, password, email, role_code, created_on FROM public.user_login"
var QueryUserLoginInsert string = "INSERT INTO public.user_login(username, password, email, role_code, created_on)"
var QueryUserLoginUpdate string = "UPDATE public.user_login "
var QueryUserLoginDelete string = "DELETE FROM public.user_login WHERE id= '"
var QueryUserSearch string = "SELECT id, username, password, email, role_code, created_on FROM public.user_login WHERE email ="

var QueryDashboardKosong = `select 30-count(*) as jumlah_meja_kosong from pesanan
where not state='Cetak QR' AND not state='Done'`
var QueryDashboardSelesai = `select count(*) as jumlah_pesanan_selesai from pesanan
where state='Done'`
var QueryDashboardBatal = `
select count(*) as jumlah_pesanan_batal from pesanan
where state='BATAL PESAN'`
var QueryJmlhBahanHabis = `SELECT COUNT(*)
FROM(
	SELECT 
		id,
		description,
		create_at,
		quantity_tambah-quantity_kurang as quantity,
		satuan
	FROM(
		SELECT 
		b.id,	
		description,
			quantity,
			CASE
				WHEN bt.quantity_tambah is NULL 
					THEN 0
				ELSE bt.quantity_tambah
			END AS quantity_tambah,
			CASE
				WHEN bk.quantity_kurang is NULL 
					THEN 0
				ELSE bk.quantity_kurang
			END AS quantity_kurang,
			b.satuan,
			b.create_at
		FROM public.bahan_makanan as b
		LEFT JOIN(
			SELECT id_bahan_makanan,sum(quantity) AS quantity_tambah
			FROM public.bahan_makanan_history_penambahan
			GROUP BY id_bahan_makanan
		) as bt on b.id = bt.id_bahan_makanan
		LEFT JOIN(
			SELECT id_bahan_makanan,sum(quantity) AS quantity_kurang
			FROM public.bahan_makanan_history_pengurangan
			GROUP BY id_bahan_makanan
		)as bk on b.id = bk.id_bahan_makanan
	) as fx 
) as fxb
where quantity=0`
var QueryEmployeeHari = `select e.full_name, e.staff_status, e.staff_position, s.start_at ::varchar,s.end_at::varchar,s.description from employee_schedule as es
left join employee as e on e.id = es.employee_id
left join schedule as s on s.id = es.schedule_id
where e.staff_active_status = 'Aktif' `
var QueryMenuDash = `select * from menu
order by created_at DESC
limit 8`
var QueryBahanDash = `SELECT 
id,
description,
create_at,
quantity_tambah-quantity_kurang as quantity,
satuan
FROM(
SELECT 
b.id,	
description,
	quantity,
	CASE
		WHEN bt.quantity_tambah is NULL 
			THEN 0
		ELSE bt.quantity_tambah
	END AS quantity_tambah,
	CASE
		WHEN bk.quantity_kurang is NULL 
			THEN 0
		ELSE bk.quantity_kurang
	END AS quantity_kurang,
	b.satuan,
	b.create_at
FROM public.bahan_makanan as b
LEFT JOIN(
	SELECT id_bahan_makanan,sum(quantity) AS quantity_tambah
	FROM public.bahan_makanan_history_penambahan
	GROUP BY id_bahan_makanan
) as bt on b.id = bt.id_bahan_makanan
LEFT JOIN(
	SELECT id_bahan_makanan,sum(quantity) AS quantity_kurang
	FROM public.bahan_makanan_history_pengurangan
	GROUP BY id_bahan_makanan
)as bk on b.id = bk.id_bahan_makanan
) as fx 
where quantity=0`
