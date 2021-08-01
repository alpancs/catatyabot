package app

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const HelpMessage = `*== cara @catatyabot membantu anda ==*

1. undang @catatyabot ke grup Telegram anda, atau boleh juga langsung _chat_ ke bot 😊

2. panggil bot menggunakan perintah berikut,

    - /catat membuat bot bersiap mencatat 📝, lalu balas dengan menuliskan satu atau beberapa catatan sekaligus

    - _reply_ pesan #catatan dari bot dengan /hapus untuk menghapusnya 🗑️

    - gunakan /lihat untuk melihat daftar catatan 👀

    - dan gunakan /rangkum untuk merangkum catatan selama beberapa waktu terakhir 📈📉

satu pesan boleh berisi lebih dari satu catatan lho, dan bisa pakai satuan ribu/rb/k/juta/jt juga 🙂
contoh pesannya seperti ini 👇

_sayur kangkung 2500_
_ayam 1 kg 27k_
_susu 86 ribu_

selain itu anda juga dapat mengubah catatan yang sudah ditulis oleh bot. cukup _reply_ pesan #catatan yang ingin diubah dengan nama & harga barang yang baru.`

func help(msg *telegram.Message) (processed bool, err error) {
	cmd := msg.Command()
	if cmd == "start" || cmd == "bantuan" {
		_, err = sendMessage(msg.Chat.ID, HelpMessage, 0)
		return true, err
	}

	return false, nil
}
