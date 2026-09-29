#!/bin/bash
# Lokales Deutsch für diesen wire-pod-Fork: Piper + Thorsten.
# Braucht kein Root. Danach wire-pod selbst bauen und im Webinterface de-DE wählen.
set -euo pipefail

DATA="${HOME}/wire-pod-data"
PIPER_DIR="${DATA}/piper"
VOICE_DIR="${DATA}/voices"
PIPER_TAG="2023.11.14-2"
ARCH="$(uname -m)"

case "${ARCH}" in
  x86_64|amd64) ASSET="piper_linux_x86_64.tar.gz" ;;
  aarch64|arm64) ASSET="piper_linux_aarch64.tar.gz" ;;
  armv7l|armv6l) ASSET="piper_linux_armv7l.tar.gz" ;;
  *) echo "Architektur ${ARCH} wird nicht unterstützt."; exit 1 ;;
esac

mkdir -p "${PIPER_DIR}" "${VOICE_DIR}"

if [[ ! -x "${PIPER_DIR}/piper" ]]; then
  echo "Lade Piper ${PIPER_TAG} (${ASSET})..."
  TMP="$(mktemp -d)"
  curl -fL --retry 3 -o "${TMP}/piper.tgz" \
    "https://github.com/rhasspy/piper/releases/download/${PIPER_TAG}/${ASSET}"
  tar -xzf "${TMP}/piper.tgz" -C "${TMP}"
  rm -rf "${PIPER_DIR}"
  mv "${TMP}/piper" "${PIPER_DIR}"
  rm -rf "${TMP}"
  chmod +x "${PIPER_DIR}/piper"
fi

VOICE="${VOICE_DIR}/de_DE-thorsten-medium.onnx"
BASE="https://huggingface.co/rhasspy/piper-voices/resolve/v1.0.0/de/de_DE/thorsten/medium"
if [[ ! -f "${VOICE}" ]]; then
  echo "Lade Stimme Thorsten..."
  curl -fL --retry 3 -o "${VOICE}" "${BASE}/de_DE-thorsten-medium.onnx"
  curl -fL --retry 3 -o "${VOICE}.json" "${BASE}/de_DE-thorsten-medium.onnx.json"
fi

cat > "${DATA}/german.env" <<EOF
PIPER_BIN=${PIPER_DIR}/piper
PIPER_MODEL=${VOICE}
TTS_SERVICE=piper
KNOWLEDGE_MODEL=gpt-6-luna
STT_LANGUAGE=de-DE
EOF

echo "Teste Piper..."
echo "Hallo. Ich bin Vector und spreche Deutsch." | "${PIPER_DIR}/piper" \
  --model "${VOICE}" --output_file "${DATA}/piper-test.wav"
echo "Fertig: ${DATA}/piper-test.wav"
echo "german.env: ${DATA}/german.env"
echo
echo "Als Nächstes im wire-pod-Webinterface:"
echo "  1. Sprache German (DE)"
echo "  2. Knowledge Graph: OpenAI, Modell gpt-6-luna, Intent-Graph an"
echo "Siehe DEUTSCH.md"
