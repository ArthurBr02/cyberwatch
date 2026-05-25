<script setup lang="ts">
import { ref } from 'vue'
import type { Source } from '../types'
import SourceList from '../components/SourceList.vue'
import SourceForm from '../components/SourceForm.vue'
import ArticleList from '../components/ArticleList.vue'

const showForm = ref(false)
const editingSource = ref<Source | null>(null)
const selectedSourceId = ref<number | undefined>(undefined)
const refreshKey = ref(0)

const onEdit = (source: Source) => {
  editingSource.value = source
  showForm.value = true
}

const onCreate = () => {
  editingSource.value = null
  showForm.value = true
}

const onSaved = () => {
  showForm.value = false
  editingSource.value = null
  refreshKey.value++
}

const onCancel = () => {
  showForm.value = false
  editingSource.value = null
}

const onSelectSource = (id: number | undefined) => {
  selectedSourceId.value = id
}
</script>

<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h2>Dashboard</h2>
    </div>

    <SourceList
      :key="'sources-' + refreshKey"
      @edit="onEdit"
      @create="onCreate"
    />

    <ArticleList
      :key="'articles-' + refreshKey"
      :source-id="selectedSourceId"
    />

    <SourceForm
      v-if="showForm"
      :source="editingSource"
      @saved="onSaved"
      @cancel="onCancel"
    />
  </div>
</template>

<style scoped>
.dashboard-header {
  margin-bottom: 1.5rem;
}
.dashboard-header h2 {
  font-size: 1.5rem;
  color: #c9d1d9;
}
</style>