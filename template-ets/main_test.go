package templateets

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func kirim(h http.Handler, metode, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metode, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestBersamaan(t *testing.T) {
	h := NewHandler()
	kirim(h, "POST", "/sumber-daya", `{"nama":"Aula A","jenis":"ruang","stok":1}`)
	body := `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":3,"slot_selesai":4,"jumlah":1}`

	var wg sync.WaitGroup
	var mu sync.Mutex
	sukses := 0
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if kirim(h, "POST", "/peminjaman", body).Code == http.StatusCreated {
				mu.Lock()
				sukses++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if sukses != 1 {
		t.Fatalf("seharusnya tepat 1 yang berhasil, dapat %d", sukses)
	}
}

func TestPaketAtomik(t *testing.T) {
	h := NewHandler()
	kirim(h, "POST", "/sumber-daya", `{"nama":"Aula A","jenis":"ruang","stok":1}`)
	kirim(h, "POST", "/sumber-daya", `{"nama":"Proyektor","jenis":"alat","stok":3}`)
	kirim(h, "POST", "/peminjaman", `{"sumber_daya_id":2,"peminjam":"X","tanggal":"2026-11-02","slot_mulai":2,"slot_selesai":4,"jumlah":3}`)

	paket := `{"peminjam":"Panitia","tanggal":"2026-11-02","slot_mulai":2,"slot_selesai":4,"item":[{"sumber_daya_id":1,"jumlah":1},{"sumber_daya_id":2,"jumlah":1}]}`
	if rec := kirim(h, "POST", "/paket-peminjaman", paket); rec.Code != http.StatusConflict {
		t.Fatalf("harus 409, dapat %d", rec.Code)
	}
	rec := kirim(h, "GET", "/peminjaman?sumber_daya_id=1", "")
	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("aula harus kosong, dapat %s", rec.Body.String())
	}
	ok := `{"peminjam":"Panitia","tanggal":"2026-11-02","slot_mulai":2,"slot_selesai":4,"item":[{"sumber_daya_id":1,"jumlah":1}]}`
	rec = kirim(h, "POST", "/paket-peminjaman", ok)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"paket_id":1`) {
		t.Fatalf("paket pertama yang berhasil harus punya paket_id 1: %d %s", rec.Code, rec.Body.String())
	}
}
