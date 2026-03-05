<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { GetReleases, InstallRelease, UninstallRelease, OpenInstallationFolder, OpenURL } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime';

// Tipagem para alinhar com o backend Go
interface ProtonRelease {
  version: string;
  url: string; // Mapeado de DownloadURL no JSON
  size: number;
  date: string;
  major: string;
  installed: boolean;
}

interface Status {
  status: string;
  progress: number;
  message: string;
}

const releases = ref<ProtonRelease[]>([]);
const statuses = reactive<Record<string, Status>>({});
const appLog = ref("Pronto para começar.");
const isLoading = ref(true);

onMounted(async () => {
  appLog.value = "Buscando versões do Proton...";
  try {
    const result = await GetReleases();
    releases.value = result || [];
    appLog.value = "Versões carregadas com sucesso.";
  } catch (err: any) {
    appLog.value = `Erro ao carregar: ${err}`;
  } finally {
    isLoading.value = false;
  }
});

EventsOn("download-progress", (data: { version: string, percentage: number }) => {
  statuses[data.version] = {
    status: "downloading",
    progress: data.percentage,
    message: `Baixando... ${Math.round(data.percentage)}%`
  };
});

EventsOn("install-status", (data: { version: string, status: string, message: string }) => {
  const currentStatus = statuses[data.version];
  statuses[data.version] = {
    status: data.status,
    progress: (data.status === 'installing' || data.status === 'installed') ? 100 : currentStatus?.progress || 0,
    message: data.message
  };

  appLog.value = data.message;

  if (data.status === 'installed' || data.status === 'uninstalled') {
    // Aguarda um pouco para o usuário ver a mensagem de sucesso e depois atualiza a lista
    setTimeout(async () => {
      const result = await GetReleases();
      releases.value = result || [];
      delete statuses[data.version];
    }, 2000);
  }
});

async function handleInstall(release: ProtonRelease) {
  statuses[release.version] = { status: "queued", progress: 0, message: "Na fila..." };
  try {
    await InstallRelease(release.version, release.url, release.size);
  } catch (err: any) {
    appLog.value = `Erro ao instalar ${release.version}: ${err}`;
    delete statuses[release.version];
  }
}

