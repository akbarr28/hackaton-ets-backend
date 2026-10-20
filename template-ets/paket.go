package templateets

import "net/http"

type paketItemInput struct {
	SumberDayaID *int `json:"sumber_daya_id"`
	Jumlah       *int `json:"jumlah"`
}

type paketInput struct {
	Peminjam    *string          `json:"peminjam"`
	Tanggal     *string          `json:"tanggal"`
	SlotMulai   *int             `json:"slot_mulai"`
	SlotSelesai *int             `json:"slot_selesai"`
	Item        []paketItemInput `json:"item"`
}

func (s *Store) handleBuatPaket(w http.ResponseWriter, r *http.Request) {
	var in paketInput
	if !decodeBody(r, &in) {
		writeError(w, http.StatusBadRequest, "body harus berupa JSON yang valid")
		return
	}

	// TAHAP 1: validasi semuanya (400). Belum ada yang disimpan.
	peminjam, pesan := validasiTeks(in.Peminjam, "peminjam", 50)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	tanggal, pesan := validasiTanggal(in.Tanggal)
	if pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	if pesan = validasiSlot(in.SlotMulai, in.SlotSelesai); pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	if len(in.Item) < 1 || len(in.Item) > 10 {
		writeError(w, http.StatusBadRequest, "item harus berisi 1 sampai 10 elemen")
		return
	}

	dipakai := map[int]bool{}
	calon := make([]Peminjaman, 0, len(in.Item))
	for _, it := range in.Item {
		if it.SumberDayaID == nil || it.Jumlah == nil {
			writeError(w, http.StatusBadRequest, "setiap item wajib punya sumber_daya_id dan jumlah")
			return
		}
		sd, ada := s.cariSumberDaya(*it.SumberDayaID)
		if !ada {
			writeError(w, http.StatusBadRequest, "sumber daya tidak ditemukan")
			return
		}
		if *it.Jumlah < 1 || *it.Jumlah > sd.Stok {
			writeError(w, http.StatusBadRequest, "jumlah harus antara 1 dan stok sumber daya")
			return
		}
		if dipakai[sd.ID] {
			writeError(w, http.StatusBadRequest, "satu sumber daya tidak boleh muncul dua kali dalam satu paket")
			return
		}
		dipakai[sd.ID] = true
		calon = append(calon, Peminjaman{
			SumberDayaID: sd.ID,
			Peminjam:     peminjam,
			Tanggal:      tanggal,
			SlotMulai:    *in.SlotMulai,
			SlotSelesai:  *in.SlotSelesai,
			Jumlah:       *it.Jumlah,
		})
	}

	// TAHAP 2: cek bentrok semua item (409). Masih belum ada yang disimpan.
	for _, c := range calon {
		if s.adaBentrok(c, 0) {
			writeError(w, http.StatusConflict, "salah satu item bentrok, paket dibatalkan seluruhnya")
			return
		}
	}

	// TAHAP 3: baru simpan. Id paket dan id peminjaman dibagikan HANYA di sini, jadi paket yang gagal tidak menghabiskan nomor id.
	paketID := s.idPaketBerikut
	s.idPaketBerikut++
	s.paketAda[paketID] = true

	hasil := make([]Peminjaman, 0, len(calon))
	for _, c := range calon {
		c.PaketID = paketID
		hasil = append(hasil, s.tambahPeminjaman(c))
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"paket_id":   paketID,
		"peminjaman": hasil,
	})
}

func (s *Store) handleHapusPaket(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return
	}
	if !s.paketAda[id] {
		writeError(w, http.StatusNotFound, "paket tidak ditemukan")
		return
	}
	sisa := make([]Peminjaman, 0, len(s.peminjaman))
	for _, p := range s.peminjaman {
		if p.PaketID != id {
			sisa = append(sisa, p)
		}
	}
	s.peminjaman = sisa
	delete(s.paketAda, id)
	w.WriteHeader(http.StatusNoContent)
}
