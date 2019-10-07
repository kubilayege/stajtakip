package cfg

import (
	"bufio"
	"encoding/json"
	"os"
)

type Ayarlar struct {
	SunucuAdresi     string
	VeritabaniAdresi string
	LogDosyasi       string
	GerekenStajGunu  int
	MaxIstekBoyutu   int64
}

// Ayarlara varsayılan değerler ata
var ayarlar Ayarlar = Ayarlar{
	SunucuAdresi:     "127.0.0.1:8080",
	VeritabaniAdresi: "root@localhost/stajtakip",
	LogDosyasi:       "stajtakip.log",
	GerekenStajGunu:  57,
	MaxIstekBoyutu:   32000000,
}

func AyarlariOku(dosya string) error {
	// Ayarlar dosyasını aç
	f, err := os.Open(dosya)
	if err != nil {
		return err
	}
	defer f.Close()

	// Dosyayı tamponlu okuyucu ile aç
	r := bufio.NewReader(f)

	// Okuyucuyu JSON çözücüye ver
	dec := json.NewDecoder(r)

	// JSON verilerini struct'a aktar
	err = dec.Decode(&ayarlar)
	if err != nil {
		return err
	}

	return nil
}

func SunucuAdresi() string {
	return ayarlar.SunucuAdresi
}

func VeritabaniAdresi() string {
	return ayarlar.VeritabaniAdresi
}

func LogDosyasi() string {
	return ayarlar.LogDosyasi
}

func GerekenStajGunu() int {
	return ayarlar.GerekenStajGunu
}

func MaxIstekBoyutu() int64 {
	return ayarlar.MaxIstekBoyutu
}
