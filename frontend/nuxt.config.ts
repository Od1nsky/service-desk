import { fileURLToPath } from 'node:url'
import { dirname, resolve } from 'node:path'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)
const scssDir = resolve(__dirname, 'assets/scss')

export default defineNuxtConfig({
  ssr: false,
  modules: ['@pinia/nuxt'],
  css: ['~/assets/scss/main.scss'],

  app: {
    head: {
      title: 'ServiceDesk',
      meta: [{ name: 'viewport', content: 'width=device-width, initial-scale=1' }],
      link: [{ rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    },
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',
    },
  },

  vite: {
    css: {
      preprocessorOptions: {
        scss: {
          // @import даёт плоский скоуп: переменные из variables доступны в mixins
          additionalData: `
            @use "sass:color";
            @use "${scssDir}/_variables.scss" as *;
            @use "${scssDir}/_mixins.scss" as *;
          `,
        },
      },
    },
  },
})
