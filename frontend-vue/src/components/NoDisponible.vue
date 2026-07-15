<script setup>
defineProps({ nombre: { type: String, required: true } });
</script>

<template>
  <div class="nodispo">
    <div class="nodispo-icon">🚧</div>
    <h2>Sección no disponible</h2>
    <p>«{{ nombre }}» no está disponible en esta demo.</p>
  </div>
</template>
