package templateets

import "sync"

type Store struct {
	mu              sync.Mutex
	sumberDaya      []SumberDaya
	peminjaman      []Peminjaman
	idSDBerikut     int
	idPinjamBerikut int
}

func NewStore() *Store {
	return &Store{idSDBerikut: 1, idPinjamBerikut: 1}
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