<script setup>
import { computed, onMounted, reactive, ref } from 'vue';
import {
  listCertificaciones, createCertificacion, updateCertificacion, deleteCertificacion,
} from '../lib/api.js';

const COLUMNAS = [
  { id: 'conocimiento', label: 'Conocimiento / Certificación', sortable: true, izq: true },
  { id: 'empresaEmisora', label: 'Empresa emisora', sortable: true },
  { id: 'fecha', label: 'Fecha', sortable: true },
];

const items = ref([]);
const cargando = ref(true);
const error = ref('');

const buscar = ref('');
const pageSize = ref(10);
const page = ref(1);
const sort = reactive({ campo: 'conocimiento', dir: 'asc' });

const editando = ref(null); // null | {id?, conocimiento, ...}
const borrando = ref(null); // null | item

// Modal de formulario
const form = reactive({ conocimiento: '', empresaEmisora: '', fecha: '' });
const formBusy = ref(false);
const formError = ref('');

async function recargar() {
  cargando.value = true;
  error.value = '';
  try {
    items.value = await listCertificaciones();
  } catch {
    error.value = 'No se pudieron cargar las certificaciones';
  } finally {
    cargando.value = false;
  }
}

onMounted(recargar);

const filtradas = computed(() => {
  const q = buscar.value.trim().toLowerCase();
  let arr = items.value;
  if (q) {
    arr = arr.filter((c) =>
      [c.conocimiento, c.empresaEmisora, c.fecha]
        .some((v) => (v ?? '').toLowerCase().includes(q)));
  }
  arr = [...arr].sort((a, b) => {
    const r = String(a[sort.campo] ?? '').localeCompare(String(b[sort.campo] ?? ''), 'es', { numeric: true });
    return sort.dir === 'asc' ? r : -r;
  });
  return arr;
});

const total = computed(() => filtradas.value.length);
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)));
const pageActual = computed(() => Math.min(page.value, totalPages.value));
const inicio = computed(() => (total.value === 0 ? 0 : (pageActual.value - 1) * pageSize.value));
const visibles = computed(() => filtradas.value.slice(inicio.value, inicio.value + pageSize.value));
const paginas = computed(() => Array.from({ length: totalPages.value }, (_, i) => i + 1));

function cambiarSort(col) {
  if (!col.sortable) return;
  if (sort.campo === col.id) {
    sort.dir = sort.dir === 'asc' ? 'desc' : 'asc';
  } else {
    sort.campo = col.id;
    sort.dir = 'asc';
  }
}

function abrirNuevo() {
  editando.value = {};
  Object.assign(form, { conocimiento: '', empresaEmisora: '', fecha: '' });
  formError.value = '';
}

function abrirEditar(c) {
  editando.value = c;
  Object.assign(form, {
    conocimiento: c.conocimiento ?? '',
    empresaEmisora: c.empresaEmisora ?? '',
    fecha: c.fecha ?? '',
  });
  formError.value = '';
}

async function guardar() {
  if (!form.conocimiento.trim()) {
    formError.value = 'El conocimiento / certificación es obligatorio';
    return;
  }
  formBusy.value = true;
  formError.value = '';
  try {
    const datos = { ...form };
    if (editando.value.id) {
      await updateCertificacion(editando.value.id, datos);
    } else {
      await createCertificacion(datos);
    }
    editando.value = null;
    await recargar();
  } catch {
    formError.value = 'No se pudo guardar';
  } finally {
    formBusy.value = false;
  }
}

async function confirmarBorrado() {
  await deleteCertificacion(borrando.value.id);
  borrando.value = null;
  await recargar();
}
</script>

