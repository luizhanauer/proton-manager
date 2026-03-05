package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	registryURL = "https://luizhanauer.github.io/proton-registry/api/smart_index.json"
	steamPath   = ".local/share/Steam/compatibilitytools.d"
)

// App struct
type App struct {
	ctx context.Context
}

// ProtonRelease é a estrutura para cada versão do Proton, com um campo extra para o status.
type ProtonRelease struct {
	Version     string `json:"version"`
	DownloadURL string `json:"url"`
	Size        int64  `json:"size"`
	Date        string `json:"date"`
	Major       string `json:"major"`
	Installed   bool   `json:"installed"`
}

// WriteCounter é usado para monitorar o progresso do download.
type WriteCounter struct {
	Total      int64
	Downloaded int64
	Version    string
	ctx        context.Context
}

// Write implementa io.Writer para contar bytes e emitir eventos de progresso.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Downloaded += int64(n)
	percentage := float64(wc.Downloaded) / float64(wc.Total) * 100

	// Emite um evento para o frontend com o progresso
	wailsruntime.EventsEmit(wc.ctx, "download-progress", map[string]interface{}{
		"version":    wc.Version,
		"percentage": percentage,
	})
	return n, nil
}

// NewApp cria uma nova instância da App
func NewApp() *App {
	return &App{}
}

// startup é chamado quando a aplicação inicia.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) getSteamInstallPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, steamPath), nil
}

// GetReleases busca a lista de versões do Proton e verifica quais estão instaladas.
func (a *App) GetReleases() ([]ProtonRelease, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao registro: %w", err)
	}
	defer resp.Body.Close()

	var releases []ProtonRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("falha ao decodificar dados do registro: %w", err)
	}

	installPath, err := a.getSteamInstallPath()
	if err != nil {
		return nil, err
	}

	installedVersions := a.getInstalledVersions(installPath)
	for i := range releases {
		if _, ok := installedVersions[releases[i].Version]; ok {
			releases[i].Installed = true
		}
	}

	// Ordena por versão, do mais novo para o mais antigo (assumindo que a API já faz isso)
	// Se não, uma lógica de ordenação mais robusta seria necessária aqui.

	return releases, nil
}

// InstallRelease instala uma versão específica do Proton.
func (a *App) InstallRelease(version, downloadURL string, size int64) error {
	installPath, err := a.getSteamInstallPath()
	if err != nil {
		return err
	}

	// Garante que o diretório de instalação exista
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		os.MkdirAll(installPath, 0755)
	}

	wailsruntime.EventsEmit(a.ctx, "install-status", map[string]interface{}{"version": version, "status": "downloading", "message": "Baixando " + version})

	tempFile := filepath.Join(os.TempDir(), "proton_download.tar.gz")
	defer os.Remove(tempFile)

	if err := a.downloadFile(version, downloadURL, tempFile, size); err != nil {
		return fmt.Errorf("erro no download: %w", err)
	}

	wailsruntime.EventsEmit(a.ctx, "install-status", map[string]interface{}{"version": version, "status": "installing", "message": "Extraindo " + version})

	if err := a.extractTarGz(tempFile, installPath); err != nil {
		// Limpa a pasta de destino em caso de falha na extração para evitar instalação corrompida
		os.RemoveAll(filepath.Join(installPath, version))
		return fmt.Errorf("erro na extração: %w", err)
	}

	wailsruntime.EventsEmit(a.ctx, "install-status", map[string]interface{}{"version": version, "status": "installed", "message": version + " instalado com sucesso!"})
	return nil
}

// UninstallRelease desinstala uma versão específica do Proton.
func (a *App) UninstallRelease(version string) error {
	installPath, err := a.getSteamInstallPath()
	if err != nil {
		return err
	}
	targetDir := filepath.Join(installPath, version)
	wailsruntime.EventsEmit(a.ctx, "install-status", map[string]interface{}{"version": version, "status": "uninstalling", "message": "Removendo " + version})

	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("falha ao remover %s: %w", version, err)
	}

	wailsruntime.EventsEmit(a.ctx, "install-status", map[string]interface{}{"version": version, "status": "uninstalled", "message": version + " removido."})
	return nil
}

// OpenInstallationFolder abre a pasta de compatibilidade da Steam no gerenciador de arquivos.
func (a *App) OpenInstallationFolder() error {
	installPath, err := a.getSteamInstallPath()
	if err != nil {
		return fmt.Errorf("não foi possível obter o caminho de instalação: %w", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", installPath)
	case "darwin":
		cmd = exec.Command("open", installPath)
	default: // "linux" e outros sistemas Unix-like
		cmd = exec.Command("xdg-open", installPath)
	}

	// .Start() executa o comando de forma não-bloqueante e retorna um erro se não conseguir iniciar.
	return cmd.Start()
}

// OpenURL abre uma URL no navegador padrão do sistema.
func (a *App) OpenURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// No Windows, 'start' é um comando interno do cmd.exe
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		// No macOS, 'open' é o comando para abrir URLs e arquivos.
		cmd = exec.Command("open", url)
	default: // "linux" e outros sistemas Unix-like
		// 'xdg-open' é o padrão para a maioria dos ambientes de desktop Linux.
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func (a *App) getInstalledVersions(path string) map[string]bool {
	versions := make(map[string]bool)
	entries, err := os.ReadDir(path)
	if err != nil {
		return versions
	}
	for _, e := range entries {
		if e.IsDir() {
			versions[e.Name()] = true
		}
	}
	return versions
}

func (a *App) downloadFile(version, url, destPath string, size int64) error {
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download falhou: status %s", resp.Status)
	}

	counter := &WriteCounter{
		Total:   size,
		Version: version,
		ctx:     a.ctx,
	}

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	return err
}

func (a *App) extractTarGz(tarPath, targetDir string) error {
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

		// Proteção contra "Zip Slip"
		if !strings.HasPrefix(target, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("conteúdo de tarball inválido: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}