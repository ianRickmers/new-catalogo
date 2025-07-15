<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

function irASolicitudes() {
  if (router.currentRoute.value.path === '/solicitudes') {
    window.location.reload()
  } else {
    router.push('/solicitudes')
  }
}
</script>

<template>
  <v-app-bar app>
    <img src="@/assets/img/Usach_P1.png" alt="Logo USACH" class="header-logo" @click="irASolicitudes"/>
    <div class="vertical-line"></div>
    <h2 class="header-title" @click="irASolicitudes">Catálogo</h2>
  </v-app-bar>
</template>
