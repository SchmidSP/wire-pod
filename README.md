# wire-pod Deutsch

Fork von [kercre123/wire-pod](https://github.com/kercre123/wire-pod) für den Anki Vector.

Vector versteht Deutsch und spricht KI-Antworten lokal mit **Piper / Thorsten**. **GPT-6 Luna** läuft über **Azure AI Foundry** (Provider Custom), nicht über api.openai.com.

Die Idee stammt von [Kolle1979/wire-pod-german](https://github.com/Kolle1979/wire-pod-german). Dort liegen die deutschen Module nur neben dem Code. Hier sind sie im echten Sprachpfad von wire-pod.

Ausführlicher steht dasselbe in [DEUTSCH.md](DEUTSCH.md). Das englische Original-Wiki gilt weiter für Installation, Zertifikate und Vector-Setup: [Installation](https://github.com/kercre123/wire-pod/wiki/Installation).

## Was dieser Fork ändert

- Spracherkennung **German (DE)** über das Vosk-Modell `de-DE`.
- Timer verstehen deutsche Zahlen, auch „fünfundzwanzig Sekunden“.
- Knowledge Graph und Intent Graph antworten auf Deutsch.
- Diese Antworten und Lua-`sayText` spricht **Thorsten**, offline. Die Audiodatei geht als 16-kHz-PCM an Vector.
- OpenAI-kompatibler **Custom**-Provider für Azure Foundry. Für `gpt-6-luna` gehen nur `max_completion_tokens` raus, kein `max_tokens`, keine Sampling-Parameter. Leere Azure-Stream-Chunks (`choices: []`) werden übersprungen. Ein Responses-Link wird auf `/openai/v1` gekürzt. Der Python-Proxy ist damit nicht mehr nötig.

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
2. Knowledge Graph: Anbieter **Custom**, Endpoint `https://ksgptsweden.cognitiveservices.azure.com/openai/v1`, Modell **gpt-6-luna**, Azure-Key aus Foundry. Nicht den Link `/openai/responses?api-version=...` eintragen.
3. **Intent-Graph** an, **LLM-Actions** aus (sonst schickt das Modell `{{playAnimation…}}` und bricht die Ansage ab).
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
| `KNOWLEDGE_PROVIDER` | `custom` für Azure Foundry |
| `KNOWLEDGE_ENDPOINT` | `https://<resource>.cognitiveservices.azure.com/openai/v1` |
| `KNOWLEDGE_MODEL` | Deployment-Name, hier `gpt-6-luna` |

## Credits

- [Digital Dream Labs](https://github.com/digital-dream-labs) für chipper und Escape Pod
- [kercre123/wire-pod](https://github.com/kercre123/wire-pod) und die genannten Mitwirkenden: bliteknight, dietb, fforchino, xanathon
- [Kolle1979/wire-pod-german](https://github.com/Kolle1979/wire-pod-german) für die deutsche Piper-Idee
- Piper und die Stimme Thorsten: Rhasspy / Thorsten Müller
- Vosk: Alpha Cephei
