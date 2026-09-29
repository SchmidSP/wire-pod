# wire-pod Deutsch

Fork von [kercre123/wire-pod](https://github.com/kercre123/wire-pod) für den Anki Vector.

Vector versteht Deutsch und spricht KI-Antworten lokal mit **Piper / Thorsten**. Als Wissensquelle ist **GPT-6 Luna** (`gpt-6-luna`) voreingestellt.

Die Idee stammt von [Kolle1979/wire-pod-german](https://github.com/Kolle1979/wire-pod-german). Dort liegen die deutschen Module nur neben dem Code. Hier sind sie im echten Sprachpfad von wire-pod.

Ausführlicher steht dasselbe in [DEUTSCH.md](DEUTSCH.md). Das englische Original-Wiki gilt weiter für Installation, Zertifikate und Vector-Setup: [Installation](https://github.com/kercre123/wire-pod/wiki/Installation).

## Was dieser Fork ändert

- Spracherkennung **German (DE)** über das Vosk-Modell `de-DE`.
- Timer verstehen deutsche Zahlen, auch „fünfundzwanzig Sekunden“.
- Knowledge Graph und Intent Graph antworten auf Deutsch.
- Diese Antworten und Lua-`sayText` spricht **Thorsten**, offline. Die Audiodatei geht als 16-kHz-PCM an Vector.
- OpenAI-Modell standardmäßig **gpt-6-luna**, zusätzlich gpt-6-sol und die älteren Modelle. Für GPT-6 gilt `reasoning_effort=low` und `verbosity=low`. Fällt das Modell aus, wird `gpt-4o-mini` benutzt.

Feste Ansagen der Vector-Firmware (viele eingebaute Befehle) bleiben die Originalstimme. Die kann kein Deutsch.

## Windows

Die fertige Installation von kercre123 enthält diesen Code nicht. Dieser Fork muss selbst gebaut werden.

```powershell
git clone https://github.com/SchmidSP/wire-pod.git
cd wire-pod
powershell -ExecutionPolicy Bypass -File .\install-german.ps1
```

Das Skript lädt Piper und die Stimme Thorsten nach `%USERPROFILE%\wire-pod-data\` und schreibt `german.env`. Diese Datei liest wire-pod beim Start. Danach chipper wie im [Wiki](https://github.com/kercre123/wire-pod/wiki/Installation) bauen und starten.

## Linux und Raspberry Pi

```bash
git clone https://github.com/SchmidSP/wire-pod.git
cd wire-pod
bash setup-german.sh
sudo STT=vosk ./setup.sh
```

Piper und die Stimme landen in `~/wire-pod-data/`. Das Setup-Skript braucht kein Root.

## Im Webinterface

1. Sprache **German (DE)**. Dabei wird das deutsche Vosk-Modell geladen.
2. Knowledge Graph: Anbieter **OpenAI**, Modell **gpt-6-luna**, eigenen Key eintragen.
3. **Intent-Graph** einschalten, damit freie Sätze wie „Erzähl mir einen Witz“ an die KI gehen.
4. Vector neu verbinden.

## Sätze zum Ausprobieren

```text
Hey Vector, wie ist das Wetter
Hey Vector, wie spät ist es
Hey Vector, stell einen Timer für fünf Minuten
Hey Vector, ich habe eine Frage
Hey Vector, erzähl mir einen Witz
```

Der letzte Satz ist kein fester Befehl. Dafür muss der Intent-Graph an sein.

## Variablen

| Variable | Bedeutung |
| --- | --- |
| `PIPER_BIN` | Pfad zu `piper` oder `piper.exe` |
| `PIPER_MODEL` | Pfad zu `de_DE-thorsten-medium.onnx` |
| `TTS_SERVICE=piper` | Piper auch dann nutzen, wenn die STT-Sprache nicht de-DE ist |
| `KNOWLEDGE_MODEL` | wird übernommen, wenn im Knowledge Graph noch kein Modell steht |

## Credits

- [Digital Dream Labs](https://github.com/digital-dream-labs) für chipper und Escape Pod
- [kercre123/wire-pod](https://github.com/kercre123/wire-pod) und die genannten Mitwirkenden: bliteknight, dietb, fforchino, xanathon
- [Kolle1979/wire-pod-german](https://github.com/Kolle1979/wire-pod-german) für die deutsche Piper-Idee
- Piper und die Stimme Thorsten: Rhasspy / Thorsten Müller
- Vosk: Alpha Cephei
