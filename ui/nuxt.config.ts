// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-04-03",
  devtools: { enabled: true },
  app: {
    head: {
      title: 'Lifery'
    }
  },
  modules: ["@vueuse/nuxt", "@vee-validate/nuxt", "@nuxt/ui", '@nuxtjs/i18n'],
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_API_BASE_URL?.trim() || "http://localhost:8080",
      cloudinaryCloudName: process.env.NUXT_PUBLIC_CLOUDINARY_CLOUD_NAME,
      cloudinaryUploadPreset: process.env.NUXT_PUBLIC_CLOUDINARY_UPLOAD_PRESET,
    }
  },
  devServer: {
    port: 8081,
  },
  i18n: {
    strategy: 'prefix_except_default',
    defaultLocale: 'tr',
    lazy: true,
    langDir: './locales',
    locales: [
      {
        code: 'tr',
        name: 'Türkçe',
        file: 'tr.json'
      },
      {
        code: 'en',
        name: 'English',
        file: 'en.json'
      }
    ]
  }
});
