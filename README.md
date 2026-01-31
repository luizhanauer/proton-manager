# 🚀 Proton Manager

> O jeito mais rápido e seguro de gerenciar versões do GE-Proton no Linux.

![GitHub release (latest by date)](https://img.shields.io/github/v/release/luizhanauer/proton-manager)
![Platform](https://img.shields.io/badge/platform-linux-linux)
![License](https://img.shields.io/github/license/luizhanauer/proton-manager)
![Go Version](https://img.shields.io/github/go-mod/go-version/luizhanauer/proton-manager)

O **Proton Manager** é uma ferramenta CLI (Linha de Comando) interativa escrita em Go. Ela se conecta ao [Proton Registry](https://github.com/luizhanauer/proton-registry) para baixar, instalar e limpar versões do GE-Proton na sua Steam ou Heroic Games Launcher, contornando limitações da API do GitHub e automatizando o processo.

## ✨ Funcionalidades

* **🖥️ Interface Interativa (TUI):** Selecione múltiplas versões com caixas de seleção, navegue com setas e veja o status de instalação.
* **⚡ Downloads Otimizados:** Barra de progresso visual com velocidade de download e tempo estimado.
* **🧹 Sincronização Inteligente:** O que você **desmarca** no menu é removido do disco. Mantém seu sistema limpo automaticamente.
* **🤖 Modo Automação:** Flag `--latest` para atualizar para a última versão sem interação humana (ideal para scripts de inicialização).
* **🛡️ Segurança:** Proteção contra exclusão acidental de pastas da Steam (`LegacyRuntime`, `Soldier`, etc). Apenas versões gerenciadas são tocadas.
* **📂 Alinhamento Perfeito:** Colunas formatadas dinamicamente para facilitar a leitura no terminal.

## 📦 Instalação

### Opção 1: Binário Pré-compilado (Recomendado)
Baixe a última versão na aba [Releases](https://github.com/luizhanauer/proton-manager/releases) e instale:

```bash
# Exemplo (ajuste a versão conforme necessário)
wget [https://github.com/luizhanauer/proton-manager/releases/download/v0.1.0/proton-manager](https://github.com/luizhanauer/proton-manager/releases/download/v0.1.0/proton-manager)
chmod +x proton-manager
sudo mv proton-manager /usr/local/bin/