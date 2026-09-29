# Lokales Deutsch für diesen wire-pod-Fork unter Windows: Piper + Thorsten.
# Aufruf: powershell -ExecutionPolicy Bypass -File .\install-german.ps1
$ErrorActionPreference = "Stop"

$Data = Join-Path $env:USERPROFILE "wire-pod-data"
$PiperDir = Join-Path $Data "piper"
$VoiceDir = Join-Path $Data "voices"
$PiperTag = "2023.11.14-2"
$ZipUrl = "https://github.com/rhasspy/piper/releases/download/$PiperTag/piper_windows_amd64.zip"
$VoiceBase = "https://huggingface.co/rhasspy/piper-voices/resolve/v1.0.0/de/de_DE/thorsten/medium"

New-Item -ItemType Directory -Force -Path $PiperDir, $VoiceDir | Out-Null

$PiperExe = Join-Path $PiperDir "piper.exe"
if (-not (Test-Path $PiperExe)) {
    Write-Host "Lade Piper $PiperTag..."
    $zip = Join-Path $env:TEMP "piper_windows_amd64.zip"
    Invoke-WebRequest -Uri $ZipUrl -OutFile $zip
    $extract = Join-Path $env:TEMP "piper-extract"
    if (Test-Path $extract) { Remove-Item -Recurse -Force $extract }
    Expand-Archive -Path $zip -DestinationPath $extract -Force
    $found = Get-ChildItem -Path $extract -Filter piper.exe -Recurse | Select-Object -First 1
    if (-not $found) { throw "piper.exe nicht im Archiv" }
    Copy-Item -Path (Join-Path $found.Directory.FullName "*") -Destination $PiperDir -Recurse -Force
}

$Voice = Join-Path $VoiceDir "de_DE-thorsten-medium.onnx"
if (-not (Test-Path $Voice)) {
    Write-Host "Lade Stimme Thorsten..."
    Invoke-WebRequest -Uri "$VoiceBase/de_DE-thorsten-medium.onnx" -OutFile $Voice
    Invoke-WebRequest -Uri "$VoiceBase/de_DE-thorsten-medium.onnx.json" -OutFile "$Voice.json"
}

$EnvFile = Join-Path $Data "german.env"
@"
PIPER_BIN=$PiperExe
PIPER_MODEL=$Voice
TTS_SERVICE=piper
STT_LANGUAGE=de-DE
KNOWLEDGE_PROVIDER=custom
KNOWLEDGE_MODEL=gpt-6-luna
KNOWLEDGE_ENDPOINT=https://ksgptsweden.cognitiveservices.azure.com/openai/v1
"@ | Set-Content -Path $EnvFile -Encoding ascii

Write-Host "Fertig."
Write-Host "Piper: $PiperExe"
Write-Host "Stimme: $Voice"
Write-Host "Umgebung: $EnvFile"
Write-Host ""
Write-Host "Wire-pod aus diesem Fork neu bauen und starten."
Write-Host "Im Webinterface: Sprache German (DE), Knowledge Graph Provider Custom, Endpoint die Azure-v1-URL, Modell gpt-6-luna, Azure-Key eintragen. Intent-Graph an, LLM-Actions aus."
Write-Host "Siehe DEUTSCH.md"
