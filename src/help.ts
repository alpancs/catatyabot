export async function sendHelpMessage(send: SendTextFn) {
    return send(`*=== CARA MENGGUNAKAN @catatyabot ===*

1. undang @catatyabot ke grup Telegram anda, atau boleh juga langsung _chat_ ke bot 😊

2. panggil bot menggunakan perintah berikut,

    - /catat membuat bot bersiap mencatat 📝, lalu balas dengan menuliskan satu atau beberapa catatan sekaligus

    - balas pesan #catatan dari bot dengan /hapus untuk menghapusnya 🗑️

    - gunakan /lihat untuk melihat daftar catatan 👀

    - dan gunakan /rangkum untuk merangkum catatan selama beberapa waktu terakhir 📈📉

satu pesan boleh berisi lebih dari satu catatan lho, dan bisa pakai satuan ribu/rb/k/juta/jt juga 🙂
contoh pesannya seperti ini 👇

_sayur kangkung 2500_
_ayam 1 kg 27k_
_susu 86 ribu_

selain itu anda juga dapat mengubah catatan yang sudah ditulis oleh bot. cukup balas pesan #catatan yang ingin diubah dengan nama & harga barang yang baru.


_ℹ️ saat ini bot sedang dalam fase pengembangan ulang dan belum berfungsi normal_`);
}
