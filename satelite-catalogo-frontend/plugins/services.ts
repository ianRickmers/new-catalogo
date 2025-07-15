// plugins/fetchModule.ts
import { defineNuxtPlugin } from 'nuxt/app'
import { Fetch } from '../common/fetchModule'

export default defineNuxtPlugin(() => {
  const fetchModule = new Fetch()
  return {
    provide: {
      fetchModule,
    },
  }
})