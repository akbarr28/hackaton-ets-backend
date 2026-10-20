package templateets

import (
	"log"
	"net/http"
)

func NewHandler() http.Handler {
	s := NewStore() // setiap NewHandler() punya penyimpanan baru yang kosong
	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", s.handleBuatSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.handleDaftarSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.handleDetailSumberDaya)

	mux.HandleFunc("POST /peminjaman", s.handleBuatPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.handleDaftarPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.handleDetailPeminjaman)

	// Supaya path yang tidak dikenal tetap membalas JSON, bukan teks biasa.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "endpoint tidak ditemukan")
	})
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", NewHandler()))
}