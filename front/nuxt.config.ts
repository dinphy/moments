// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03',
    devtools: {enabled: false},
    modules: ["@nuxt/ui", '@vueuse/nuxt', 'dayjs-nuxt'],
    ssr: false,
    dayjs: {
        locales: ['zh'],
        defaultLocale: 'zh'
    },
    ui: {
        icons: ['carbon']
    },
    tailwindcss: {
        safelist: [
            'grid-cols-1',
            'grid-cols-3',
        ]
    },
    vue: {
        compilerOptions: {
            isCustomElement: (tag:string) => ['meting-js'].includes(tag),
        },
    },
    app: {
        head: {
            meta: [
                { name: "viewport", content: "width=device-width, initial-scale=1, user-scalable=no" },
                { charset: "utf-8" },
            ],
            link: [
                {href: `/css/APlayer.min.css`, rel: 'stylesheet'},
            ],
            script: [
                {src: `/js/APlayer.min.js`, type: 'text/javascript', async: true, defer: true},
                {src: `/js/Meting.min.js`, type: 'text/javascript', async: true, defer: true},
                {src: `/js/main.js`, type: 'text/javascript', async: true, defer: true},
            ]
        }
    },
    vite: {
        server: {
            proxy: {
                "/api": {
                    target: "http://localhost:37892",
                    changeOrigin: true,
                    // rewrite: (path) => path.replace(/^\/api/, ""),
                },
                "/upload": {
                    target: "http://localhost:37892",
                    changeOrigin: true,
                }
            },
        },
        build: {
            rollupOptions: {
                output: {
                    hashCharacters: 'base36'
                }
            }
        }
    }
})