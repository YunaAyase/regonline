export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxt/ui'],
  ssr: false,
  runtimeConfig: {
    public: {
      apiBase: '/api',
    },
  },
  devServer: {
    host: '::',
  },
  nitro: {
    devProxy: {
      '/api/': {
        target: 'http://localhost:8080/api/',
        changeOrigin: true,
      },
    },
  },
  css: ['~/assets/css/main.css'],
  app: {
    head: {
      title: '在线报名系统',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '在线报名系统' },
      ],
    },
  },
})
