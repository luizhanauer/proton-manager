#!/bin/sh
set -e

# Garante que o script rode no diretório onde está localizado
cd "$(dirname "$0")"

# --- Configurações de Caminhos ---
BIN_DIR="$HOME/.local/bin"
APPS_DIR="$HOME/.local/share/applications"
ICONS_DIR="$HOME/.local/share/icons"
FINAL_BIN="$BIN_DIR/proton-manager"
DESKTOP_FILE="/tmp/protonmanager.desktop"

echo ">>> 🚀 Instalando Proton Manager..."

# 1. Preparação do ambiente
mkdir -p "$BIN_DIR"
mkdir -p "$APPS_DIR"
mkdir -p "$ICONS_DIR"

# 2. Instalação do Binário
if [ ! -f "bin/proton-manager" ]; then
    echo "❌ Erro: Binário não encontrado na pasta extraída (bin/proton-manager)."
    exit 1
fi

cp bin/proton-manager "$FINAL_BIN"
chmod +x "$FINAL_BIN"

# 3. Gerenciamento do Ícone
# Se houver um ícone local, instalamos com um nome único para evitar conflitos
ICON_PATH="applications-games" 
if [ -f "appicon.png" ]; then
    cp appicon.png "$ICONS_DIR/proton-manager-icon.png"
    # Usar o caminho absoluto ajuda o GNOME/KDE a localizar o arquivo fora do tema padrão
    ICON_PATH="$ICONS_DIR/proton-manager-icon.png"
fi

# 4. Criação do Atalho (.desktop)
# Removido espaços desnecessários e garantido que o Icon aponte para o local correto
cat <<EOF > "$DESKTOP_FILE"
[Desktop Entry]
Type=Application
Name=Proton Manager
Comment=Gerenciador de versões customizadas do Proton-GE
Exec="$FINAL_BIN"
Icon=$ICON_PATH
Terminal=false
Categories=Utility;System;Game;
Keywords=proton;steam;gaming;manager;
StartupNotify=true
EOF

cp "$DESKTOP_FILE" "$APPS_DIR/protonmanager.desktop"

# 5. Atualização dos Bancos de Dados do Sistema
# Adicionado update-icon-cache para garantir que o sistema veja o novo ícone
update-desktop-database "$APPS_DIR" 2>/dev/null || true
gtk-update-icon-cache -f -t "$ICONS_DIR" 2>/dev/null || true

echo "--------------------------------------------------------"
echo "✅ Proton Manager instalado com sucesso!"
echo "📂 Menu: Disponível em 'Mostrar Aplicativos'"
echo "--------------------------------------------------------"