# Wire-pod auf Deutsch

Dieser Fork von [kercre123/wire-pod](https://github.com/kercre123/wire-pod) spricht und versteht Deutsch. Die Idee entspricht [wire-pod-german](https://github.com/Kolle1979/wire-pod-german) (Piper, Stimme Thorsten, Vosk `de-DE`), ist aber in den echten Sprachpfad von wire-pod eingebaut und nicht nur als lose Dateien abgelegt.

Zusätzlich spricht der Knowledge Graph über **Azure AI Foundry**. Deployment **gpt-6-luna** auf `https://ksgptsweden.cognitiveservices.azure.com/openai/v1`, Provider **Custom**. Der alte Responses-Link (`/openai/responses?api-version=2025-04-01-preview`) wird, falls doch eingefügt, auf `/openai/v1` gekürzt. Auth bleibt `Authorization: Bearer`.

GPT-6 Luna lehnt `max_tokens`, `temperature`, `top_p` und die Penalty-Felder ab. Dieser Fork schickt sie nicht. Azure beginnt den Stream mit `choices: []` (Content-Filter). Das wird übersprungen, der Text danach geht an Piper. Action-Tags `{{…}}` werden aus der gesprochenen Antwort entfernt, solange LLM-Actions aus sind. Ein extra Python-Proxy ist nicht nötig.

## Was Vector auf Deutsch kann

- Spracherkennung über das vorhandene Vosk-Modell `de-DE` (Wetter, Timer, Uhrzeit, Tanzen, Namen, …).
- Timer verstehen jetzt deutsche Zahlen: „Stell einen Timer für fünf Minuten“, auch „fünfundzwanzig Sekunden“.
- KI-Antworten (Knowledge Graph und Intent Graph) kommen auf Deutsch zurück.
- Diese Antworten spricht **Piper mit der Stimme Thorsten**, lokal, ohne OpenAI-TTS. Die Audiodatei geht als 16-kHz-PCM an Vector.
- Lua-`sayText` nutzt denselben Weg.
- Feste Ansagen, die die Vector-Firmware selbst abspielt (viele eingebaute Befehle), bleiben die Originalstimme. Die kennt kein Deutsch.

## Windows

Im Ordner des geklonten Repos, in PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File .\install-german.ps1
```

Das lädt Piper und Thorsten nach `%USERPROFILE%\wire-pod-data\` und schreibt `german.env`. Danach wire-pod **aus diesem Fork neu bauen** und starten (die fertige Original-Installation von kercre123 enthält diesen Code nicht).

Im Webinterface (Port 8080):

1. STT-Sprache **German (DE)**. Das deutsche Vosk-Modell wird dabei heruntergeladen.
2. Knowledge Graph: Anbieter **Custom**, Endpoint `https://ksgptsweden.cognitiveservices.azure.com/openai/v1`, Modell **gpt-6-luna**, Azure-Key. **Intent-Graph** an, **LLM-Actions** aus.
3. **Intent-Graph** einschalten, wenn Vector auch freie Sätze wie „Erzähl mir einen Witz“ an die KI geben soll.
4. Vector neu verbinden.

`german.env` wird beim Start automatisch gelesen. Schon gesetzte Umgebungsvariablen bleiben unangetastet.

## Linux und Raspberry Pi

```bash
bash setup-german.sh
```

Das Skript braucht kein Root. Piper und die Stimme landen in `~/wire-pod-data/`. Danach wire-pod mit `STT=vosk` bauen (`setup.sh` im Repo). Im Webinterface **German (DE)** und Knowledge Graph **Custom** mit **gpt-6-luna** wählen. Den Key nur dort eintragen, nicht in `german.env`.

## Nützliche Sätze

```text
Hey Vector, wie ist das Wetter
Hey Vector, wie spät ist es
Hey Vector, stell einen Timer für fünf Minuten
Hey Vector, ich habe eine Frage
Hey Vector, erzähl mir einen Witz
```

Der letzte Satz ist kein fester Befehl. Dafür muss der Intent-Graph an sein, dann antwortet GPT-6 Luna und Thorsten spricht es.

## Variablen

| Variable | Bedeutung |
| --- | --- |
| `PIPER_BIN` | Pfad zu `piper` bzw. `piper.exe` |
| `PIPER_MODEL` | Pfad zu `de_DE-thorsten-medium.onnx` |
| `TTS_SERVICE=piper` | Piper auch dann, wenn die STT-Sprache nicht de-DE ist |
| `KNOWLEDGE_PROVIDER` | `custom` für Azure Foundry |
| `KNOWLEDGE_ENDPOINT` | Basis-URL, ohne `/chat/completions` |
| `KNOWLEDGE_MODEL` | Deployment-Name `gpt-6-luna`, wenn im Knowledge Graph noch nichts steht |

## Credits

- wire-pod: [kercre123](https://github.com/kercre123/wire-pod)
- Deutsche Piper-Idee: [Kolle1979/wire-pod-german](https://github.com/Kolle1979/wire-pod-german)
- Stimme Thorsten, Piper: Rhasspy / Thorsten Müller
- Vosk: Alpha Cephei
