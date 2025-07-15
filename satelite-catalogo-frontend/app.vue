<template>
  <div class="layout-container"> 
    <v-app>
      <Header />
      <v-main>
        <NuxtPage />
        <v-snackbar
          v-model="mostrar"
          :color="color"
          location="bottom right"
          timeout="3000"
          elevation="10"
          rounded
        >
          {{ mensaje }}
        </v-snackbar>
      </v-main>
      <Footer />
    </v-app>
  </div> 
</template>

<script setup>
import { useNotifier } from '@/composables/useNotifier'

const { mostrar, mensaje, color } = useNotifier()
</script>
