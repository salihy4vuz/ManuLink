<script setup>
import { ref } from 'vue'
import { UploadCloud, Loader2 } from 'lucide-vue-next'

const props = defineProps({
  isCalculating: Boolean
})

const emit = defineEmits(['file-dropped'])

const isHovering = ref(false)
const fileInput = ref(null)

const onDragOver = (e) => {
  e.preventDefault()
  if (!props.isCalculating) isHovering.value = true
}

const onDragLeave = () => {
  isHovering.value = false
}

const onDrop = (e) => {
  e.preventDefault()
  isHovering.value = false
  if (props.isCalculating) return

  const files = e.dataTransfer.files
  if (files.length > 0) {
    handleFile(files[0])
  }
}

const onClick = () => {
  if (!props.isCalculating) {
    fileInput.value.click()
  }
}

const onFileSelect = (e) => {
  const files = e.target.files
  if (files.length > 0) {
    handleFile(files[0])
  }
}

const handleFile = (file) => {
  const validExts = ['.step', '.stp', '.dxf']
  const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
  
  if (!validExts.includes(ext)) {
    alert("Please upload a valid 3D CAD file (.step, .stp, .dxf)")
    return
  }
  
  emit('file-dropped', file)
}
</script>

<template>
  <div 
    class="upload-zone"
    :class="{ 'hover-active': isHovering, 'calculating': isCalculating }"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
    @click="onClick"
  >
    <input 
      type="file" 
      ref="fileInput" 
      class="hidden-input" 
      accept=".step,.stp,.dxf"
      @change="onFileSelect"
    />
    
    <div v-if="!isCalculating">
      <div class="icon-wrap">
        <UploadCloud class="lucide-icon" size="56" />
      </div>
      <h3>Drop Engineering Files</h3>
      <p class="text-muted mono">Supported: .STEP, .STP, .DXF</p>
    </div>
    
    <div v-else class="loading-state">
      <Loader2 class="spinner lucide-icon" size="48" />
      <h3 class="text-accent">Parsing Geometry Data...</h3>
      <p class="text-muted mono">Extracting topological faces & bounding boxes.</p>
    </div>
  </div>
</template>

<style scoped>
.upload-zone {
  border: 2px dashed #cbd5e1;
  border-radius: 8px;
  padding: 60px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
  background: #f8fafc;
  position: relative;
}

.upload-zone:hover, .hover-active {
  border-color: #1d6bf3;
  background-color: #eff6ff;
}

.calculating {
  cursor: wait;
  border-color: #1d6bf3;
  background-color: #eff6ff;
}

.hidden-input {
  display: none;
}

.icon-wrap {
  margin-bottom: 1rem;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  padding: 15px;
  border-radius: 50%;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
}

.lucide-icon {
  color: #1d6bf3;
}

h3 {
  color: #0f172a;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.spinner {
  animation: spin 1.5s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.text-accent {
  color: #1d6bf3;
  margin: 0 0 0.5rem 0;
  font-weight: 600;
}

.text-muted {
  color: #64748b;
  margin: 0;
  font-size: 0.95rem;
}
</style>