async function handleUninstall(version: string) {
  statuses[version] = { status: "queued", progress: 0, message: "Removendo..." };
  try {
    await UninstallRelease(version);
  } catch (err: any) {
    appLog.value = `Erro ao remover ${version}: ${err}`;
    delete statuses[version];
  }
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${parseFloat((bytes / Math.pow(1024, i)).toFixed(0))} ${['B', 'KB', 'MB', 'GB'][i]}`;
}

async function handleOpenFolder() {
  try {
    await OpenInstallationFolder();
  } catch (err: any) {
    appLog.value = `Erro ao abrir pasta: ${err}`;
  }
}

async function handleOpenURL(url: string) {
  try {
    await OpenURL(url);
  } catch (err: any) {
    appLog.value = `Erro ao abrir URL: ${err}`;
  }
}
</script>

<template>
  <div id="app-container">
    <header>
      <div class="header-content">
        <div class="title-container">
          <img src="/logo.svg" alt="Proton Manager Logo" class="logo" />
          <h1>Proton Manager</h1>
        </div>
        <button class="btn-icon" @click="handleOpenFolder" title="Abrir pasta de ferramentas de compatibilidade">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
        </button>
      </div>
    </header>

    <main>
      <div class="restart-warning">
        <p>⚠️ Lembre-se de reiniciar a Steam, Heroic ou outro cliente de jogos para que as alterações sejam aplicadas.</p>
      </div>

      <div v-if="isLoading">
        <p>Carregando...</p>
      </div>
      <div v-else class="release-list">
        <div v-for="release in releases" :key="release.version" class="release-card">
          
          <div class="release-main-info">
            <h2>{{ release.version }}</h2>
            <div class="release-info">
              <span>Data: {{ release.date }}</span> |
              <span>Tamanho: {{ formatSize(release.size) }}</span>
            </div>
          </div>

          <div class="release-actions-container">
            <div v-if="statuses[release.version]" class="status-display">
              <p>{{ statuses[release.version]?.message }}</p>
              <div v-if="statuses[release.version]?.status === 'downloading' || statuses[release.version]?.status === 'installing'" class="progress-bar">
                <div class="progress-bar-inner" :style="{ width: statuses[release.version]?.progress + '%' }"></div>
              </div>
            </div>

            <div v-else class="release-actions">
              <button v-if="release.installed" class="btn btn-uninstall" @click="handleUninstall(release.version)">
                Remover
              </button>
              <button v-else class="btn btn-install" @click="handleInstall(release)">
                Instalar
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>

    <footer>
      <p class="app-log">{{ appLog }}</p>
      <div class="footer-content">
        <div class="footer-links">
          <a href="https://github.com/luizhanauer/proton-manager" @click.prevent="handleOpenURL('https://github.com/luizhanauer/proton-manager')" class="btn-footer" title="Ver repositório no GitHub">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.91 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
            <span>Repositório</span>
          </a>
          <a href="https://www.paypal.com/donate/?hosted_button_id=SFR785YEYHC4E" @click.prevent="handleOpenURL('https://www.paypal.com/donate/?hosted_button_id=SFR785YEYHC4E')">
            <img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" class="donation-button">
          </a>
        </div>
        <p class="credits">Desenvolvido por <strong>Luiz Hanauer</strong></p>
      </div>
    </footer>
  </div>
</template>

<style>
/* Estilos Globais - Removi 'scoped' para aplicar o tema em toda a aplicação */
:root {
  font-family: "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;

  --color-bg: #121212;
  --color-bg-light: #1e1e1e;
  --color-bg-lighter: #2a2a2a;
  --color-text: #e0e0e0;
  --color-text-dim: #888;
  --color-primary: #00ff7f; /* Verde Neon */
  --color-primary-glow: rgba(0, 255, 127, 0.4);
  --color-danger: #ff4136;
  --color-danger-glow: rgba(255, 65, 54, 0.4);
  --color-border: #333;
}

body {
  margin: 0;
  background-color: var(--color-bg);
  color: var(--color-text);
  overflow-y: scroll;
}

#app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

main {
  padding: 2rem;
  flex-grow: 1;
}

header {
  background-color: var(--color-bg-light);
  padding: 1rem 2rem;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo {
  height: 2rem;
  width: auto;
}

.btn-icon {
  padding: 0.5rem;
  background-color: var(--color-bg-lighter);
  border: 1px solid var(--color-border);
  color: var(--color-text-dim);
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon:hover {
  background-color: var(--color-border);
  color: var(--color-text);
}

header h1 {
  margin: 0;
  font-size: 1.5rem;
  color: var(--color-primary);
  text-shadow: 0 0 5px var(--color-primary-glow);
}

footer {
  background-color: var(--color-bg-light);
  padding: 0.5rem 2rem;
  border-top: 1px solid var(--color-border);
  font-size: 0.8rem;
  color: var(--color-text-dim);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.app-log {
  margin: 0;
  text-align: center;
}

.footer-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
}

.footer-links {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.restart-warning {
  background-color: rgba(255, 193, 7, 0.1);
  border: 1px solid #ffc107;
  color: #ffc107;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
  text-align: center;
}

.restart-warning p {
  margin: 0;
}

.btn-footer {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.4rem 0.8rem;
  background-color: var(--color-bg-lighter);
  border: 1px solid var(--color-border);
  color: var(--color-text-dim);
  border-radius: 5px;
  text-decoration: none;
  transition: all 0.2s ease;
  font-size: 0.8rem;
}

.btn-footer:hover {
  background-color: var(--color-border);
  color: var(--color-primary);
  border-color: var(--color-primary);
}

.btn-footer svg {
  fill: currentColor;
}

.donation-button {
  height: 40px !important;
  width: 150px !important;
}

.release-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.release-card {
  background-color: var(--color-bg-light);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 0.75rem 1.5rem;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  transition: all 0.2s ease-in-out;
  gap: 1rem;
}

.release-card:hover {
  border-color: var(--color-primary);
  box-shadow: 0 0 10px var(--color-primary-glow);
}
.release-main-info {
  display: flex;
  align-items: center;
  flex-grow: 1;
  min-width: 0;
}

.release-card h2 {
  margin: 0 1rem 0 0;
  font-size: 1.1rem;
  color: var(--color-text);
  white-space: nowrap;
}

.release-info {
  color: var(--color-text-dim);
  font-size: 0.85rem;
  margin-bottom: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.release-actions-container {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-shrink: 0;
  width: 240px;
}

.release-actions {
  display: flex;
  gap: 0.5rem;
}

.status-display {
  width: 100%;
  text-align: right;
}

.status-display p {
  margin: 0;
  font-size: 0.9rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 5px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s ease;
  text-transform: uppercase;
  letter-spacing: 1px;
  flex-grow: 1;
  background-color: transparent;
  border: 1px solid;
}

.btn-install {
  border-color: var(--color-primary);
  color: var(--color-primary);
  box-shadow: 0 0 3px var(--color-primary-glow), inset 0 0 3px var(--color-primary-glow);
}

.btn-install:hover:not(:disabled) {
  background-color: var(--color-primary-glow);
  color: white;
  box-shadow: 0 0 12px var(--color-primary-glow), inset 0 0 8px var(--color-primary-glow);
}

.btn-uninstall {
  border-color: var(--color-danger);
  color: var(--color-danger);
  box-shadow: 0 0 3px var(--color-danger-glow), inset 0 0 3px var(--color-danger-glow);
}

.btn-uninstall:hover:not(:disabled) {
  background-color: var(--color-danger-glow);
  color: white;
  box-shadow: 0 0 12px var(--color-danger-glow), inset 0 0 8px var(--color-danger-glow);
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.progress-bar {
  width: 100%;
  height: 5px;
  background-color: var(--color-border);
  border-radius: 5px;
  margin-top: 0.5rem;
  overflow: hidden;
}

.progress-bar-inner {
  height: 100%;
  background-color: var(--color-primary);
  transition: width 0.2s ease-out;
  box-shadow: 0 0 5px var(--color-primary-glow);
}

.credits {
  font-size: 0.75rem;
  margin: 0;
}
</style>
