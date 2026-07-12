<script setup>
// Logo "n" con la elipse, recreado en SVG para la demo.
</script>

<template>
  <svg viewBox="0 0 160 110" role="img" aria-label="Nunegal Consulting">
    <defs>
      <linearGradient id="swoosh" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0%" stop-color="#1a237e" />
        <stop offset="60%" stop-color="#9fa8da" />
        <stop offset="100%" stop-color="#cfd8dc" />
      </linearGradient>
    </defs>
    <ellipse
      cx="80" cy="55" rx="70" ry="34"
      fill="none" stroke="url(#swoosh)" stroke-width="7"
      transform="rotate(-14 80 55)"
    />
    <text
      x="80" y="78" text-anchor="middle"
      font-size="72" font-weight="700" font-family="Georgia, serif"
      fill="#1a237e"
    >
      n
    </text>
  </svg>
</template>
