<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { navigateTo } from '#app'
import useAuthStore from '../store/useAuthStore'

const authStore = useAuthStore()
const router = useRouter()

const esAdmin = computed(() => {
  return authStore.getRole?.includes('Administrador')
})

function irASolicitudes() {
  navigateTo('/solicitudes')
}

function irAAprobar() {
  navigateTo('/solicitudes/aprobar')
}

function logOut() {
  authStore.logOut()
  location.reload()
}
</script>
<template>
  <v-navigation-drawer
    id="navbar"
    width="105"
    permanent
    class="bg-grey-lighten-5"
  >

    <v-list class="navbar-buttons">
   
      <v-list-item>
        <v-tooltip location="end">
          <template #activator="{ props }">
            <v-btn icon flat v-bind="props" @click="irASolicitudes">
              <v-icon>mdi-home</v-icon>
            </v-btn>
          </template>
          <span>Todas las solicitudes</span>
        </v-tooltip>
      </v-list-item>

      <v-list-item>
        <v-tooltip location="end">
          <template #activator="{ props }">
            <v-btn icon flat v-bind="props" @click="logOut">
              <v-icon>mdi-logout</v-icon>
            </v-btn>
          </template>
          <span>Cerrar Sesión</span>
        </v-tooltip>
      </v-list-item>
  
      <v-list-item v-if="esAdmin">
        <v-tooltip location="end">
          <template #activator="{ props }">
            <v-btn icon flat v-bind="props" @click="irAAprobar">
              <v-icon>mdi-check</v-icon>
            </v-btn>
          </template>
          <span>Aprobar</span>
        </v-tooltip>
      </v-list-item>
    </v-list>

  </v-navigation-drawer>
</template>


