package app

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const StartMessage = `*Cara @catatyabot membantu anda:*
- undang @catatyabot ke grup Telegram keluarga anda
- /catat untuk memanggil bot supaya bersiap mencatat 📝
- /hapus untuk menghapus catatan 🗑️
- /lihat untuk melihat catatan 👀
- /rangkum untuk melihat rangkuman catatan 💰

Selain itu, anda juga dapat mengubah catatan dengan cara membalas/_reply_ ke pesan #catatan yang ingin diubah, lalu sebutkan nama barang serta harga barang yang baru.`

func commandStart(msg *telegram.Message) (bool, error) {
	if msg.Command() != "start" {
		return false, nil
	}

	_, err := sendMessage(msg.Chat.ID, StartMessage, 0)
	return true, err
}
