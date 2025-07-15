<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  // si el archivo se llamara login.vue tendria que ser /login/login
  // ya que dejaria de ser el archivo principal de esta car
  router.push('/login')
})
</script>
