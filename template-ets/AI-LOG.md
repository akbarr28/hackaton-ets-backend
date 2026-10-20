# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Model data dan penyimpanan di memori (B1-B3)
- Prompt: "bantu saya mengerjakan persona D (panitia acara) sesuai kontrak di requirement, pisahkan jawaban per commit, commit 1 untuk B1-B3"
- Dipakai: struct SumberDaya, Peminjaman, sumberDayaInput, peminjamanInput di model.go; Store, NewStore(), cariSumberDaya(), tambahSumberDaya(), indeksPeminjaman(), tambahPeminjaman() di store.go
- Diubah: <isi sesuai kenyataan, lihat pilihan di bawah>

## Validasi input (B4 sebagian)
- Prompt: sama dengan di atas (kode validasi ikut dihasilkan dalam jawaban yang sama)
- Dipakai: validasiTeks(), validasiTanggal(), validasiSlot(), validasiSumberDaya(), validasiPeminjaman() di validasi.go
- Diubah: <isi sesuai kenyataan>

## Endpoint sumber daya dan peminjaman (B1-B3)
- Prompt: sama dengan di atas
- Dipakai: writeJSON(), writeError(), decodeBody(), parseID(), handleBuatSumberDaya(), handleDaftarSumberDaya(), handleDetailSumberDaya(), handleBuatPeminjaman(), handleDaftarPeminjaman(), handleDetailPeminjaman() di handler.go
- Diubah: <isi sesuai kenyataan>

## Routing
- Prompt: sama dengan di atas
- Dipakai: NewHandler() dan main() di main.go
- Diubah: <isi sesuai kenyataan>