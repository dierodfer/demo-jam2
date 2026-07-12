<script setup>
import { reactive, ref } from 'vue';
import { updateMe } from '../lib/api.js';

const props = defineProps({ user: { type: Object, required: true } });
const emit = defineEmits(['update']);

const CAMPOS = [
  { id: 'nombre', label: 'Nombre' },
  { id: 'email', label: 'Email' },
  { id: 'telefono', label: 'Teléfono' },
  { id: 'puesto', label: 'Puesto' },
  { id: 'departamento', label: 'Departamento' },
  { id: 'direccion', label: 'Dirección' },
  { id: 'foto', label: 'Foto (URL)' },
];

const form = reactive({
  nombre: props.user.nombre ?? '',
  email: props.user.email ?? '',
  telefono: props.user.telefono ?? '',
  puesto: props.user.puesto ?? '',
  departamento: props.user.departamento ?? '',
  direccion: props.user.direccion ?? '',
  foto: props.user.foto ?? '',
});

const estado = ref(null); // null | 'ok' | 'error'
const busy = ref(false);

async function guardar() {
  busy.value = true;
  try {
    const actualizado = await updateMe({ ...form });
    emit('update', actualizado);
    estado.value = 'ok';
  } catch {
    estado.value = 'error';
  } finally {
    busy.value = false;
  }
}
</script>

<template>
  <div class="perfil">
    <div class="perfil-foto-col">
      <img
        class="perfil-foto"
        :src="form.foto || 'data:image/svg+xml,<svg xmlns=&quot;http://www.w3.org/2000/svg&quot;/>'"
        :alt="form.nombre || 'Foto de perfil'"
      />
    </div>
    <form class="perfil-form" @submit.prevent="guardar">
      <div v-for="c in CAMPOS" :key="c.id" class="perfil-campo">
        <label :for="`campo-${c.id}`">{{ c.label }}</label>
        <input :id="`campo-${c.id}`" v-model="form[c.id]" @input="estado = null" />
      </div>
      <div class="perfil-acciones">
        <button type="submit" :disabled="busy">Guardar</button>
        <span v-if="estado === 'ok'" class="perfil-ok">Cambios guardados ✓</span>
        <span v-if="estado === 'error'" class="perfil-error">No se pudo guardar</span>
      </div>
    </form>
  </div>
</template>
