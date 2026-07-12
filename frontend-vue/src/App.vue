<script setup>
import { onMounted, ref, computed } from 'vue';
import { getMe, logout } from './lib/api.js';
import Logo from './components/Logo.vue';
import Login from './components/Login.vue';
import DatosEmpleado from './components/DatosEmpleado.vue';
import Vacaciones from './components/Vacaciones.vue';
import Conocimientos from './components/Conocimientos.vue';
import NoDisponible from './components/NoDisponible.vue';

const SECCIONES = [
  { id: 'datos', label: 'Datos del empleado' },
  { id: 'nominas', label: 'Nóminas' },
  { id: 'retenciones', label: 'Retenciones' },
  { id: 'documentos', label: 'Otros documentos' },
  { id: 'nunenews', label: 'Nunenews' },
  { id: 'vacaciones', label: 'Vacaciones' },
  { id: 'horas', label: 'Registro de horas' },
  { id: 'ticket', label: 'Ticket Restaurant' },
  { id: 'inventario', label: 'Inventario' },
  { id: 'formacion', label: 'Formación Interna' },
  { id: 'conocimientos', label: 'Conocimientos / Certificaciones' },
  { id: 'salas', label: 'Reserva de salas' },
];

const user = ref(null);
const checking = ref(true);
const seccion = ref('datos');
const menuOpen = ref(false);

const activa = computed(() => SECCIONES.find((s) => s.id === seccion.value));

onMounted(async () => {
  try {
    user.value = await getMe();
  } catch {
    user.value = null;
  } finally {
    checking.value = false;
  }
});

async function cerrarSesion() {
  menuOpen.value = false;
  try {
    await logout();
  } finally {
    user.value = null;
    seccion.value = 'datos';
  }
}
</script>

<template>
  <template v-if="!checking">
    <Login v-if="!user" @login="user = $event" />
    <template v-else>
      <header class="header">
        <Logo class="header-logo" />
        <span class="header-title">Portal empleado</span>
        <div class="header-menu">
          <button class="header-burger" aria-label="Menú" @click="menuOpen = !menuOpen">☰</button>
          <div v-if="menuOpen" class="header-dropdown">
            <button @click="cerrarSesion">Cerrar sesión</button>
          </div>
        </div>
      </header>

      <nav class="nav">
        <button
          v-for="s in SECCIONES"
          :key="s.id"
          class="nav-item"
          :class="{ active: s.id === seccion }"
          @click="seccion = s.id"
        >
          {{ s.label }}
        </button>
      </nav>

      <main class="main">
        <DatosEmpleado v-if="seccion === 'datos'" :user="user" @update="user = $event" />
        <Vacaciones v-else-if="seccion === 'vacaciones'" />
        <Conocimientos v-else-if="seccion === 'conocimientos'" />
        <NoDisponible v-else :nombre="activa.label" />
      </main>
    </template>
  </template>
</template>
