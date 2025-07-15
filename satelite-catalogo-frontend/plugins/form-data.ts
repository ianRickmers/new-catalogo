import { defineNuxtPlugin } from '#app'
import FormData from 'form-data'

export default defineNuxtPlugin(() => {
  if (process.server && !(globalThis as any).FormData) {
    (globalThis as any).FormData = FormData as unknown as typeof globalThis.FormData
  }
})
