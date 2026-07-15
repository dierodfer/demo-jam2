<script setup>
import { ref } from 'vue';
import {
  MESES, DIAS_SEMANA, DEMO, buildMonth, dateKey, dayFlags, todayKey,
} from '../lib/vacaciones.js';

const year = ref(DEMO.year);
const toast = ref(false);
const hoy = todayKey();
const resumen = DEMO.resumen;

function clasesDia(monthIndex, day) {
  const key = dateKey(year.value, monthIndex, day);
  const f = dayFlags(key, year.value);
  return {
    'cal-dia': true,
    festivo: f.holiday,
    usado: f.used,
    marcado: f.selected,
    hoy: key === hoy,
  };
}

function solicitar() {
  toast.value = true;
  setTimeout(() => { toast.value = false; }, 2500);
}
</script>

<template>
  <div>
    <div class="vac-yearbar">
      <button class="vac-arrow" aria-label="Año anterior" @click="year--">‹</button>
      <div class="vac-years">
        <button class="vac-year-side" @click="year--">{{ year - 1 }}</button>
        <span class="vac-year">{{ year }}</span>
        <button class="vac-year-side" @click="year++">{{ year + 1 }}</button>
      </div>
      <button class="vac-arrow" aria-label="Año siguiente" @click="year++">›</button>
    </div>

    <div class="vac-grid">
      <div v-for="(mes, m) in MESES" :key="m" class="cal-mes">
        <div class="cal-mes-nombre">{{ mes }}</div>
        <div class="cal-tabla">
          <div v-for="d in DIAS_SEMANA" :key="d" class="cal-cab">{{ d }}</div>
          <template v-for="(day, i) in buildMonth(year, m)" :key="i">
            <div v-if="day === null" />
            <div v-else :class="clasesDia(m, day)">{{ day }}</div>
          </template>
        </div>
      </div>
    </div>

    <div class="vac-resumen">
      <div class="vac-fila">
        <div class="vac-seg amarillo">Acumulados año anterior no disfrutados: {{ resumen.acumulados }}</div>
        <div class="vac-seg azul">Naturales: {{ resumen.naturales }} días</div>
      </div>
      <div class="vac-fila">
        <div class="vac-seg rojo">Utilizados: {{ resumen.utilizados }} días</div>
        <div class="vac-seg verde">Disponibles: {{ resumen.disponibles }} días</div>
      </div>
    </div>

    <div class="vac-leyenda">
      <span><i class="vac-punto amarillo" /> Acumulados</span>
      <span><i class="vac-punto azul" /> Totales</span>
      <span><i class="vac-punto rojo" /> Utilizados</span>
      <span><i class="vac-punto verde" /> Disponibles</span>
    </div>

    <div class="vac-solicitar">
      <button @click="solicitar">Solicitar</button>
      <div v-if="toast" class="vac-toast">La solicitud de vacaciones no está disponible en esta demo.</div>
    </div>
  </div>
</template>
