package main

import (
	"net/http"
	"stajtakip/database"
	"stajtakip/routes"
)

type StajServer struct {
	addr   string
	server *http.Server
	mux    *http.ServeMux
}

func (s *StajServer) SetHandlers(conn *database.Connection) {
	// TODO giriş yapmadan erişilemesin?

	s.mux.Handle("/assets/", routes.Assets)
	s.mux.Handle("/uploads/", routes.Uploads)

	// Raporlar
	s.mux.Handle("/", routes.Raporlar{conn})

	// Öğrenci
	s.mux.Handle("/ogrenci-ekle", routes.OgrenciEkle{conn})
	s.mux.Handle("/ogrenci-listele", routes.OgrenciListesi{conn}) // Stajı bitenler
	s.mux.Handle("/ogrenci-listele-pdf", routes.OgrenciListesiPDF{conn})
	s.mux.Handle("/ogrenci-listele-xls", routes.OgrenciListesiXLS{conn})
	s.mux.Handle("/ogrenci-belge-ekle", routes.OgrenciBelge{conn}) // DGS için pdf
	s.mux.Handle("/ogrenci-bilgi", routes.OgrenciBilgi{conn})        // OgrNo -> Bilgiler, stajlar, staj bitim durumu

	// Staj
	s.mux.Handle("/staj-ekle", routes.StajEkle{conn})
	s.mux.Handle("/konular", routes.KonuListesi{conn}) // Konu listesi

	// Mülakat
	s.mux.Handle("/mulakat", routes.MulakatListesi{conn}) // Tarih/saat, komisyon -> öğrenci, ...
	s.mux.Handle("/mulakat-pdf", routes.MulakatListesiPDF{conn})
	s.mux.Handle("/mulakat-xlsx", routes.MulakatListesiXLSX{conn})
	s.mux.Handle("/sonuc-listele", routes.SonucListele{conn})   // Mülakat sonuçları
	s.mux.Handle("/sonuc-listele-pdf", routes.SonucListesiPDF{conn})
	s.mux.Handle("/sonuc-listele-xlsx", routes.SonucListesiXLSX{conn})
	s.mux.Handle("/staj-degerlendir", routes.StajDegerlendir{conn})    // Mülakat sonucu ekle
	s.mux.Handle("/komisyon", routes.KomisyonListesi{conn}) // Komisyon Listesi
}

func (s *StajServer) Run() error {
	return s.server.ListenAndServe()
}

func NewStajServer(addr string) *StajServer {
	mux := http.NewServeMux()
	server := &http.Server{Addr: addr, Handler: mux}

	return &StajServer{
		addr:   addr,
		server: server,
		mux:    mux,
	}
}
