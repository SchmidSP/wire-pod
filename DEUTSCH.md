# Wire-pod auf Deutsch

Dieser Fork von [kercre123/wire-pod](https://github.com/kercre123/wire-pod) spricht und versteht Deutsch. Die Idee entspricht [wire-pod-german](https://github.com/Kolle1979/wire-pod-german) (Piper, Stimme Thorsten, Vosk `de-DE`), ist aber in den echten Sprachpfad von wire-pod eingebaut und nicht nur als lose Dateien abgelegt.

Zusätzlich ist **GPT-6 Luna** (`gpt-6-luna`) das OpenAI-Standardmodell für den Knowledge Graph. GPT-6 Sol und die älteren Modelle bleiben auswählbar. Für Luna/Sol setzt wire-pod `reasoning_effort=low` und `verbosity=low`, damit Vector nicht ewig nachdenkt. Wenn das Modell mit dem Key nicht geht, fällt der Server auf `gpt-4o-mini` zurück.

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
2. Knowledge Graph: Anbieter **OpenAI**, Modell **gpt-6-luna**, eigenen Key eintragen.
3. **Intent-Graph** einschalten, wenn Vector auch freie Sätze wie „Erzähl mir einen Witz“ an die KI geben soll.
4. Vector neu verbinden.

`german.env` wird beim Start automatisch gelesen. Schon gesetzte Umgebungsvariablen bleiben unangetastet.

## Linux und Raspberry Pi

```bash
bash setup-german.sh
```

Das Skript braucht kein Root. Piper und die Stimme landen in `~/wire-pod-data/`. Danach wire-pod ganz normal mit `STT=vosk` bauen (`setup.sh` im Repo) und im Webinterface **German (DE)** sowie **gpt-6-luna** wählen.

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
| `KNOWLEDGE_MODEL` | wird übernommen, wenn im Knowledge Graph noch kein Modell steht |

## Credits

- wire-pod: [kercre123](https://github.com/kercre123/wire-pod)
- Deutsche Piper-Idee: [Kolle1979/wire-pod-german](https://github.com/Kolle1979/wire-pod-german)
- Stimme Thorsten, Piper: Rhasspy / Thorsten Müller
- Vosk: Alpha Cephei
