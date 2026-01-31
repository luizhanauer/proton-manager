package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/schollz/progressbar/v3"
)

var (
	version = "dev" // Será substituído pelo GitHub Actions (ex: v0.1.0)
	date    = "unknown"
)

const (
	registryURL = "https://raw.githubusercontent.com/luizhanauer/proton-registry/main/api/smart_index.json"
	steamPath   = ".local/share/Steam/compatibilitytools.d"
)

type ProtonEntry struct {
	Version     string `json:"version"`
	DownloadURL string `json:"url"`
	Size        int64  `json:"size"`
	Date        string `json:"date"`
	Major       string `json:"major"`
}

func main() {
	latestFlag := flag.Bool("latest", false, "Instala automaticamente a última versão")

	versionFlag := flag.Bool("version", false, "Exibe a versão atual do gerenciador")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("🦀 Proton Manager %s (Compilado em: %s)\n", version, date)
		return
	}

	homeDir, _ := os.UserHomeDir()
	installPath := filepath.Join(homeDir, steamPath)

	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		os.MkdirAll(installPath, 0755)
	}

	fmt.Println("🛰️  Conectando ao Proton Registry...")
	releases, err := fetchFromRegistry()
	if err != nil {
		fmt.Printf("❌ Erro ao buscar dados: %v\n", err)
		return
	}

	installed := getInstalledVersions(installPath)

	// --- MODO AUTOMÁTICO ---
	if *latestFlag {
		latest := releases[0]
		if installed[latest.Version] {
			fmt.Printf("✅ A última versão (%s) já está instalada.\n", latest.Version)
			return
		}
		processInstall(latest, installPath)
		return
	}

	// --- MODO INTERATIVO ---

	// 1. Calcular o tamanho máximo da string de versão para alinhamento perfeito
	maxLen := 0
	for _, r := range releases {
		if len(r.Version) > maxLen {
			maxLen = len(r.Version)
		}
	}
	// Cria um formato dinâmico (ex: "%-18s | %-10s | %-8s | %s")
	formatStr := fmt.Sprintf("%%-%ds | %%-10s | %%-8s | %%s", maxLen+2)

	var options []string
	var defaultSelected []string
	var labelToVersion = make(map[string]string)
	var versionToEntry = make(map[string]ProtonEntry)

	for _, r := range releases {
		sizeMB := fmt.Sprintf("%.0f MB", float64(r.Size)/(1024*1024))
		status := ""
		if installed[r.Version] {
			status = "✅ INSTALADO"
		}

		// Usa o formato calculado para garantir colunas retas
		label := fmt.Sprintf(formatStr, r.Version, r.Date, sizeMB, status)

		options = append(options, label)
		labelToVersion[label] = r.Version
		versionToEntry[r.Version] = r

		if installed[r.Version] {
			defaultSelected = append(defaultSelected, label)
		}
	}

	// 2. Exibir Instruções Claras ANTES do menu
	printHeader()

	var selectedLabels []string
	prompt := &survey.MultiSelect{
		Message:  "Gerenciar versões:",
		Options:  options,
		Default:  defaultSelected,
		PageSize: 15,
		// Removemos o 'Help' escondido e deixamos explícito acima
	}

	// Captura a seleção
	if err := survey.AskOne(prompt, &selectedLabels); err != nil {
		fmt.Println("Operação cancelada.")
		return
	}

	// 3. Processamento
	toKeep := make(map[string]bool)

	// Instalação
	for _, label := range selectedLabels {
		version := labelToVersion[label]
		toKeep[version] = true

		if !installed[version] {
			processInstall(versionToEntry[version], installPath)
		}
	}

	// Limpeza (Com proteção)
	for v := range installed {
		_, isRegistryVersion := versionToEntry[v]
		if isRegistryVersion && !toKeep[v] {
			fmt.Printf("🗑️  Removendo versão desmarcada: %s\n", v)
			os.RemoveAll(filepath.Join(installPath, v))
		}
	}

	fmt.Println("\n✨ Sincronização concluída!")
}

// --- Função Visual ---
func printHeader() {
	fmt.Println("\n📋  INSTRUÇÕES DE USO:")
	fmt.Println("   • [Espaço] marcar/desmarcar versões.")
	fmt.Println("   • [Enter]  confirmar as alterações.")
	fmt.Println("\n⚠️  ATENÇÃO: Versões desmarcadas serão DELETADAS do disco.")
	fmt.Println("\n⚠️  ATENÇÃO: Versões marcadas serão INSTALADAS do disco.")
	fmt.Println(strings.Repeat("-", 60))
}

// --- Funções Auxiliares (Inalteradas) ---

func fetchFromRegistry() ([]ProtonEntry, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var entries []ProtonEntry
	return entries, json.NewDecoder(resp.Body).Decode(&entries)
}

func getInstalledVersions(path string) map[string]bool {
	m := make(map[string]bool)
	entries, _ := os.ReadDir(path)
	for _, e := range entries {
		if e.IsDir() {
			m[e.Name()] = true
		}
	}
	return m
}

func processInstall(entry ProtonEntry, targetDir string) {
	finalDest := filepath.Join(targetDir, entry.Version)

	if _, err := os.Stat(finalDest); !os.IsNotExist(err) {
		fmt.Printf("🧹 Limpando instalação anterior de %s...\n", entry.Version)
		if err := os.RemoveAll(finalDest); err != nil {
			fmt.Printf("❌ Erro ao limpar: %v (Feche a Steam)\n", err)
			return
		}
	}

	tempFile := filepath.Join(os.TempDir(), "proton_download.tar.gz")

	if err := downloadWithProgress(entry, tempFile); err != nil {
		fmt.Printf("❌ Erro no download: %v\n", err)
		return
	}

	fmt.Printf("\n📦 Extraindo...\n")
	if err := extractTarGz(tempFile, targetDir); err != nil {
		fmt.Printf("❌ Erro na extração: %v\n", err)
		os.Remove(tempFile)
		return
	}

	os.Remove(tempFile)
	fmt.Printf("✅ %s instalado!\n", entry.Version)
}

func downloadWithProgress(entry ProtonEntry, destPath string) error {
	req, _ := http.NewRequest("GET", entry.DownloadURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, _ := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetDescription(fmt.Sprintf("⬇️  Baixando %-15s", entry.Version)),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(20),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() { fmt.Fprint(os.Stderr, "\n") }),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer: "=", SaucerHead: ">", SaucerPadding: " ", BarStart: "[", BarEnd: "]",
		}),
	)

	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	return err
}

func extractTarGz(tarPath, targetDir string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(targetDir, header.Name)
		if !strings.HasPrefix(target, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			continue
		}
		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			io.Copy(f, tr)
			f.Close()
		}
	}
	return nil
}
