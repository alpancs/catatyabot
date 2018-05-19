package telegram

import (
	"os"
)

type Update struct {
	Message struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID int
		}
		Text string
	}
}

type Response struct {
	ChatID           int    `json:"chat_id"`
	Text             string `json:"text"`
	ParseMode        string `json:"parse_mode,omitempty"`
	ReplyToMessageId int    `json:"reply_to_message_id,omitempty"`
}

var (
	Username        = os.Getenv("TELEGRAM_BOT_USERNAME")
	DefaultResponse = Response{ParseMode: "Markdown"}
	ConfusingMsgs   = []string{"ngomong apa to bos? 🤔", "mbuh bos, gak ngerti 😒", "aku orak paham boooss 😔"}
	StartMsg        = `*Cara Catatan Belanja Bot Membantu Anda*
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

func Respond(u Update) (*Response, error) {
	switch u.Message.Text {
	case "/start":
		return responseStart(u)
	default:
		return responseConfusing(u)
	}
}

func responseStart(u Update) (*Response, error) {
	res := DefaultResponse
	res.ChatID = u.Message.Chat.ID
	res.Text = StartMsg
	return &res, nil
}

func responseConfusing(u Update) (*Response, error) {
	res := DefaultResponse
	res.ChatID = u.Message.Chat.ID
	res.Text = sample(ConfusingMsgs)
	res.ReplyToMessageId = u.Message.MessageID
	return &res, nil
}
