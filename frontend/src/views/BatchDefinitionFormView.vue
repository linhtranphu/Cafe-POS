<template>
  <BatchDefinitionForm 
    :definition="definition"
    :mode="mode"
    @close="handleClose"
    @saved="handleSaved"
  />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useBatchDefinitionStore } from '../stores/batchDefinition'
import BatchDefinitionForm from '../components/batch/BatchDefinitionForm.vue'

const route = useRoute()
const router = useRouter()
const batchStore = useBatchDefinitionStore()

const definition = ref(null)
const mode = ref('create')

onMounted(async () => {
  if (route.params.id) {
    mode.value = 'edit'
    // Fetch the definition
    await batchStore.fetchDefinitions()
    definition.value = batchStore.definitions.find(d => d.id === route.params.id)
  }
})

const handleClose = () => {
  router.push('/batch/definitions')
}

const handleSaved = () => {
  router.push('/batch/definitions')
}
</script>
