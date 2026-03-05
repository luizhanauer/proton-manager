#!/bin/sh
set -e

REPO="luizhanauer/proton-manager"
BINARY_NAME="proton-manager"

ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ASSET_ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
    ASSET_ARCH="arm64"
else
    echo "❌ Arquitetura não suportada: $ARCH"
    exit 1
fi

echo ">>> 📦 Instalador via Rede: Proton Manager"
echo ">>> Fonte: https://github.com/$REPO"

if ! command -v curl >/dev/null; then 
    echo "❌ Erro: 'curl' é necessário."
    exit 1
fi

if ! command -v tar >/dev/null; then 
    echo "❌ Erro: 'tar' é necessário."
    exit 1
fi

TMP_DIR=$(mktemp -d)

# Garante que a pasta temporária será removida ao sair, mesmo em caso de erro
trap 'rm -rf "$TMP_DIR"' EXIT

FILENAME="${BINARY_NAME}_linux_${ASSET_ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/latest/download/${FILENAME}"

echo ">>> ⬇️  Baixando release..."
if ! curl -f -L "$URL" -o "$TMP_DIR/$FILENAME"; then
    echo "❌ Erro ao baixar release. Verifique se a tag de release existe no GitHub."
    exit 1
fi

echo ">>> 📂 Extraindo..."
tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

echo ">>> 🚀 Iniciando script de instalação..."
cd "$TMP_DIR"
chmod +x install.sh
sh ./install.sh

echo ">>> ✅ Setup finalizado."