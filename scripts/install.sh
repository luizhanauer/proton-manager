#!/bin/sh
set -e

cd "$(dirname "$0")"

BIN_DIR="$HOME/.local/bin"
APPS_DIR="$HOME/.local/share/applications"
ICONS_DIR="$HOME/.local/share/icons"

FINAL_BIN="$BIN_DIR/proton-manager"

echo ">>> 🚀 Instalando Proton Manager..."

mkdir -p "$BIN_DIR"
mkdir -p "$APPS_DIR"
mkdir -p "$ICONS_DIR"

if [ ! -f "bin/proton-manager" ]; then
    echo "❌ Erro: Binário não encontrado na pasta extraída."
    exit 1
fi

cp bin/proton-manager "$FINAL_BIN"
chmod +x "$FINAL_BIN"

ICON_NAME="applications-games"
if [ -f "appicon.png" ]; then
    cp appicon.png "$ICONS_DIR/proton-manager.png"
    ICON_NAME="$ICONS_DIR/proton-manager.png"
fi

cat <<EOF > /tmp/protonmanager.desktop
[Desktop Entry]
Type=Application
Name=Proton Manager
Comment=Gerenciador de versões customizadas do Proton-GE
Exec=$FINAL_BIN
Icon=$ICON_NAME
Terminal=false
Categories=Utility;System;Game;
Keywords=proton;steam;gaming;manager;
Hidden=false
NoDisplay=false
EOF

cp /tmp/protonmanager.desktop "$APPS_DIR/protonmanager.desktop"

update-desktop-database "$APPS_DIR" 2>/dev/null || true

echo "--------------------------------------------------------"
echo "✅ Proton Manager instalado com sucesso!"
echo "📂 Menu: Disponível em 'Mostrar Aplicativos'"
echo "--------------------------------------------------------"