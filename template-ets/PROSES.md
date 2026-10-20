# PROSES

## 1. Persona dan masalahnya
Pengguna persona saya adalah panitia acara yang membutuhkan beberapa hal sekaligus untuk satu acara, misalnya aula, proyektor, dan sound system. Masalahnya, kalau salah satu sudah dipinjam orang lain, acaranya tidak bisa jalan, dan panitia tidak mau aulanya sudah terpesan padahal alatnya tidak ada. Fitur khasnya adalah peminjaman paket yang sifatnya semua berhasil atau semuanya batal.

## 2. Rancangan struktur data
Saya menyimpan semuanya di dalam satu struct `Store`:
- slice `sumberDaya` (isi: struct `SumberDaya`) dan slice `peminjaman` (isi: struct `Peminjaman`). Saya memilih slice karena daftar harus terurut menurut id naik, dan `append` otomatis menjaga urutan itu. Pencarian per id cukup dengan perulangan karena datanya kecil.
- dua counter id, `idSDBerikut` dan `idPinjamBerikut`, yang dimulai dari 1 dan naik satu setiap ada data baru di jenisnya masing-masing, sesuai kontrak.
- `sync.Mutex` sudah saya siapkan di struct, tetapi belum dipakai pada commit 1. Mutex baru dipasang untuk B7 pada commit 2.
- `Peminjaman` sudah memiliki field `PaketID` sejak awal untuk fitur persona D (nilai 0 bila bukan dari paket), tetapi endpoint paketnya belum dibuat.
- Untuk membaca body JSON saya memakai struct input terpisah (`sumberDayaInput`, `peminjamanInput`) dengan field bertipe pointer.
Rencana untuk commit berikutnya: ditambah counter id paket dan penanda paket yang ada. (Catat di sini bila rancangan berubah.)

## 3. Satu keputusan desain
Saya memakai struct input dengan field pointer (`*int`, `*string`) untuk membaca body JSON, bukan struct biasa. Pilihan lain yang saya pertimbangkan adalah langsung men-decode ke struct `Peminjaman` dengan tipe biasa. Saya tidak memilihnya karena kontrak mewajibkan semua field ada, dan dengan tipe biasa, field yang tidak dikirim menjadi 0 atau string kosong sehingga tidak bisa dibedakan dari nilai yang sengaja dikirim 0 atau "". Dengan pointer, field yang tidak dikirim bernilai `nil` sehingga bisa langsung ditolak 400 dengan pesan "wajib diisi".

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
