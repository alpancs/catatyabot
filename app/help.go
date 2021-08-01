package app

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const StartMessage = `*== cara @catatyabot membantu anda ==*

1. undang @catatyabot ke grup Telegram anda, atau boleh juga langsung _chat_ ke bot 😊
2. panggil @catatyabot menggunakan perintah berikut,
   - /catat membuat bot bersiap mencatat 📝, lalu balas dengan menuliskan satu atau beberapa catatan
   - _reply_ pesan #catatan dari bot dengan /hapus untuk menghapusnya 🗑️
   - gunakan /lihat untuk melihat daftar catatan 👀
   - dan gunakan /rangkum untuk merangkum catatan selama beberapa waktu terakhir 📈📉

satu pesan boleh berisi lebih dari satu catatan lho 🙂
contoh pesannya seperti ini 👇

_sayur kangkung 2500_
_ayam 1 kg 27k_
_susu 86 ribu_
_token listrik 200rb_
_sofa ruang tamu 6 jt_

selain itu anda juga dapat mengubah catatan yang sudah ditulis oleh bot.
cukup _reply_ pesan #catatan yang ingin diubah dengan nama & harga barang yang baru.`

func help(msg *telegram.Message) error {
	_, err := sendMessage(msg.Chat.ID, StartMessage, 0)
	return err
}
