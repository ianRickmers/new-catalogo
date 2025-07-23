export default defineNuxtPlugin(() => {
  if (process.client) {
    globalThis.FormData = window.FormData
  }
})