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
- Diubah: Tidak ada perubahan logika pada commit 1. Pada commit 2 setiap handler dibungkus s.kunci() (lihat entri B5-B7).

## Deteksi bentrok per slot, PUT/DELETE, dan mutex (B5-B7)
- Prompt: sama dengan di atas, untuk commit 2 (B4-B7 dan fitur khas).
- Dipakai: unitTerpakai(), adaBentrok(), kunci() di store.go; pemeriksaan 409 di handleBuatPeminjaman(); handlePerbaruiPeminjaman() dan handleHapusPeminjaman() di handler.go; TestBersamaan di main_test.go.
- Diubah: Tidak ada perubahan logika. Saya menjalankan TestBersamaan (50 request ke slot yang sama, hasilnya tepat satu 201) dan mencoba melepas mutex untuk melihat sendiri bahwa hasilnya tidak lagi satu 201.

## Peminjaman paket, fitur khas persona D
- Prompt: sama dengan di atas.
- Dipakai: paket.go (handleBuatPaket, handleHapusPaket), rute /paket-peminjaman di main.go, TestPaketAtomik di main_test.go.
- Diubah: Tidak ada perubahan logika. Saya mempelajari tiga tahap (validasi semua item, cek bentrok semua item, baru simpan) dan alasan id paket baru dibagikan di tahap terakhir.

## Change request B: bentrok_dengan pada 409
- Prompt: "bantu saya mengerjakan checkpoint 3 langkah demi langkah sesuai requirement perubahan pada file yang saya attach" (lembar CR-B).
- Dipakai: cariPenghalang() di store.go; writeBentrok() dan perubahan handleBuatPeminjaman() di handler.go; TestBentrokDenganDaftarPenghalang di bentrok_test.go.
- Diubah: Tidak ada perubahan logika. Saya memastikan PUT dan paket tidak ikut berubah, dan menguji contoh dari lembar (slot 3-5 hasilnya [1, 2], slot 1-2 hasilnya [1]).