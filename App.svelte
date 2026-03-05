<script>
  import { onMount } from 'svelte';
  import { writable } from 'svelte/store';
  import { GetReleases, InstallRelease, UninstallRelease } from '../../wailsjs/go/main/App';
  import { EventsOn } from '../../wailsjs/runtime';

  let releases = writable([]);
  let statuses = writable({}); // { "version": { status: "downloading", progress: 50, message: "..." } }
  let appLog = writable("Pronto para começar.");
  let isLoading = writable(true);

  onMount(async () => {
    $appLog = "Buscando versões do Proton...";
    try {
      const result = await GetReleases();
      releases.set(result || []);
      $appLog = "Versões carregadas com sucesso.";
    } catch (err) {
      $appLog = `Erro ao carregar: ${err}`;
    } finally {
      isLoading.set(false);
    }
  });

  EventsOn("download-progress", (data) => {
    statuses.update(s => {
      s[data.version] = { 
        status: "downloading", 
        progress: data.percentage,
        message: `Baixando... ${Math.round(data.percentage)}%`
      };
      return s;
    });
  });

  EventsOn("install-status", (data) => {
    statuses.update(s => {
      s[data.version] = { 
        status: data.status, 
        progress: (data.status === 'installing' || data.status === 'installed') ? 100 : s[data.version]?.progress || 0,
        message: data.message
      };
      return s;
    });
    
    $appLog = data.message;

    if (data.status === 'installed' || data.status === 'uninstalled') {
      // Aguarda um pouco para o usuário ver a mensagem de sucesso e depois limpa o status
      setTimeout(async () => {
        const result = await GetReleases();
        releases.set(result || []);
        statuses.update(s => {
          delete s[data.version];
          return s;
        });
      }, 2000);
    }
  });

  async function handleInstall(release) {
    statuses.update(s => {
      s[release.version] = { status: "queued", progress: 0, message: "Na fila..." };
      return s;
    });
    try {
      await InstallRelease(release.version, release.url, release.size);
    } catch (err) {
      $appLog = `Erro ao instalar ${release.version}: ${err}`;
      statuses.update(s => {
        delete s[release.version];
        return s;
      });
    }
  }

  async function handleUninstall(version) {
    statuses.update(s => {
      s[version] = { status: "queued", progress: 0, message: "Removendo..." };
      return s;
    });
    try {
      await UninstallRelease(version);
    } catch (err) {
      $appLog = `Erro ao remover ${version}: ${err}`;
      statuses.update(s => {
        delete s[version];
        return s;
      });
    }
  }

  function formatSize(bytes) {
    if (bytes === 0) return '0 B';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${parseFloat((bytes / Math.pow(1024, i)).toFixed(0))} ${['B', 'KB', 'MB', 'GB'][i]}`;
  }
</script>

<div id="app">
  <header>
    <h1>🦀 Proton Manager</h1>
  </header>

  <main>
    {#if $isLoading}
      <p>Carregando...</p>
    {:else}
      <div class="release-list">
        {#each $releases as release (release.version)}
          <div class="release-card">
            <h2>{release.version}</h2>
            <div class="release-info">
              <span>Data: {release.date}</span> | 
              <span>Tamanho: {formatSize(release.size)}</span>
            </div>
            
            {#if $statuses[release.version]}
              <div class="status-display">
                <p>{$statuses[release.version].message}</p>
                {#if $statuses[release.version].status === 'downloading' || $statuses[release.version].status === 'installing'}
                  <div class="progress-bar">
                    <div class="progress-bar-inner" style="width: {$statuses[release.version].progress}%"></div>
                  </div>
                {/if}
              </div>
            {:else}
              <div class="release-actions">
                {#if release.installed}
                  <button class="btn btn-uninstall" on:click={() => handleUninstall(release.version)}>
                    Remover
                  </button>
                {:else}
                  <button class="btn btn-install" on:click={() => handleInstall(release)}>
                    Instalar
                  </button>
                {/if}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </main>

  <footer>
    <p>{$appLog}</p>
  </footer>
</div>