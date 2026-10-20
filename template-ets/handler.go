package templateets

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, pesan string) {
	writeJSON(w, status, map[string]string{"error": pesan})
}

func decodeBody(r *http.Request, dst any) bool {
	return json.NewDecoder(r.Body).Decode(dst) == nil
}

func parseID(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	return id, err == nil
}

// ---------- B1: sumber daya ----------

func (s *Store) handleBuatSumberDaya(w http.ResponseWriter, r *http.Request) {
	var in sumberDayaInput
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "body harus berupa JSON yang valid")
		return
	}
	sd, pesan := validasiSumberDaya(in)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	writeJSON(w, http.StatusCreated, s.tambahSumberDaya(sd))
}

func (s *Store) handleDaftarSumberDaya(w http.ResponseWriter, r *http.Request) {
	hasil := make([]SumberDaya, len(s.sumberDaya)) // bukan nil, jadi tampil sebagai []
	copy(hasil, s.sumberDaya)
	writeJSON(w, http.StatusOK, hasil)
}

func (s *Store) handleDetailSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	sd, ada := s.cariSumberDaya(id)
	if !ada {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, sd)
}

// ---------- B2 & B3: peminjaman ----------

func (s *Store) handleBuatPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanInput
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "body harus berupa JSON yang valid")
		return
	}
	p, pesan := s.validasiPeminjaman(in)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	if s.adaBentrok(p, 0) {
		writeError(w, http.StatusConflict, "kapasitas sumber daya pada slot tersebut terlampaui")
		return
	}
	writeJSON(w, http.StatusCreated, s.tambahPeminjaman(p))
}

func (s *Store) handleDaftarPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	tanggal := strings.TrimSpace(q.Get("tanggal"))
	status := strings.TrimSpace(q.Get("status"))

	filterSD, adaFilterSD := 0, false
	if v := strings.TrimSpace(q.Get("sumber_daya_id")); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus berupa angka")
			return
		}
		filterSD, adaFilterSD = id, true
	}

	hasil := []Peminjaman{}
	for _, p := range s.peminjaman {
		if adaFilterSD && p.SumberDayaID != filterSD {
			continue
		}
		if tanggal != "" && p.Tanggal != tanggal {
			continue
		}
		if status != "" && p.Status != status {
			continue
		}
		hasil = append(hasil, p)
	}
	writeJSON(w, http.StatusOK, hasil)
}

func (s *Store) handleDetailPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	idx := s.indeksPeminjaman(id)
	if idx < 0 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, s.peminjaman[idx])
}

func (s *Store) handlePerbaruiPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	idx := s.indeksPeminjaman(id)
	if idx < 0 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	var in peminjamanInput
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "body harus berupa JSON yang valid")
		return
	}
	p, pesan := s.validasiPeminjaman(in)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	if s.adaBentrok(p, id) {
		writeError(w, http.StatusConflict, "kapasitas sumber daya pada slot tersebut terlampaui")
		return
	}

	p.ID = id
	p.Status = s.peminjaman[idx].Status
	p.PaketID = s.peminjaman[idx].PaketID
	s.peminjaman[idx] = p
	writeJSON(w, http.StatusOK, p)
}

func (s *Store) handleHapusPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	idx := s.indeksPeminjaman(id)
	if idx < 0 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	s.peminjaman = append(s.peminjaman[:idx], s.peminjaman[idx+1:]...)
	w.WriteHeader(http.StatusNoContent)
}
