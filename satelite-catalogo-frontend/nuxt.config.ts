import vuetify from 'vite-plugin-vuetify'
import { defineNuxtConfig } from "nuxt/config";

export default defineNuxtConfig({
  modules: [
    '@pinia/nuxt',
    '@sidebase/nuxt-session',
  ],

  typescript: {
    strict: true,
  },

  app: {
    head: {
      title: 'Catalogo USACH',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Template de prueba para postulantes' },
      ],
    },
  },

  css: [
    'vuetify/styles',
    '@mdi/font/css/materialdesignicons.min.css', // iconos opcionales
    '@/assets/scss/main.scss',
  ],

  plugins: [
    '~/plugins/form-data',
  ],

  build: {
    transpile: ['vuetify'],
  },

  vite: {
    plugins: [vuetify()],
  },

  runtimeConfig: {
    public: {
      backBaseUrl: process.env.NUXT_PUBLIC_BACK_BASE_URL,
      baseURL: process.env.NUXT_PUBLIC_BACK_BASE_URL,
    },
  },

  imports: {
    dirs: ['./store'],
  },


 	session: {
		session: {
			expiryInSeconds: 3600,
			cookieSameSite: 'lax',
			cookieSecure: true,
			cookieHttpOnly: true,
			storageOptions: {
				driver: 'memory',
				options: {},
			},
			domain: false,
		},
	},
})
