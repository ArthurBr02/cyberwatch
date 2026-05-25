<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import type { Article } from '../types'
import { getArticles } from '../services/api'

const props = defineProps<{
  sourceId?: number
}>()

const articles = ref<Article[]>([])
const loading = ref(false)
const error = ref('')

const loadArticles = async () => {
  loading.value = true
  try { articles.value = await getArticles(props.sourceId) }
  catch (e) { error.value = 'Failed to load articles' }
  finally { loading.value = false }
}

watch(() => props.sourceId, loadArticles)
onMounted(loadArticles)
</script>

<template>
  <div class="article-list">
    <h2>Articles ({{ articles.length }})</h2>
    <p v-if="error" class="error">{{ error }}</p>
    <div v-if="loading" class="loading">Loading...</div>
    <table v-else class="table">
      <thead>
        <tr>
          <th>Title</th>
          <th>Date</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="a in articles" :key="a.id">
          <td>
            <a :href="a.url" target="_blank" rel="noopener" class="title-link">{{ a.title }}</a>
            <p v-if="a.summary" class="summary">{{ a.summary }}</p>
          </td>
          <td class="date-cell">{{ a.published_at ? new Date(a.published_at).toLocaleDateString() : '-' }}</td>
          <td>
            <span :class="['badge', a.is_sent ? 'sent' : 'pending']">
              {{ a.is_sent ? 'Sent' : 'Pending' }}
            </span>
          </td>
        </tr>
        <tr v-if="articles.length === 0">
          <td colspan="3" class="empty">No articles found</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.article-list { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 1.5rem; margin-top: 1.5rem; }
.article-list h2 { font-size: 1.2rem; margin-bottom: 1rem; color: #c9d1d9; }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: left; padding: 0.75rem 0.5rem; border-bottom: 1px solid #30363d; color: #8b949e; font-size: 0.8rem; text-transform: uppercase; }
.table td { padding: 0.75rem 0.5rem; border-bottom: 1px solid #21262d; }
.title-link { color: #58a6ff; text-decoration: none; font-weight: 500; }
.title-link:hover { text-decoration: underline; }
.summary { color: #8b949e; font-size: 0.8rem; margin-top: 0.2rem; }
.date-cell { white-space: nowrap; color: #8b949e; font-size: 0.85rem; }
.badge { display: inline-block; padding: 0.15rem 0.5rem; border-radius: 12px; font-size: 0.75rem; }
.badge.sent { background: #0d5320; color: #2ea043; border: 1px solid #2ea043; }
.badge.pending { background: #21262d; color: #d29922; border: 1px solid #d29922; }
.loading { text-align: center; padding: 2rem; color: #8b949e; }
.error { color: #da3633; margin-bottom: 1rem; }
.empty { text-align: center; color: #8b949e; padding: 2rem; }
</style>