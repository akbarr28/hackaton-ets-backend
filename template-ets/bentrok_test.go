package templateets

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestBentrokDenganDaftarPenghalang(t *testing.T) {
	h := NewHandler()
	kirim(h, "POST", "/sumber-daya", `{"nama":"Aula A","jenis":"ruang","stok":1}`)
	// id 1 (slot 2-3) dan id 2 (slot 5-6) di 2026-11-02
	kirim(h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"A","tanggal":"2026-11-02","slot_mulai":2,"slot_selesai":3,"jumlah":1}`)
	kirim(h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"B","tanggal":"2026-11-02","slot_mulai":5,"slot_selesai":6,"jumlah":1}`)
	// id 3: tanggal lain, tidak boleh ikut daftar
	kirim(h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"C","tanggal":"2026-11-03","slot_mulai":3,"slot_selesai":5,"jumlah":1}`)

	cek := func(body string, harapan []int) {
		t.Helper()
		rec := kirim(h, "POST", "/peminjaman", body)
		if rec.Code != http.StatusConflict {
			t.Fatalf("harus 409, dapat %d", rec.Code)
		}
		var resp struct {
			Error         string `json:"error"`
			BentrokDengan []int  `json:"bentrok_dengan"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		if resp.Error == "" {
			t.Fatal("field error harus tetap ada")
		}
		if !slices.Equal(resp.BentrokDengan, harapan) {
			t.Fatalf("bentrok_dengan harus %v, dapat %v", harapan, resp.BentrokDengan)
		}
	}

	cek(`{"sumber_daya_id":1,"peminjam":"X","tanggal":"2026-11-02","slot_mulai":3,"slot_selesai":5,"jumlah":1}`, []int{1, 2})
	cek(`{"sumber_daya_id":1,"peminjam":"X","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":2,"jumlah":1}`, []int{1})

	// 201 tidak boleh memuat bentrok_dengan
	rec := kirim(h, "POST", "/peminjaman", `{"sumber_daya_id":1,"peminjam":"X","tanggal":"2026-11-02","slot_mulai":7,"slot_selesai":8,"jumlah":1}`)
	if rec.Code != http.StatusCreated || strings.Contains(rec.Body.String(), "bentrok_dengan") {
		t.Fatalf("201 tidak boleh memuat bentrok_dengan: %d %s", rec.Code, rec.Body.String())
	}
	// 400 juga tidak boleh
	rec = kirim(h, "POST", "/peminjaman", `{"sumber_daya_id":99,"peminjam":"X","tanggal":"2026-11-02","slot_mulai":1,"slot_selesai":1,"jumlah":1}`)
	if rec.Code != http.StatusBadRequest || strings.Contains(rec.Body.String(), "bentrok_dengan") {
		t.Fatalf("400 tidak boleh memuat bentrok_dengan: %d %s", rec.Code, rec.Body.String())
	}
}
