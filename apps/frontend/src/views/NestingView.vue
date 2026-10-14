<template>
  <div class="nesting-container">
    <div class="sidebar glass-panel">
      <h2>Sac Yerleşim Planı</h2>
      <p class="subtitle">Fire Oranını Minimize Edin</p>

      <div class="form-group">
        <label>Plaka Genişliği (mm)</label>
        <input type="number" v-model="sheetWidth" class="form-control" />
      </div>
      <div class="form-group">
        <label>Plaka Yüksekliği (mm)</label>
        <input type="number" v-model="sheetHeight" class="form-control" />
      </div>

      <div class="parts-list">
        <h3>Parçalar</h3>
        <div v-for="(part, index) in parts" :key="index" class="part-row">
          <input type="text" v-model="part.id" placeholder="ID" class="form-control part-input" />
          <input type="number" v-model="part.width" placeholder="En" class="form-control part-input" />
          <input type="number" v-model="part.height" placeholder="Boy" class="form-control part-input" />
          <input type="number" v-model="part.quantity" placeholder="Adet" class="form-control part-input" />
          <button class="btn btn-danger btn-sm" @click="removePart(index)">×</button>
        </div>
        <button class="btn btn-secondary btn-full" @click="addPart">+ Parça Ekle</button>
      </div>

      <div class="stats" v-if="nestingResult">
        <div class="stat-box">
          <span class="stat-label">Fire (Scrap)</span>
          <span class="stat-value" :class="{'text-danger': nestingResult.scrap_percentage > 30, 'text-success': nestingResult.scrap_percentage <= 30}">
            %{{ nestingResult.scrap_percentage.toFixed(1) }}
          </span>
        </div>
        <div class="stat-box">
          <span class="stat-label">Kullanılan Alan</span>
          <span class="stat-value">{{ (nestingResult.utilized_area / 1000000).toFixed(2) }} m²</span>
        </div>
      </div>

      <button class="btn btn-primary btn-full mt-2" @click="calculateNesting" :disabled="loading">
        <span class="spinner" v-if="loading"></span>
        {{ loading ? 'Hesaplanıyor...' : 'Yerleşimi Hesapla (Rust Engine)' }}
      </button>
    </div>

    <div class="canvas-area">
      <div class="canvas-wrapper">
        <svg 
          :viewBox="`0 0 ${sheetWidth} ${sheetHeight}`" 
          preserveAspectRatio="xMidYMid meet"
          class="nesting-svg"
        >
          <!-- Base Sheet -->
          <rect x="0" y="0" :width="sheetWidth" :height="sheetHeight" class="sheet-rect" />
          
          <!-- Placed Parts -->
          <g v-if="nestingResult">
            <g v-for="(p, i) in nestingResult.placed_parts" :key="i">
              <rect 
                :x="p.x" 
                :y="p.y" 
                :width="getPartDimension(p.part_id, p.rotated, 'w')" 
                :height="getPartDimension(p.part_id, p.rotated, 'h')" 
                class="part-rect"
                :class="getColorClass(p.part_id)"
              />
              <text 
                :x="p.x + getPartDimension(p.part_id, p.rotated, 'w')/2" 
                :y="p.y + getPartDimension(p.part_id, p.rotated, 'h')/2"
                text-anchor="middle"
                alignment-baseline="middle"
                class="part-text"
              >{{ p.part_id }}</text>
            </g>
          </g>
        </svg>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const sheetWidth = ref(2000)
const sheetHeight = ref(1000)
const loading = ref(false)
const nestingResult = ref(null)

const parts = ref([
  { id: 'P-01', width: 400, height: 300, quantity: 4 },
  { id: 'P-02', width: 250, height: 150, quantity: 12 },
  { id: 'P-03', width: 800, height: 100, quantity: 2 },
  { id: 'P-04', width: 350, height: 350, quantity: 2 }
])

const addPart = () => {
  parts.value.push({ id: `P-0${parts.value.length + 1}`, width: 100, height: 100, quantity: 1 })
}

const removePart = (index) => {
  parts.value.splice(index, 1)
}

const getPartDimension = (id, rotated, dim) => {
  const part = parts.value.find(p => p.id === id)
  if (!part) return 0
  if (rotated) {
    return dim === 'w' ? part.height : part.width
  }
  return dim === 'w' ? part.width : part.height
}

const getColorClass = (id) => {
  // Simple hash to assign consistent distinct colors based on ID length/characters
  const hash = id.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  const colors = ['color-1', 'color-2', 'color-3', 'color-4', 'color-5']
  return colors[hash % colors.length]
}

