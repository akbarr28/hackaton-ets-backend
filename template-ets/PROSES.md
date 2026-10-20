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
- Gejala: Saat 50 request POST /peminjaman ke slot yang sama dikirim bersamaan tanpa mutex, hasilnya tidak konsisten. Kadang lebih dari satu request mendapat 201 padahal stok hanya 1, dan go test -race melaporkan DATA RACE.
- Cara menemukan: Saya menjalankan TestBersamaan setelah sengaja melepas pembungkus s.kunci() dari rute di main.go, untuk membuktikan sendiri kenapa mutex diperlukan. Karena penjadwalan goroutine tidak pasti, hasilnya tidak selalu sama, jadi saya menjalankannya berulang kali.
- Penyebab dan solusi: Cek bentrok (membaca slice) dan penyimpanan (append) adalah dua langkah terpisah. Beberapa goroutine bisa sama-sama lolos cek sebelum ada yang sempat menyimpan. Solusinya, kunci() membungkus handler dengan satu sync.Mutex sehingga cek bentrok dan penyimpanan terjadi dalam satu kunci yang sama.


### Bug 2
- Gejala: Pada contoh lembar change request (id 1 slot 2-3, id 2 slot 5-6), permintaan slot 3-5 menghasilkan bentrok_dengan [1], padahal seharusnya [1, 2]. Peminjaman id 2 tidak ikut karena slot 5 hanya "menyentuh" batasnya.
- Cara menemukan: TestBentrokDenganDaftarPenghalang gagal dengan pesan "bentrok_dengan harus [1 2], dapat [1]" ketika kondisi irisan di cariPenghalang saya tulis dengan tanda < (bukan <=).
- Penyebab dan solusi: Rentang slot bersifat inklusif (3-5 berarti slot 3, 4, dan 5), jadi dua rentang [a,b] dan [c,d] beririsan bila a <= d dan b >= c. Dengan < batas yang sama terlewat (off-by-one). Saya ganti dengan <=, dan tetap memastikan rentang yang hanya bersebelahan (1-2 dan 3-4) tidak dianggap bentrok.

## 5. Change request
Tim frontend meminta agar respons 409 pada POST /peminjaman karena slot atau kapasitas penuh memuat field baru `bentrok_dengan`, berisi id peminjaman yang menghalangi, supaya pengguna tahu siapa yang harus dihubungi.
Bagian kode yang saya ubah: saya menambah fungsi `cariPenghalang()` di store.go, helper `writeBentrok()` di handler.go, dan mengganti satu pemanggilan `writeError` di `handleBuatPeminjaman` dengan `writeBentrok`. PUT, paket, dan semua respons 400/404 tidak saya sentuh.
Keputusan desain: saya membuat fungsi baru `cariPenghalang()` yang terpisah dari `adaBentrok()`. Pilihan lain adalah mengubah `adaBentrok()` supaya sekaligus mengembalikan daftar id. Saya tidak memilihnya karena `adaBentrok()` juga dipakai PUT dan paket yang tidak boleh berubah responsnya, dan mengubah signature-nya berisiko merusak fitur lama. Dengan fungsi terpisah, `bentrok_dengan` hanya dibuat di satu tempat sehingga tidak bisa bocor ke respons lain.
