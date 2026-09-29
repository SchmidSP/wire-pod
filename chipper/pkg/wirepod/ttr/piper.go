package wirepod_ttr

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fforchino/vector-go-sdk/pkg/vector"
	"github.com/fforchino/vector-go-sdk/pkg/vectorpb"
	"github.com/kercre123/wire-pod/chipper/pkg/logger"
	"github.com/kercre123/wire-pod/chipper/pkg/vars"
)

// Piper (Thorsten) speaks German locally. Vector's onboard voice is English-only,
// which is why German answers used to sound wrong. Audio is streamed to the robot
// as 16 kHz mono PCM, the same path OpenAI TTS already uses.

func piperShouldSpeak() bool {
	lang := vars.APIConfig.STT.Language
	forced := os.Getenv("TTS_SERVICE") == "piper" || strings.TrimSpace(os.Getenv("PIPER_MODEL")) != ""
	if lang != "de-DE" && !forced {
		return false
	}
	return piperExecutable() != "" && piperModelPath() != ""
}

func piperExecutable() string {
	if p := strings.TrimSpace(os.Getenv("PIPER_BIN")); p != "" {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	if p, err := exec.LookPath("piper"); err == nil {
		return p
	}
	if p, err := exec.LookPath("piper.exe"); err == nil {
		return p
	}
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, "wire-pod-data", "piper", "piper.exe"),
		filepath.Join(home, "wire-pod-data", "piper", "piper"),
		"/usr/local/bin/piper",
		"/usr/bin/piper",
		filepath.Join("piper", "piper.exe"),
		filepath.Join("piper", "piper"),
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c
		}
	}
	return ""
}