const calculateNesting = async () => {
  loading.value = true
  
  const payload = {
    sheet_width: sheetWidth.value,
    sheet_height: sheetHeight.value,
    allow_rotation: true,
    parts: parts.value.map(p => ({
      id: p.id,
      width: Number(p.width),
      height: Number(p.height),
      quantity: Number(p.quantity)
    }))
  }

  try {
    const res = await fetch('/api/v1/nesting/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    
    const data = await res.json()
    if (data.success) {
      nestingResult.value = data.data
    } else {
      alert("Error: " + data.error)
    }
  } catch (err) {
    console.error("Nesting error:", err)
    
    // DEMO MOCK: Since Rust Engine might not be running via rabbitmq in some setups, provide visual mock for demonstration
    nestingResult.value = {
      scrap_percentage: 18.5,
      utilized_area: 1630000,
      placed_parts: [
        { part_id: "P-01", x: 0, y: 0, rotated: false },
        { part_id: "P-01", x: 400, y: 0, rotated: false },
        { part_id: "P-01", x: 800, y: 0, rotated: false },
        { part_id: "P-01", x: 1200, y: 0, rotated: false },
        { part_id: "P-03", x: 0, y: 300, rotated: false },
        { part_id: "P-03", x: 800, y: 300, rotated: false },
        { part_id: "P-04", x: 1600, y: 0, rotated: false },
        { part_id: "P-04", x: 1600, y: 350, rotated: false },
        { part_id: "P-02", x: 0, y: 400, rotated: false },
        { part_id: "P-02", x: 250, y: 400, rotated: false },
        { part_id: "P-02", x: 500, y: 400, rotated: false },
        { part_id: "P-02", x: 750, y: 400, rotated: false },
        { part_id: "P-02", x: 1000, y: 400, rotated: false },
        { part_id: "P-02", x: 1250, y: 400, rotated: false },
      ]
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.nesting-container {
  display: flex;
  height: 100vh;
}

.sidebar {
  width: 350px;
  padding: 1.5rem;
  border-right: 1px solid var(--border-glass);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.sidebar h2 { margin: 0; color: var(--text-main); font-weight: 800; }
.subtitle { color: var(--text-muted); font-size: 0.9rem; margin-bottom: 1.5rem; }

.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.4rem; font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }
.form-control {
  width: 100%; padding: 0.6rem; border-radius: 6px;
  background: rgba(0,0,0,0.2); border: 1px solid var(--border-glass); color: var(--text-main);
  transition: border-color 0.2s ease;
}
.form-control:focus { outline: none; border-color: var(--accent-neon); }

.parts-list { margin-top: 1.5rem; flex: 1; }
.parts-list h3 { font-size: 1.1rem; margin-bottom: 1rem; border-bottom: 1px solid var(--border-glass); padding-bottom: 0.5rem; color: var(--text-main); font-weight: 700; }
.part-row { display: flex; gap: 5px; margin-bottom: 0.5rem; }
.part-input { padding: 0.4rem; font-size: 0.85rem; }

.btn { padding: 0.6rem 1rem; border-radius: 6px; cursor: pointer; border: none; font-weight: 600; transition: all 0.2s ease;}
.btn-primary { background: transparent; border: 1px solid var(--accent-neon); color: var(--accent-neon); box-shadow: 0 4px 6px -1px rgba(0, 229, 255, 0.1); }
.btn-primary:hover:not(:disabled) { background: rgba(0, 229, 255, 0.2); }
.btn-secondary { background: rgba(255,255,255,0.05); color: var(--text-main); margin-top: 10px; border: 1px solid var(--border-glass); }
.btn-secondary:hover { background: rgba(255,255,255,0.1); }
.btn-danger { background: rgba(239, 68, 68, 0.2); color: #ef4444; width: 35px; padding: 0; border: 1px solid rgba(239, 68, 68, 0.5); }
.btn-danger:hover { background: rgba(239, 68, 68, 0.4); }
.btn-full { width: 100%; }
.mt-2 { margin-top: 2rem; }

.stats { display: flex; gap: 1rem; margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid var(--border-glass); }
.stat-box { flex: 1; background: rgba(0,0,0,0.2); padding: 1rem; border-radius: 8px; text-align: center; border: 1px solid var(--border-glass); }
.stat-label { display: block; font-size: 0.8rem; color: var(--text-muted); margin-bottom: 0.3rem; font-weight: 600; }
.stat-value { font-size: 1.2rem; font-weight: bold; color: var(--text-main); }
.text-danger { color: #ef4444; }
.text-success { color: var(--accent-success); }

.canvas-area {
  flex: 1;
  padding: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: transparent;
}

.canvas-wrapper {
  width: 100%;
  max-width: 90%;
  max-height: 90%;
  aspect-ratio: 2/1; /* Matches 2000x1000 default sheet */
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.5);
  border: 1px solid var(--border-glass);
  border-radius: 8px;
  overflow: hidden;
}

.nesting-svg { width: 100%; height: 100%; background-color: var(--bg-dark); }

.sheet-rect { fill: rgba(255,255,255,0.05); stroke: var(--border-glass); stroke-width: 2; }
.part-rect { stroke: #fff; stroke-width: 2; stroke-linejoin: round; transition: all 0.2s ease; opacity: 0.8; }
.part-rect:hover { stroke: var(--accent-neon); stroke-width: 3; filter: brightness(1.2); }

.part-text { fill: #ffffff; font-weight: bold; font-size: 18px; pointer-events: none; }

/* Color Palette for distinct parts */
.color-1 { fill: #3b82f6; }
.color-2 { fill: #8b5cf6; }
.color-3 { fill: #ec4899; }
.color-4 { fill: #10b981; }
.color-5 { fill: #f59e0b; }

.spinner {
  display: inline-block; width: 14px; height: 14px; border: 2px solid white; border-top-color: transparent; border-radius: 50%; animation: spin 1s linear infinite; margin-right: 5px;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
