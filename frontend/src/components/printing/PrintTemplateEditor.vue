<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-bold text-gray-800">📝 Templates</h2>
        <div class="flex gap-2">
          <button
            @click="activeTemplateTab = 'html'"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-bold transition-colors',
              activeTemplateTab === 'html'
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            ]"
          >
            🌐 Bill Template
          </button>
          <button
            @click="activeTemplateTab = 'label'"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-bold transition-colors',
              activeTemplateTab === 'label'
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            ]"
          >
            🏷️ Label Template
          </button>
          <button
            @click="activeTemplateTab = 'text'"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-bold transition-colors',
              activeTemplateTab === 'text'
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            ]"
          >
            📝 Text Templates
          </button>
        </div>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="flex-1 overflow-hidden">
      <!-- HTML Template Tab -->
      <div v-show="activeTemplateTab === 'html'" class="h-full">
        <HTMLTemplateEditor />
      </div>

      <!-- Label Template Tab -->
      <div v-show="activeTemplateTab === 'label'" class="h-full">
        <LabelTemplateEditor />
      </div>

      <!-- Text Templates Tab -->
      <div v-show="activeTemplateTab === 'text'" class="h-full">
        <TextTemplateList />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import HTMLTemplateEditor from './HTMLTemplateEditor.vue'
import LabelTemplateEditor from './LabelTemplateEditor.vue'
import TextTemplateList from './TextTemplateList.vue'

const activeTemplateTab = ref('html')
</script>
