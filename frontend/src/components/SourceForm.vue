<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Source } from '../types'
import { createSource, updateSource } from '../services/api'

const props = defineProps<{
  source?: Source | null
}>()

const emit = defineEmits<{
  (e: 'saved'): void
  (e: 'cancel'): void
}>()

const name = ref('')
const url = ref('')
const fetchType = ref<'html' | 'json'>('html')
const llmRules = ref('')
const isActive = ref(true)
const submitting = ref(false)
const error = ref('')

const isEditing = () => props.source !== null && props.source !== undefined

watch(() => props.source, (s) => {
  if (s) { 
    name.value = s.name
    url.value = s.url
    fetchType.value = s.fetch_type as 'html' | 'json'
    llmRules.value = s.llm_rules || ''
    isActive.value = s.is_active 
  }
  else { 
    name.value = ''
    url.value = ''
    fetchType.value = 'html'
    llmRules.value = ''
    isActive.value = true 
  }
}, { immediate: true })

const handleSubmit = async () => {
  if (!name.value.trim() || !url.value.trim()) { error.value = 'Name and URL are required'; return }
  submitting.value = true; error.value = ''
  try {
    const payload = { 
      name: name.value, 
      url: url.value, 
      fetch_type: fetchType.value, 
      llm_rules: llmRules.value,
      is_active: isActive.value 
    }
    if (isEditing() && props.source) {
      await updateSource(props.source.id, payload as any)
    } else {
      await createSource(payload as any)
    }
    emit('saved')
  } catch (e) { error.value = 'Failed to save source' }
  finally { submitting.value = false }
}
</script>

<template>
  <div class="form-overlay" @click.self="$emit('cancel')">
    <div class="form-card">
      <h2>{{ isEditing() ? 'Edit Source' : 'New Source' }}</h2>
      <form @submit.prevent="handleSubmit">
        <div class="field">
          <label>Name</label>
          <input v-model="name" type="text" placeholder="e.g. The Hacker News" required />
        </div>
        <div class="field">
          <label>URL</label>
          <input v-model="url" type="url" placeholder="https://example.com/rss" required />
        </div>
        <div class="field">
          <label>Fetch Type</label>
          <select v-model="fetchType">
            <option value="html">HTML</option>
            <option value="json">JSON</option>
          </select>
        </div>
        <div class="field">
          <label>LLM Rules (JSON)</label>
          <textarea v-model="llmRules" rows="5" placeholder='{"container": ".post", "title": "h2", "url": "a", "summary": "p"}'></textarea>
          <small class="hint">Leave empty to use "Generate Rules" button, or edit manually here.</small>
        </div>
        <div v-if="isEditing()" class="field">
          <label class="checkbox">
            <input v-model="isActive" type="checkbox" /> Active
          </label>
        </div>
        <p v-if="error" class="error">{{ error }}</p>
        <div class="form-actions">
          <button type="button" class="btn" @click="$emit('cancel')">Cancel</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">
            {{ submitting ? 'Saving...' : (isEditing() ? 'Update' : 'Create') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.form-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.7); display: flex; align-items: center; justify-content: center; z-index: 100; }
.form-card { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 2rem; width: 100%; max-width: 480px; max-height: 90vh; overflow-y: auto; }
.form-card h2 { margin-bottom: 1.5rem; color: #c9d1d9; }
.field { margin-bottom: 1rem; }
.field label { display: block; margin-bottom: 0.4rem; color: #8b949e; font-size: 0.9rem; }
.field input, .field select, .field textarea { width: 100%; padding: 0.6rem; background: #0d1117; border: 1px solid #30363d; border-radius: 6px; color: #c9d1d9; font-size: 0.9rem; font-family: monospace; }
.field input:focus, .field select:focus, .field textarea:focus { outline: none; border-color: #58a6ff; }
.hint { color: #8b949e; font-size: 0.75rem; margin-top: 0.3rem; display: block; }
.checkbox { display: flex; align-items: center; gap: 0.5rem; cursor: pointer; }
.checkbox input { width: auto; }
.form-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 1.5rem; }
.btn { padding: 0.5rem 1rem; border: 1px solid #30363d; border-radius: 6px; cursor: pointer; background: #21262d; color: #c9d1d9; }
.btn-primary { background: #238636; border-color: #2ea043; color: #fff; }
.btn-primary:hover { background: #2ea043; }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.error { color: #da3633; margin-bottom: 1rem; }
</style>
