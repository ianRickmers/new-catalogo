<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue';
import axios from 'axios';
// Variables reactivas
const authStore = useAuthStore();
const items = ref([]);
const currentPage = ref(1);
// const pageSize = ref(10);
const totalRecords = ref(0);
const isMobile = ref(false);
const productoSeleccionado = ref(null);
const searchQuery = reactive({
  id_product: '',
  categoria: '',
  descripcion: ''
});
const allProducts = ref([]);
const loading = ref(false)
const error = ref('')
const page = ref(1)
const pageSize = ref(10)
const totalPages = ref(0)
// Props y eventos
const props = defineProps({
  modelValue: Boolean
});

const emit = defineEmits(['update:modelValue', 'cerrar']);

const mostrar = ref(props.modelValue);

function seleccionarProducto(producto) {
  productoSeleccionado.value = producto;

  emit('update:modelValue', false); // Cierra el modal
  emit('cerrar', producto); // Envía el producto al padre
  resetearFiltros();
}

// Sincronizar cambios del padre al modal
watch(() => props.modelValue, (val) => {
  mostrar.value = val;
});
// Observar cambio de página: si hay filtros, usar filtrados; si no, traer todos
watch(page, () => {
  const hayFiltros = searchQuery.id_product || searchQuery.categoria || searchQuery.descripcion;
  if (hayFiltros) {
    getProductosFiltrados();
  } else {
    getProductos();
  }
});

// Comunicar cambios del modal al padre
watch(mostrar, (val) => {
  emit('update:modelValue', val);
});

// Cerrar modal
function cerrar() {
  mostrar.value = false;
  resetearFiltros();
  productoSeleccionado.value = null; // Limpiar selección
  emit('cerrar',productoSeleccionado.value); // Se pasa el producto seleccionado al padre vacio
}

// Comprobar si es dispositivo móvil
const checkIfMobile = () => {
  isMobile.value = window.innerWidth <= 1023;
};

// Actualizar página
const updatePage = (page: number) => {
  currentPage.value = page;
};
async function getProductos() {
  loading.value = true
  error.value = ''
  try {
    const res = await axios.get(`${useRuntimeConfig().public.baseURL}/product/paginated`, {
      params: {
        page: page.value,
        pageSize: pageSize.value
      }
    })
    allProducts.value = res.data.data
    totalPages.value = res.data.totalPages || 1
  } catch (err) {
    error.value = 'Error al cargar los productos'
    allProducts.value = []
    totalPages.value = 0
    console.error(err)
  } finally {
    loading.value = false
  }
}
// se obtienen los productos filtrados
async function getProductosFiltrados() {
  loading.value = true
  error.value = ''

  const params = {
    pageSize: pageSize.value,
    id_product: searchQuery.id_product,
    categoria: searchQuery.categoria,
    descripcion: searchQuery.descripcion
  }

  try {
    const res = await axios.get(
      `${useRuntimeConfig().public.baseURL}/product/filtradas`,
      {
        params,
        headers: {
          'Authorization': `Bearer ${authStore.getToken}`
        }
      }
    )

    allProducts.value = res.data.data
    totalPages.value = res.data.totalPages || 1
  } catch (err) {
    error.value = 'Error al filtrar solicitudes'
    allProducts.value = []
    totalPages.value = 0
    console.error(err)
  } finally {
    loading.value = false
  }
}



async function resetearFiltros() {
  searchQuery.id_product = '';
  searchQuery.categoria = '';
  searchQuery.descripcion = '';
  await getProductos();
}
onMounted(() => {
  checkIfMobile();
  window.addEventListener('resize', checkIfMobile);
  // getData();
  getProductos();

});

onUnmounted(() => {
  window.removeEventListener('resize', checkIfMobile);
});
</script>

<template>
  <!-- Modal controlado por Vuetify con v-model -->
  <v-dialog :model-value="mostrar"  @update:modelValue="$emit('update:modelValue', $event)" max-width="1200"   persistent>
    <template #default>
      <div class="modal-content">
        <div v-if="!isMobile" class="content-container">
          <div class="search-container">
            <div class="search-bar">
              <input
                type="text"
                class="search-input"
                placeholder="Buscar por ID artículo"
                v-model="searchQuery.id_product"
              />
              <div class="divider"></div>
              <input
                type="text"
                class="search-input"
                placeholder="Buscar por categoria"
                v-model="searchQuery.categoria"
              />
              <div class="divider"></div>
              <input
                type="text"
                class="search-input"
                placeholder="Buscar por Descripción"
                v-model="searchQuery.descripcion"
              />
            </div>
          </div>

          <div class="botones centrado">
            <button class="boton-outlined" @click="getProductosFiltrados">Consultar</button>
            <button class="boton-outlined" @click="resetearFiltros">Borrar</button>
            <button class="boton-outlined" @click="cerrar">Cancelar</button>
          </div>

          <table class="tabla-planillas">
            <thead>
              <tr class="table-header">
                <th scope="col">ID Artículo</th>
                <th scope="col">Categoria</th>
                <th scope="col">ID Categoria</th>
                <th scope="col">Descripción</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(producto, index) in allProducts"
                :key="index"
                :class="{ 'row-even': index % 2 === 0, 'row-odd': index % 2 !== 0 }"
                @click="seleccionarProducto(producto)"
                style="cursor: pointer;"
              >
                <td>{{ producto.id_product }}</td>
                <td>{{ producto.categoria }}</td>
                <td>{{ producto.id_categoria }}</td>
                <td>{{ producto.descripcion}}</td>
              </tr>
            </tbody>
          </table>

          <!-- PAGINACION -->
          
          <div class="pagination-container">
            <v-pagination
              v-if="totalPages > 1"
              v-model="page"
              :length="totalPages"
              class="mt-4"
            />
          </div>
            <v-alert v-if="error" type="error" class="mt-4">{{ error }}</v-alert>
            <v-progress-linear v-if="loading" indeterminate color="teal" class="mt-2" />         
        </div>
      </div>
    </template>
  </v-dialog>
</template>


