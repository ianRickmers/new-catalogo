
import { ref } from 'vue'

const mostrar = ref(false)
const mensaje = ref('')
const color = ref('info')

function notify(msg: string, tipo: 'info' | 'success' | 'error' = 'info') {
  mensaje.value = msg
  color.value = tipo
  mostrar.value = true
}

export function useNotifier() {
  return { mostrar, mensaje, color, notify }
}

