<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Source } from '../types'
import { getSources, deleteSource, generateRules, testScrape, scrapeAll } from '../services/api'

const emit = defineEmits<{
  (e: 'edit', source: Source): void
  (e: 'create'): void
}>()

const sources = ref<Source[]>([])
const loading = ref(false)
const error = ref('')
const scrapping = ref<number | null>(null)
const generatingRules = ref<number | null>(null)

const loadSources = async () => {
  loading.value = true
  try { sources.value = await getSources() }
  catch (e) { error.value = 'Failed to load sources' }
  finally { loading.value = false }
}

const handleDelete = async (id: number) => {
  if (!confirm('Delete this source?')) return
  try { await deleteSource(id); await loadSources() }
  catch (e) { error.value = 'Failed to delete source' }
}

const handleGenerateRules = async (id: number) => {
  generatingRules.value = id
  try { 
    await generateRules(id)
    await loadSources() 
    alert('Rules generated successfully!')
  }
  catch (e) { error.value = 'Failed to generate rules' }
  finally { generatingRules.value = null }
}

const handleTestScrape = async (id: number) => {
  scrapping.value = id
  try {
    const articles = await testScrape(id)
    alert(`Found ${articles.length} articles`)
  } catch (e) { error.value = 'Failed to test scrape' }
  finally { scrapping.value = null }
}

const handleScrapeAll = async () => {
  scrapping.value = -1
  try { await scrapeAll(); alert('Scraping completed!') }
  catch (e) { error.value = 'Failed to scrape all' }
  finally { scrapping.value = null }
}

onMounted(loadSources)
</script>

<template>
  <div class="source-list">
    <div class="toolbar">
      <h2>Sources ({{ sources.length }})</h2>
      <div class="actions">
        <button class="btn btn-secondary" @click="$emit('create')">+ New Source</button>
        <button class="btn btn-primary" :disabled="scrapping === -1" @click="handleScrapeAll">
          {{ scrapping === -1 ? 'Scraping...' : 'Scrape All' }}
        </button>
      </div>
    </div>
    <p v-if="error" class="error">{{ error }}</p>
    <div v-if="loading" class="loading">Loading...</div>
    <table v-else class="table">
      <thead>
        <tr>
          <th>Name</th>
          <th>URL</th>
          <th>Type</th>
          <th>Status</th>
          <th>Rules</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in sources" :key="s.id">
          <td><strong>{{ s.name }}</strong></td>
          <td class="url-cell">{{ s.url }}</td>
          <td><span class="badge">{{ s.fetch_type }}</span></td>
          <td>
            <span :class="['badge', s.is_active ? 'active' : 'inactive']">
              {{ s.is_active ? 'Active' : 'Inactive' }}
            </span>
          </td>
          <td class="rules-cell">
            <code v-if="s.llm_rules">{{ s.llm_rules }}</code>
            <span v-else class="text-muted">No rules</span>
          </td>
          <td class="actions-cell">
            <button class="btn btn-sm" @click="$emit('edit', s)">Edit</button>
            <button class="btn btn-sm btn-warning" :disabled="generatingRules === s.id" @click="handleGenerateRules(s.id)">
              {{ generatingRules === s.id ? '...' : 'Rules' }}
            </button>
            <button class="btn btn-sm btn-info" :disabled="scrapping === s.id" @click="handleTestScrape(s.id)">
              {{ scrapping === s.id ? '...' : 'Test' }}
            </button>
            <button class="btn btn-sm btn-danger" @click="handleDelete(s.id)">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.source-list { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 1.5rem; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
.toolbar h2 { font-size: 1.2rem; color: #c9d1d9; }
.actions { display: flex; gap: 0.5rem; }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: left; padding: 0.75rem 0.5rem; border-bottom: 1px solid #30363d; color: #8b949e; font-size: 0.8rem; text-transform: uppercase; }
.table td { padding: 0.75rem 0.5rem; border-bottom: 1px solid #21262d; }
.url-cell { max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #58a6ff; }
.rules-cell { max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.rules-cell code { font-size: 0.7rem; background: #0d1117; padding: 0.2rem; border-radius: 4px; color: #79c0ff; }
.text-muted { color: #8b949e; font-size: 0.75rem; }
.actions-cell { display: flex; gap: 0.3rem; flex-wrap: wrap; }
.btn { padding: 0.4rem 0.8rem; border: 1px solid #30363d; border-radius: 6px; cursor: pointer; font-size: 0.8rem; background: #21262d; color: #c9d1d9; }
.btn:hover { background: #30363d; }
.btn-primary { background: #238636; border-color: #2ea043; color: #fff; }
.btn-primary:hover { background: #2ea043; }
.btn-danger { border-color: #da3633; color: #da3633; }
.btn-danger:hover { background: #da3633; color: #fff; }
.btn-warning { border-color: #d29922; color: #d29922; }
.btn-warning:hover { background: #d29922; color: #fff; }
.btn-info { border-color: #58a6ff; color: #58a6ff; }
.btn-info:hover { background: #58a6ff; color: #fff; }
.btn-sm { padding: 0.2rem 0.5rem; font-size: 0.75rem; }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.badge { display: inline-block; padding: 0.15rem 0.5rem; border-radius: 12px; font-size: 0.75rem; background: #21262d; border: 1px solid #30363d; }
.badge.active { border-color: #2ea043; color: #2ea043; }
.badge.inactive { border-color: #8b949e; color: #8b949e; }
.loading { text-align: center; padding: 2rem; color: #8b949e; }
.error { color: #da3633; margin-bottom: 1rem; }
</style>
