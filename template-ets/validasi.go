package templateets

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Mengembalikan (nilai yang sudah di-trim, pesan error). Pesan kosong = valid.
func validasiTeks(nilai *string, label string, maks int) (string, string) {
	if nilai == nil {
		return "", label + " wajib diisi"
	}
	t := strings.TrimSpace(*nilai)
	if t == "" {
		return "", label + " tidak boleh kosong"
	}
	if utf8.RuneCountInString(t) > maks { // hitung per karakter, bukan byte
		return "", fmt.Sprintf("%s maksimal %d karakter", label, maks)
	}
	return t, ""
}

func validasiTanggal(nilai *string) (string, string) {
	if nilai == nil {
		return "", "tanggal wajib diisi"
	}
	t := strings.TrimSpace(*nilai)
	// time.Parse menolak tanggal yang tidak ada (30 Februari) dan menerima 29 Feb 2028
	if _, err := time.Parse("2006-01-02", t); err != nil {
		return "", "tanggal harus berformat YYYY-MM-DD dan benar-benar ada"
	}
	return t, ""
}

func validasiSlot(mulai, selesai *int) string {
	if mulai == nil || selesai == nil {
		return "slot_mulai dan slot_selesai wajib diisi"
	}
	if *mulai < 1 || *mulai > 8 || *selesai < 1 || *selesai > 8 {
		return "slot harus antara 1 sampai 8"
	}
	if *mulai > *selesai {
		return "slot_mulai tidak boleh lebih besar dari slot_selesai"
	}
	return ""
}

func validasiSumberDaya(in sumberDayaInput) (SumberDaya, string) {
	var sd SumberDaya
	var pesan string

	if sd.Nama, pesan = validasiTeks(in.Nama, "nama", 50); pesan != "" {
		return sd, pesan
	}
	if in.Jenis == nil {
		return sd, "jenis wajib diisi"
	}
	jenis := strings.TrimSpace(*in.Jenis)
	if jenis != "ruang" && jenis != "alat" {
		return sd, `jenis harus "ruang" atau "alat"`
	}
	sd.Jenis = jenis
	if in.Stok == nil {
		return sd, "stok wajib diisi"
	}
	if *in.Stok < 1 {
		return sd, "stok minimal 1"
	}
	sd.Stok = *in.Stok
	return sd, ""
}

// Mengembalikan Peminjaman yang sudah bersih (belum ada id/status) atau pesan error.
func (s *Store) validasiPeminjaman(in peminjamanInput) (Peminjaman, string) {
	var p Peminjaman
	var pesan string

	if in.SumberDayaID == nil {
		return p, "sumber_daya_id wajib diisi"
	}
	sd, ada := s.cariSumberDaya(*in.SumberDayaID)
	if !ada {
		return p, "sumber daya tidak ditemukan"
	}
	p.SumberDayaID = sd.ID

	if p.Peminjam, pesan = validasiTeks(in.Peminjam, "peminjam", 50); pesan != "" {
		return p, pesan
	}
	if p.Tanggal, pesan = validasiTanggal(in.Tanggal); pesan != "" {
		return p, pesan
	}
	if pesan = validasiSlot(in.SlotMulai, in.SlotSelesai); pesan != "" {
		return p, pesan
	}
	p.SlotMulai, p.SlotSelesai = *in.SlotMulai, *in.SlotSelesai

	if in.Jumlah == nil {
		return p, "jumlah wajib diisi"
	}
	if *in.Jumlah < 1 || *in.Jumlah > sd.Stok {
		return p, fmt.Sprintf("jumlah harus antara 1 dan %d", sd.Stok)
	}
	p.Jumlah = *in.Jumlah
	return p, ""
}
