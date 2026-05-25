<template>
  <div id="app">
    <header class="app-header">
      <h1>🛡️ CyberWatch</h1>
      <nav>
        <router-link to="/">Dashboard</router-link>
      </nav>
    </header>
    <main class="app-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0d1117; color: #c9d1d9; }
.app-header { background: #161b22; padding: 1rem 2rem; display: flex; align-items: center; gap: 2rem; border-bottom: 1px solid #30363d; }
.app-header h1 { font-size: 1.5rem; color: #58a6ff; }
.app-header nav a { color: #8b949e; text-decoration: none; font-size: 0.9rem; }
.app-header nav a:hover { color: #58a6ff; }
.app-main { padding: 2rem; max-width: 1200px; margin: 0 auto; }
</style>