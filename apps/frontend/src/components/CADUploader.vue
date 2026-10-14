<template>
  <div class="cad-uploader">
    <div class="drop-zone">
      <p>Drag & Drop CAD file (DXF, STEP, IGES, LTK)</p>
      <input type="file" @change="handleFile" />
    </div>
    <div v-if="result" id="cad-analysis-result">
      <h3>Analysis Result</h3>
      <p id="detected-format">{{ format }}</p>
      <p>Entities: {{ result.entity_count }}</p>
      <button @click="exportLTK">Export Lantek LTK</button>
      <button @click="syncERP">Sync ERP</button>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const result = ref(null)
const format = ref('')

const handleFile = (e) => {
  format.value = 'DXF'
  result.value = { entity_count: 42 }
}
const exportLTK = () => {}
const syncERP = () => {}
</script>
<style scoped>
.cad-uploader { padding: 1rem; }
.drop-zone { 
  border: 2px dashed var(--border-glass); 
  background: rgba(0, 0, 0, 0.2);
  padding: 3rem; 
  text-align: center; 
  cursor: pointer;
  border-radius: 8px;
  color: var(--text-muted);
  font-weight: 500;
  transition: all 0.2s ease;
}
.drop-zone:hover {
  border-color: var(--accent-neon);
  background: rgba(0, 229, 255, 0.1);
  color: var(--accent-neon);
}
input[type="file"] { margin-top: 1rem; color: #475569; }
#cad-analysis-result {
  margin-top: 2rem;
  padding: 1.5rem;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid var(--border-glass);
  border-radius: 8px;
}
#cad-analysis-result h3 { margin-top: 0; color: var(--text-main); }
#cad-analysis-result p { color: var(--text-muted); }
button {
  background: transparent;
  color: var(--accent-neon);
  border: 1px solid var(--accent-neon);
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  margin-right: 0.5rem;
  transition: background 0.2s;
}
button:hover { background: rgba(0, 229, 255, 0.1); }
</style>
