<script setup>
import { ref } from 'vue';
import { login } from '../lib/api.js';
import Logo from './Logo.vue';

const emit = defineEmits(['login']);

const username = ref('');
const password = ref('');
const error = ref('');
const busy = ref(false);

async function submit() {
  busy.value = true;
  error.value = '';
  try {
    const user = await login(username.value, password.value);
    emit('login', user);
  } catch (err) {
    error.value = err.status === 401
      ? 'Usuario o contraseña no válidos'
      : 'No se pudo conectar con el servidor';
  } finally {
    busy.value = false;
  }
}
</script>

<template>
  <div class="login-page">
    <Logo class="login-logo" />
    <h1 class="login-title">Acceso al<br />Portal del Empleado</h1>
    <form class="login-form" @submit.prevent="submit">
      <input v-model="username" placeholder="Usuario" autofocus />
      <input v-model="password" type="password" placeholder="Contraseña" />
      <button type="submit" :disabled="busy">Entrar</button>
      <p class="login-error">{{ error }}</p>
    </form>
    <div class="login-version">PEN version 2.9.5</div>
    <div class="login-copy">Nunegal Consulting ©</div>
  </div>
</template>
