package telegram

var (
	StartMsg = `*Cara Catatan Belanja Bot Membantu Anda*
- Undang @CatatanBelanjaBot ke grup Telegram keluarga anda
- Bot mencatat pengeluaran Anda, ketika ada pesan seperti
  - @CatatanBelanjaBot bahan masakan 45.000
  - @CatatanBelanjaBot tagihan listrik 200 k
  - @CatatanBelanjaBot baju lebaran 1,5jt
  - @CatatanBelanjaBot jus jambu 8rb, enak banget 😆
- Bot juga memiliki beberapa perintah, yaitu
  - /Rangkuman: rangkuman catatan belanja
  - /HariIni: daftar belanjaan hari ini
  - /Kemarin: daftar belanjaan kemarin
  - /PekanIni: daftar belanjaan pekan ini
  - /PekanLalu: daftar belanjaan pekan lalu
  - /BulanIni: daftar belanjaan bulan ini
  - /BulanLalu: daftar belanjaan bulan lalu
  - /Hapus: menghapus catatan. _reply_ ke catatan yang ingin dihapus`
)

func responseStart(u Update) (*Response, error) {
	res := DefaultResponse
	res.ChatID = u.Message.Chat.ID
	res.Text = StartMsg
	return &res, nil
}
