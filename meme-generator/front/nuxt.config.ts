import { aliases } from 'vuetify/iconsets/mdi'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  devtools: { enabled: true },
  modules: [
    '@pinia/nuxt',
    '@vueuse/nuxt',
    'vuetify-nuxt-module',
    '@nuxt/icon',
    '@nuxt/eslint',
  ],
  css: ['~/assets/styles/index.css'],
  experimental: { typedPages: true },
  typescript: { shim: false, strict: true },
  vue: { propsDestructure: true },
  vueuse: { ssrHandlers: true },
  icon: {
    clientBundle: {
      icons: Object.values(aliases).map((v) =>
        (v as string).replace(/^mdi-/, 'mdi:'),
      ),
      scan: true,
      // scan all components in the project and include icons
      // scan: true,
    },
    customCollections: [
      {
        prefix: 'custom',
        dir: './assets/icons',
      },
    ],
  },
  vite: {
    build: { sourcemap: false },
  },
  runtimeConfig: {
    public: {
      apiBaseUrl: '/api',
    },
  },
  compatibilityDate: '2024-08-05',
})
