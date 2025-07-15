<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import useAuthStore from '../store/useAuthStore'
import { useNotifier } from '@/composables/useNotifier'

const { notify } = useNotifier()

definePageMeta({
  auth: false,
})
const emailError = ref('')
const passwordError = ref('')
const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const showPassword = ref(false)

const emailRules = [
  value => {
    if (/[^0-9]/.test(value)) return true
    return 'Ingresa un correo institucional'
  },
]

const passwordRules = [
  value => {
    if (value?.length >= 3) return true
    return 'Ingresa una contraseña válida'
  },
]

const handleSubmit = async () => {
  // Limpiar errores anteriores
  emailError.value = ''
  passwordError.value = ''

  if (!email.value || !password.value) {
    notify('Por favor completa todos los campos.', 'warning')
    return
  }

  // Validar reglas del email
  const emailValidation = emailRules.map(rule => rule(email.value)).find(r => r !== true)
  if (emailValidation) {
    emailError.value = emailValidation
    return
  }

  // Validar reglas del password
  const passwordValidation = passwordRules.map(rule => rule(password.value)).find(r => r !== true)
  if (passwordValidation) {
    passwordError.value = passwordValidation
    return
  }

  try {
    await authStore.logIn({
      user: email.value,
      password: password.value,
    })

    notify('Inicio de sesión exitoso!', 'success')

    // const allowedRoles = ['Usuario', 'Administrador']
    // const roles = authStore.getRole || []

    // if (roles.some(r => allowedRoles.includes(r.trim()))) {
    //   router.push('/solicitudes')
    // }


    router.push('/solicitudes')
  } catch (error) {
    console.error(error)

    // Mostrar error directamente en los campos
    emailError.value = 'Correo o contraseña incorrectos'
    passwordError.value = 'Correo o contraseña incorrectos'

    notify('Credenciales inválidas. Intente nuevamente.', 'error')
  }
}


</script>

<template>
  <div class="login-container">
    <v-sheet class="mx-auto v-sheet-login" >
      <h2 class="login-form-title">iniciar sesión</h2>
      <v-form fast-fail @submit.prevent="handleSubmit">
        <v-text-field 
          class="login-input"
          v-model="email"
          :rules="emailRules"
          label="correo institucional"
          placeholder="correo institucional (sin @usach.cl)"
          variant="solo-filled"
          prepend-inner-icon="mdi-account"
          :error-messages="emailError ? [emailError] : []"
        />
        <v-text-field 
          class="login-input"
          v-model="password"
          :rules="passwordRules"
          :type="showPassword ? 'text' : 'password'"
          :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
          @click:append-inner="showPassword = !showPassword"
          label="contraseña"
          variant="solo-filled"
          prepend-inner-icon="mdi-lock"
          :error-messages="passwordError ? [passwordError] : []"
        />
        <v-btn class="mt-2 login-btn" type="submit" block>Ingresar</v-btn>
      </v-form>
    </v-sheet>
  </div>
</template>
