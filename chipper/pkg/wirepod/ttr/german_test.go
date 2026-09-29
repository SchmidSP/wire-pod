package wirepod_ttr

import (
	"testing"

	"github.com/kercre123/wire-pod/chipper/pkg/vars"
)

func TestParseGermanCardinal(t *testing.T) {
	cases := map[string]int{
		"fünf":             5,
		"eins":             1,
		"eine":             1,
		"fünfundzwanzig":   25,
		"fünf und zwanzig": 25,
		"einundzwanzig":    21,
		"dreißig":          30,
		"hundert":          100,
	}
	for in, want := range cases {
		got, ok := parseGermanCardinal(in)
		if !ok || got != want {
			t.Fatalf("%q: got %d ok=%v, want %d", in, got, ok, want)
		}
	}
}

func TestGermanTimerPhrase(t *testing.T) {
	vars.APIConfig.STT.Language = "de-DE"
	if got := words2num("stelle einen timer für fünf minuten"); got != "300" {
		t.Fatalf("fünf minuten: got %s", got)
	}
	if got := words2num("timer für 5 minuten und 10 sekunden"); got != "310" {
		t.Fatalf("5 minuten 10 sekunden: got %s", got)
	}
	if got := words2num("stell einen timer für fünfundzwanzig sekunden"); got != "25" {
		t.Fatalf("fünfundzwanzig sekunden: got %s", got)
	}
}

func TestWavRoundTripHeader(t *testing.T) {
	pcm := []byte{0x00, 0x00, 0xff, 0x7f, 0x00, 0x80, 0x10, 0x00}
	wav := makeWav(22050, pcm)
	got, rate, err := wavToMono16(wav)
	if err != nil {
		t.Fatal(err)
	}
	if rate != 22050 {
		t.Fatalf("rate %d", rate)
	}
	if string(got) != string(pcm) {
		t.Fatalf("pcm mismatch")
	}
	out := resamplePCM16(pcm, 22050, 22050)
	if len(out) != len(pcm) {
		t.Fatalf("resample identity len %d", len(out))
	}
}

func makeWav(rate int, pcm []byte) []byte {
	// minimal PCM WAV
	dataSize := len(pcm)
	riffSize := 36 + dataSize
	buf := make([]byte, 44+dataSize)
	copy(buf[0:], []byte("RIFF"))
	put32(buf[4:], riffSize)
	copy(buf[8:], []byte("WAVE"))
	copy(buf[12:], []byte("fmt "))
	put32(buf[16:], 16)
	buf[20] = 1
	buf[22] = 1
	put32(buf[24:], rate)
	put32(buf[28:], rate*2)
	buf[32] = 2
	buf[34] = 16
	copy(buf[36:], []byte("data"))
	put32(buf[40:], dataSize)
	copy(buf[44:], pcm)
	return buf
}

func put32(b []byte, v int) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}