<template>
  <div>
    <div class="cert-top">
      <div class="cert-mostrar">
        Mostrar
        <select v-model.number="pageSize" @change="page = 1">
          <option v-for="n in [10, 25, 50]" :key="n" :value="n">{{ n }}</option>
        </select>
        registros
      </div>
      <label class="cert-buscar">
        Buscar:
        <input v-model="buscar" @input="page = 1" />
      </label>
    </div>

    <table class="cert-tabla">
      <thead>
        <tr>
          <th
            v-for="col in COLUMNAS"
            :key="col.id"
            :class="[{ sortable: col.sortable }, col.izq ? 'cert-col-conocimiento' : '']"
            @click="cambiarSort(col)"
          >
            {{ col.label }}
            <span v-if="col.sortable" class="cert-sort" :class="{ activo: sort.campo === col.id }">
              {{ sort.campo === col.id ? (sort.dir === 'asc' ? '▲' : '▼') : '⇅' }}
            </span>
          </th>
          <th>Opciones</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="cargando"><td class="cert-vacio" colspan="4">Cargando…</td></tr>
        <tr v-else-if="error"><td class="cert-vacio" colspan="4">{{ error }}</td></tr>
        <tr v-else-if="visibles.length === 0"><td class="cert-vacio" colspan="4">No hay registros disponibles</td></tr>
        <tr v-else v-for="c in visibles" :key="c.id">
          <td class="cert-col-conocimiento">{{ c.conocimiento }}</td>
          <td>{{ c.empresaEmisora }}</td>
          <td>{{ c.fecha }}</td>
          <td>
            <div class="cert-acciones">
              <button class="cert-icono" title="Editar" @click="abrirEditar(c)">✏️</button>
              <button class="cert-icono borrar" title="Eliminar" @click="borrando = c">🗑️</button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <div class="cert-footer">
      <div class="cert-info">
        <template v-if="total === 0">Mostrando 0 de 0 registros</template>
        <template v-else>
          Mostrando {{ inicio + 1 }} a {{ Math.min(inicio + pageSize, total) }} de un total de {{ total }} registros
        </template>
      </div>
      <div class="cert-paginacion">
        <button class="cert-pag" :disabled="pageActual <= 1" @click="page = pageActual - 1">‹</button>
        <button
          v-for="n in paginas"
          :key="n"
          class="cert-pag"
          :class="{ activo: n === pageActual }"
          @click="page = n"
        >
          {{ n }}
        </button>
        <button class="cert-pag" :disabled="pageActual >= totalPages" @click="page = pageActual + 1">›</button>
      </div>
    </div>

    <button class="cert-anadir" @click="abrirNuevo">Añadir</button>

    <!-- Modal alta/edición -->
    <div v-if="editando" class="modal-overlay" @click="editando = null">
      <form class="modal" @click.stop @submit.prevent="guardar">
        <div class="modal-head">{{ editando.id ? 'Editar certificación' : 'Nueva certificación' }}</div>
        <div class="modal-body">
          <div class="modal-campo">
            <label for="m-conocimiento">Conocimiento / Certificación</label>
            <input id="m-conocimiento" v-model="form.conocimiento" autofocus />
          </div>
          <div class="modal-campo">
            <label for="m-empresa">Empresa emisora</label>
            <input id="m-empresa" v-model="form.empresaEmisora" />
          </div>
          <div class="modal-campo">
            <label for="m-fecha">Fecha</label>
            <input id="m-fecha" type="date" v-model="form.fecha" />
          </div>
          <span v-if="formError" class="perfil-error">{{ formError }}</span>
        </div>
        <div class="modal-foot">
          <button type="button" class="btn" @click="editando = null">Cancelar</button>
          <button type="submit" class="btn btn-primary" :disabled="formBusy">Guardar</button>
        </div>
      </form>
    </div>

    <!-- Modal borrado -->
    <div v-if="borrando" class="modal-overlay" @click="borrando = null">
      <div class="modal" @click.stop>
        <div class="modal-head">Eliminar certificación</div>
        <div class="modal-texto">¿Seguro que quieres eliminar «{{ borrando.conocimiento }}»?</div>
        <div class="modal-foot">
          <button class="btn" @click="borrando = null">Cancelar</button>
          <button class="btn btn-danger" @click="confirmarBorrado">Eliminar</button>
        </div>
      </div>
    </div>
  </div>
</template>
