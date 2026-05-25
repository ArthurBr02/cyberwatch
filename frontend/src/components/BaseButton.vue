<script setup lang="ts">
defineProps<{
  variant?: 'primary' | 'danger' | 'warning' | 'info' | 'secondary'
  size?: 'sm' | 'md'
  disabled?: boolean
}>()
</script>

<template>
  <button
    class="btn"
    :class="[variant ? `btn-${variant}` : '', size ? `btn-${size}` : '']"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>