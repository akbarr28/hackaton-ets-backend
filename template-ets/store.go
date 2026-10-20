package templateets

import (
	"net/http"
	"sync"
)

type Store struct {
	mu              sync.Mutex
	sumberDaya      []SumberDaya
	peminjaman      []Peminjaman
	paketAda        map[int]bool
	idSDBerikut     int
	idPinjamBerikut int
	idPaketBerikut  int
}

func NewStore() *Store {
	return &Store{
		paketAda:        map[int]bool{},
		idSDBerikut:     1,
		idPinjamBerikut: 1,
		idPaketBerikut:  1,
	}
}

func (s *Store) cariSumberDaya(id int) (SumberDaya, bool) {
	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			return sd, true
		}
	}
	return SumberDaya{}, false
}

func (s *Store) tambahSumberDaya(sd SumberDaya) SumberDaya {
	sd.ID = s.idSDBerikut
	s.idSDBerikut++
	s.sumberDaya = append(s.sumberDaya, sd)
	return sd
}

func (s *Store) indeksPeminjaman(id int) int {
	for i, p := range s.peminjaman {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func (s *Store) tambahPeminjaman(p Peminjaman) Peminjaman {
	p.ID = s.idPinjamBerikut
	s.idPinjamBerikut++
	p.Status = "aktif"
	s.peminjaman = append(s.peminjaman, p)
	return p
}

func (s *Store) unitTerpakai(sdID int, tanggal string, slot int, kecualiID int) int {
	total := 0
	for _, p := range s.peminjaman {
		if p.Status != "aktif" || p.ID == kecualiID {
			continue
		}
		if p.SumberDayaID != sdID || p.Tanggal != tanggal {
			continue
		}
		if slot >= p.SlotMulai && slot <= p.SlotSelesai {
			total += p.Jumlah
		}
	}
	return total
}

func (s *Store) adaBentrok(p Peminjaman, kecualiID int) bool {
	sd, ada := s.cariSumberDaya(p.SumberDayaID)
	if !ada {
		return false
	}
	for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
		if s.unitTerpakai(p.SumberDayaID, p.Tanggal, slot, kecualiID)+p.Jumlah > sd.Stok {
			return true
		}
	}
	return false
}

func (s *Store) kunci(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		h(w, r)
	}
}
