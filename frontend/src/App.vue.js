import { ref, onMounted, reactive } from 'vue';
import { GetReleases, InstallRelease, UninstallRelease, OpenInstallationFolder, OpenURL } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime';
const releases = ref([]);
const statuses = reactive({});
const appLog = ref("Pronto para começar.");
const isLoading = ref(true);
onMounted(async () => {
    appLog.value = "Buscando versões do Proton...";
    try {
        const result = await GetReleases();
        releases.value = result || [];
        appLog.value = "Versões carregadas com sucesso.";
    }
    catch (err) {
        appLog.value = `Erro ao carregar: ${err}`;
    }
    finally {
        isLoading.value = false;
    }
});
EventsOn("download-progress", (data) => {
    statuses[data.version] = {
        status: "downloading",
        progress: data.percentage,
        message: `Baixando... ${Math.round(data.percentage)}%`
    };
});
EventsOn("install-status", (data) => {
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
async function handleInstall(release) {
    statuses[release.version] = { status: "queued", progress: 0, message: "Na fila..." };
    try {
        await InstallRelease(release.version, release.url, release.size);
    }
    catch (err) {
        appLog.value = `Erro ao instalar ${release.version}: ${err}`;
        delete statuses[release.version];
    }
}
async function handleUninstall(version) {
    statuses[version] = { status: "queued", progress: 0, message: "Removendo..." };
    try {
        await UninstallRelease(version);
    }
    catch (err) {
        appLog.value = `Erro ao remover ${version}: ${err}`;
        delete statuses[version];
    }
}
function formatSize(bytes) {
    if (bytes === 0)
        return '0 B';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${parseFloat((bytes / Math.pow(1024, i)).toFixed(0))} ${['B', 'KB', 'MB', 'GB'][i]}`;
}
async function handleOpenFolder() {
    try {
        await OpenInstallationFolder();
    }
    catch (err) {
        appLog.value = `Erro ao abrir pasta: ${err}`;
    }
}
async function handleOpenURL(url) {
    try {
        await OpenURL(url);
    }
    catch (err) {
        appLog.value = `Erro ao abrir URL: ${err}`;
    }
}
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    id: "app-container",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.header, __VLS_intrinsics.header)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "header-content" },
});
/** @type {__VLS_StyleScopedClasses['header-content']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "title-container" },
});
/** @type {__VLS_StyleScopedClasses['title-container']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.img)({
    src: "/logo.svg",
    alt: "Proton Manager Logo",
    ...{ class: "logo" },
});
/** @type {__VLS_StyleScopedClasses['logo']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
    ...{ onClick: (__VLS_ctx.handleOpenFolder) },
    ...{ class: "btn-icon" },
    title: "Abrir pasta de ferramentas de compatibilidade",
});
/** @type {__VLS_StyleScopedClasses['btn-icon']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.svg, __VLS_intrinsics.svg)({
    xmlns: "http://www.w3.org/2000/svg",
    width: "20",
    height: "20",
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    'stroke-width': "2",
    'stroke-linecap': "round",
    'stroke-linejoin': "round",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.path, __VLS_intrinsics.path)({
    d: "M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.main, __VLS_intrinsics.main)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "restart-warning" },
});
/** @type {__VLS_StyleScopedClasses['restart-warning']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({});
if (__VLS_ctx.isLoading) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({});
}
else {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "release-list" },
    });
    /** @type {__VLS_StyleScopedClasses['release-list']} */ ;
    for (const [release] of __VLS_vFor((__VLS_ctx.releases))) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            key: (release.version),
            ...{ class: "release-card" },
        });
        /** @type {__VLS_StyleScopedClasses['release-card']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "release-main-info" },
        });
        /** @type {__VLS_StyleScopedClasses['release-main-info']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.h2, __VLS_intrinsics.h2)({});
        (release.version);
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "release-info" },
        });
        /** @type {__VLS_StyleScopedClasses['release-info']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
        (release.date);
        __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
        (__VLS_ctx.formatSize(release.size));
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "release-actions-container" },
        });
        /** @type {__VLS_StyleScopedClasses['release-actions-container']} */ ;
        if (__VLS_ctx.statuses[release.version]) {
            __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
                ...{ class: "status-display" },
            });
            /** @type {__VLS_StyleScopedClasses['status-display']} */ ;
            __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({});
            (__VLS_ctx.statuses[release.version]?.message);
            if (__VLS_ctx.statuses[release.version]?.status === 'downloading' || __VLS_ctx.statuses[release.version]?.status === 'installing') {
                __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
                    ...{ class: "progress-bar" },
                });
                /** @type {__VLS_StyleScopedClasses['progress-bar']} */ ;
                __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
                    ...{ class: "progress-bar-inner" },
                    ...{ style: ({ width: __VLS_ctx.statuses[release.version]?.progress + '%' }) },
                });
                /** @type {__VLS_StyleScopedClasses['progress-bar-inner']} */ ;
            }
        }
        else {
            __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
                ...{ class: "release-actions" },
            });
            /** @type {__VLS_StyleScopedClasses['release-actions']} */ ;
            if (release.installed) {
                __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
                    ...{ onClick: (...[$event]) => {
                            if (!!(__VLS_ctx.isLoading))
                                return;
                            if (!!(__VLS_ctx.statuses[release.version]))
                                return;
                            if (!(release.installed))
                                return;
                            __VLS_ctx.handleUninstall(release.version);
                            // @ts-ignore
                            [handleOpenFolder, isLoading, releases, formatSize, statuses, statuses, statuses, statuses, statuses, handleUninstall,];
                        } },
                    ...{ class: "btn btn-uninstall" },
                });
                /** @type {__VLS_StyleScopedClasses['btn']} */ ;
                /** @type {__VLS_StyleScopedClasses['btn-uninstall']} */ ;
            }
            else {
                __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
                    ...{ onClick: (...[$event]) => {
                            if (!!(__VLS_ctx.isLoading))
                                return;
                            if (!!(__VLS_ctx.statuses[release.version]))
                                return;
                            if (!!(release.installed))
                                return;
                            __VLS_ctx.handleInstall(release);
                            // @ts-ignore
                            [handleInstall,];
                        } },
                    ...{ class: "btn btn-install" },
                });
                /** @type {__VLS_StyleScopedClasses['btn']} */ ;
                /** @type {__VLS_StyleScopedClasses['btn-install']} */ ;
            }
        }
        // @ts-ignore
        [];
    }
}
__VLS_asFunctionalElement1(__VLS_intrinsics.footer, __VLS_intrinsics.footer)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
    ...{ class: "app-log" },
});
/** @type {__VLS_StyleScopedClasses['app-log']} */ ;
(__VLS_ctx.appLog);
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "footer-content" },
});
/** @type {__VLS_StyleScopedClasses['footer-content']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "footer-links" },
});
/** @type {__VLS_StyleScopedClasses['footer-links']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({
    ...{ onClick: (...[$event]) => {
            __VLS_ctx.handleOpenURL('https://github.com/luizhanauer/proton-manager');
            // @ts-ignore
            [appLog, handleOpenURL,];
        } },
    href: "https://github.com/luizhanauer/proton-manager",
    ...{ class: "btn-footer" },
    title: "Ver repositório no GitHub",
});
/** @type {__VLS_StyleScopedClasses['btn-footer']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.svg, __VLS_intrinsics.svg)({
    xmlns: "http://www.w3.org/2000/svg",
    width: "18",
    height: "18",
    viewBox: "0 0 24 24",
    fill: "currentColor",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.path)({
    d: "M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.91 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({
    ...{ onClick: (...[$event]) => {
            __VLS_ctx.handleOpenURL('https://www.paypal.com/donate/?hosted_button_id=SFR785YEYHC4E');
            // @ts-ignore
            [handleOpenURL,];
        } },
    href: "https://www.paypal.com/donate/?hosted_button_id=SFR785YEYHC4E",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.img)({
    src: "https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png",
    alt: "Buy Me A Coffee",
    ...{ class: "donation-button" },
});
/** @type {__VLS_StyleScopedClasses['donation-button']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
    ...{ class: "credits" },
});
/** @type {__VLS_StyleScopedClasses['credits']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.strong, __VLS_intrinsics.strong)({});
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
//# sourceMappingURL=App.vue.js.map