#!/bin/bash
set -e

# 1. Configurações de Caminhos (Padrão XDG)
BIN_DIR="$HOME/.local/bin"
APPS_DIR="$HOME/.local/share/applications"
ICONS_DIR="$HOME/.local/share/icons"
FINAL_BIN="$BIN_DIR/proton-manager"
DESKTOP_FILE="$APPS_DIR/proton-manager.desktop"

echo ">>> 🚀 Iniciando instalação do Proton Manager..."

# 2. Garantir que os diretórios existam
mkdir -p "$BIN_DIR" "$APPS_DIR" "$ICONS_DIR"

# 3. Instalar o Binário (Caminho corrigido para bin/proton-manager)
if [ -f "bin/proton-manager" ]; then
    cp "bin/proton-manager" "$FINAL_BIN"
    chmod +x "$FINAL_BIN"
    echo "✅ Binário instalado em $FINAL_BIN"
else
    echo "❌ Erro: Arquivo 'bin/proton-manager' não encontrado."
    exit 1
fi

# 4. Instalar o Ícone (Caminho corrigido para appicon.png)
if [ -f "appicon.png" ]; then
    cp "appicon.png" "$ICONS_DIR/proton-manager.png"
    ICON_PATH="$ICONS_DIR/proton-manager.png"
    echo "✅ Ícone instalado em $ICONS_DIR"
else
    echo "⚠️ Aviso: appicon.png não encontrada. Usando ícone genérico."
    ICON_PATH="applications-games"
fi

# 5. Criar o Atalho (.desktop)
# Adicionado StartupWMClass para garantir o ícone na dock durante a execução
cat <<EOF > "$DESKTOP_FILE"
[Desktop Entry]
Type=Application
Name=Proton Manager
Comment=Gerenciador de versões Proton
Exec="$FINAL_BIN"
Icon=$ICON_PATH
Terminal=false
Categories=Utility;Game;
Keywords=proton;steam;manager;
StartupNotify=true
StartupWMClass=proton-manager
EOF

chmod +x "$DESKTOP_FILE"

# 6. Atualizar Base de Dados e Cache de Ícones
update-desktop-database "$APPS_DIR" 2>/dev/null || true
gtk-update-icon-cache -f -t "$ICONS_DIR" 2>/dev/null || true

echo "--------------------------------------------------------"
echo "✨ Instalação concluída! Procure por 'Proton Manager' no menu."
echo "--------------------------------------------------------"