func piperModelPath() string {
	if p := strings.TrimSpace(os.Getenv("PIPER_MODEL")); p != "" {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	home, _ := os.UserHomeDir()
	names := []string{
		"de_DE-thorsten-medium.onnx",
		"de_DE-thorsten-high.onnx",
		"de_DE-karlsson-low.onnx",
	}
	dirs := []string{
		filepath.Join(home, "wire-pod-data", "voices"),
		"/opt/wire-pod/voices",
		"voices",
		filepath.Join("..", "voices"),
	}
	for _, d := range dirs {
		for _, n := range names {
			p := filepath.Join(d, n)
			if st, err := os.Stat(p); err == nil && !st.IsDir() {
				return p
			}
		}
	}
	return ""
}

func cleanForGermanTTS(text string) string {
	repl := []struct{ from, to string }{
		{"WLAN", "W-LAN"},
		{"WiFi", "Wei-Fai"},
		{"WIFI", "Wei-Fai"},
		{"OK", "O-K"},
		{"API", "A-P-I"},
		{"URL", "U-R-L"},
		{"USB", "U-S-B"},
		{"GPS", "G-P-S"},
		{"SMS", "S-M-S"},
		{"PDF", "P-D-F"},
	}
	for _, r := range repl {
		text = strings.ReplaceAll(text, r.from, r.to)
	}
	text = strings.ReplaceAll(text, "&", " und ")
	text = strings.ReplaceAll(text, "*", "")
	text = strings.ReplaceAll(text, "#", "")
	text = strings.ReplaceAll(text, "@", "")
	text = strings.ReplaceAll(text, "^", "")
	return strings.TrimSpace(text)
}

func DoSayText_Piper(robot *vector.Vector, input string) error {
	text := cleanForGermanTTS(input)
	if text == "" {
		return nil
	}
	bin := piperExecutable()
	model := piperModelPath()
	if bin == "" || model == "" {
		return fmt.Errorf("piper ist nicht installiert")
	}
	tmp, err := os.CreateTemp("", "wirepod-piper-*.wav")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpName)

	cmd := exec.Command(bin, "--model", model, "--output_file", tmpName)
	cmd.Dir = filepath.Dir(bin)
	cmd.Stdin = strings.NewReader(text)
	cmd.Env = append(os.Environ(),
		"LD_LIBRARY_PATH="+filepath.Dir(bin),
		"ESPEAK_DATA_PATH="+filepath.Join(filepath.Dir(bin), "espeak-ng-data"),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("piper: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	wav, err := os.ReadFile(tmpName)
	if err != nil {
		return err
	}
	pcm, rate, err := wavToMono16(wav)
	if err != nil {
		return err
	}
	pcm16 := resamplePCM16(pcm, rate, 16000)
	pcm16 = scalePCM16(pcm16, 1.6)
	logger.Println(fmt.Sprintf("Piper TTS Thorsten (%d bytes, %d Hz -> 16 kHz)", len(pcm16), rate))
	return streamAudioChunks(robot, chunkPCM(pcm16))
}

func streamAudioChunks(robot *vector.Vector, audioChunks [][]byte) error {
	if robot == nil || robot.Conn == nil {
		return fmt.Errorf("kein Roboter verbunden")
	}
	vclient, err := robot.Conn.ExternalAudioStreamPlayback(context.Background())
	if err != nil {
		return err
	}
	if err := vclient.Send(&vectorpb.ExternalAudioStreamRequest{
		AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamPrepare{
			AudioStreamPrepare: &vectorpb.ExternalAudioStreamPrepare{
				AudioFrameRate: 16000,
				AudioVolume:    100,
			},
		},
	}); err != nil {
		return err
	}
	var all []byte
	for _, chunk := range audioChunks {
		all = append(all, chunk...)
	}
	go func() {
		for _, chunk := range audioChunks {
			_ = vclient.Send(&vectorpb.ExternalAudioStreamRequest{
				AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamChunk{
					AudioStreamChunk: &vectorpb.ExternalAudioStreamChunk{
						AudioChunkSizeBytes: 1024,
						AudioChunkSamples:   chunk,
					},
				},
			})
			time.Sleep(time.Millisecond * 25)
		}
		_ = vclient.Send(&vectorpb.ExternalAudioStreamRequest{
			AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamComplete{
				AudioStreamComplete: &vectorpb.ExternalAudioStreamComplete{},
			},
		})
	}()
	time.Sleep(pcmLength(all) + (time.Millisecond * 80))
	return nil
}

func chunkPCM(pcm []byte) [][]byte {
	var audioChunks [][]byte
	for len(pcm) > 0 {
		if len(pcm) < 1024 {
			chunk := make([]byte, 1024)
			copy(chunk, pcm)
			audioChunks = append(audioChunks, chunk)
			break
		}
		chunk := make([]byte, 1024)
		copy(chunk, pcm[:1024])
		audioChunks = append(audioChunks, chunk)
		pcm = pcm[1024:]
	}
	return audioChunks
}

func wavToMono16(data []byte) ([]byte, int, error) {
	if len(data) < 44 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return nil, 0, fmt.Errorf("keine WAV-Datei")
	}
	var rate, channels, bits int
	var pcm []byte
	off := 12
	for off+8 <= len(data) {
		id := string(data[off : off+4])
		size := int(binary.LittleEndian.Uint32(data[off+4 : off+8]))
		off += 8
		if off+size > len(data) {
			return nil, 0, fmt.Errorf("kaputte WAV-Datei")
		}
		chunk := data[off : off+size]
		switch id {
		case "fmt ":
			if len(chunk) < 16 {
				return nil, 0, fmt.Errorf("fmt-Chunk zu kurz")
			}
			format := binary.LittleEndian.Uint16(chunk[0:2])
			channels = int(binary.LittleEndian.Uint16(chunk[2:4]))
			rate = int(binary.LittleEndian.Uint32(chunk[4:8]))
			bits = int(binary.LittleEndian.Uint16(chunk[14:16]))
			if format != 1 || bits != 16 {
				return nil, 0, fmt.Errorf("piper-wav ist nicht 16-bit PCM (format=%d bits=%d)", format, bits)
			}
		case "data":
			pcm = append([]byte(nil), chunk...)
		}
		off += size
		if size%2 == 1 {
			off++
		}
	}
	if rate == 0 || len(pcm) < 2 {
		return nil, 0, fmt.Errorf("wav ohne audiodaten")
	}
	if channels > 1 {
		mono := make([]byte, 0, len(pcm)/channels)
		frame := channels * 2
		for i := 0; i+frame <= len(pcm); i += frame {
			mono = append(mono, pcm[i], pcm[i+1])
		}
		pcm = mono
	}
	return pcm, rate, nil
}

func resamplePCM16(pcm []byte, fromRate, toRate int) []byte {
	if fromRate <= 0 || toRate <= 0 || fromRate == toRate || len(pcm) < 2 {
		return pcm
	}
	in := bytesToInt16s(pcm)
	outLen := int(math.Round(float64(len(in)) * float64(toRate) / float64(fromRate)))
	if outLen < 1 {
		outLen = 1
	}
	out := make([]int16, outLen)
	ratio := float64(fromRate) / float64(toRate)
	last := len(in) - 1
	for i := 0; i < outLen; i++ {
		pos := float64(i) * ratio
		idx := int(pos)
		if idx >= last {
			out[i] = in[last]
			continue
		}
		frac := pos - float64(idx)
		sample := float64(in[idx])*(1-frac) + float64(in[idx+1])*frac
		out[i] = int16(sample)
	}
	return int16sToBytes(out)
}

func scalePCM16(pcm []byte, factor float64) []byte {
	if factor == 1 || len(pcm) < 2 {
		return pcm
	}
	samples := bytesToInt16s(pcm)
	for i := range samples {
		scaled := float64(samples[i]) * factor
		if scaled > math.MaxInt16 {
			samples[i] = math.MaxInt16
		} else if scaled < math.MinInt16 {
			samples[i] = math.MinInt16
		} else {
			samples[i] = int16(scaled)
		}
	}
	return int16sToBytes(samples)
}
