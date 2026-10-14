<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { 
  Play, Pause, Square, AlertTriangle, Thermometer, 
  Zap, Gauge, Wind, Activity, CheckCircle2, Factory, ChevronRight
} from 'lucide-vue-next'
import NavBar from '../components/NavBar.vue'
import ThreeDViewer from '../components/ThreeDViewer.vue'

// CNC Machine Telemetry Simulation
const telemetry = ref({
  status: 'RUNNING', // RUNNING, IDLE, ERROR
  laserPower: 3200, // Watts
  gasPressure: 12.4, // Bar
  temperature: 42, // Celsius
  feedRate: 1500, // mm/min
  progress: 68 // %
})

// Current Job
const activeJob = ref({
  id: 'JOB-9024',
  part: 'Chassis_Main.dxf',
  material: 'Stainless Steel 304',
  thickness: '3.0mm',
  qty: 150,
  completed: 102,
  startTime: '08:45 AM',
  estCompletion: '11:20 AM'
})

// Queue
const queue = ref([
  { id: 'JOB-9025', part: 'Bracket_Angle.step', material: 'Aluminum 6061', thickness: '2.0mm', qty: 50, time: '45m' },
  { id: 'JOB-9026', part: 'Enclosure_Top.dxf', material: 'Carbon Steel', thickness: '1.5mm', qty: 200, time: '1h 20m' },
  { id: 'JOB-9027', part: 'Custom_Flange.dxf', material: 'Stainless 316', thickness: '5.0mm', qty: 10, time: '15m' }
])

let interval;
onMounted(() => {
  // Simulate live progress
  interval = setInterval(() => {
    if (telemetry.value.status === 'RUNNING') {
      telemetry.value.progress += 0.5
      if (telemetry.value.progress > 100) {
        telemetry.value.progress = 0
        activeJob.value.completed++
      }
      // Fluctuate telemetry
      telemetry.value.laserPower = 3100 + Math.random() * 200
      telemetry.value.gasPressure = 12 + Math.random()
    }
  }, 1000)
})

onUnmounted(() => clearInterval(interval))

const actionCycleStart = () => telemetry.value.status = 'RUNNING'
const actionFeedHold = () => telemetry.value.status = 'IDLE'
const actionEStop = () => telemetry.value.status = 'ERROR'
</script>

<template>
  <div class="cnc-dashboard">
    <NavBar />

    <main class="dashboard-grid">
      <!-- TOP STATUS BAR (HMI Style) -->
      <header class="status-bar glass-panel">
        <div class="machine-id">
          <Factory class="icon text-blue-400" size="24" />
          <div>
            <h2>APEX LASER CELL 01</h2>
            <span class="model">TRUMPF / 6KW FIBER</span>
          </div>
        </div>

        <div class="status-indicator" :class="telemetry.status.toLowerCase()">
          <div class="pulse-dot"></div>
          <span class="status-text">{{ telemetry.status }}</span>
        </div>

        <div class="telemetry-gauges">
          <div class="gauge">
            <Zap class="text-yellow-400" size="18" />
            <div class="val">{{ telemetry.laserPower.toFixed(0) }} <span>W</span></div>
          </div>
          <div class="gauge">
            <Wind class="text-cyan-400" size="18" />
            <div class="val">{{ telemetry.gasPressure.toFixed(1) }} <span>Bar</span></div>
          </div>
          <div class="gauge">
            <Thermometer :class="telemetry.temperature > 50 ? 'text-red-500' : 'text-orange-400'" size="18" />
            <div class="val">{{ telemetry.temperature }} <span>°C</span></div>
          </div>
          <div class="gauge">
            <Gauge class="text-green-400" size="18" />
            <div class="val">{{ telemetry.feedRate }} <span>mm/m</span></div>
          </div>
        </div>
      </header>

      <div class="main-content">
        <!-- LEFT: ACTIVE JOB & 3D -->
        <section class="active-job-section glass-panel">
          <div class="section-header">
            <h3>ACTIVE PROGRAM</h3>
            <span class="job-id">{{ activeJob.id }}</span>
          </div>

          <div class="job-meta">
            <h1 class="part-name">{{ activeJob.part }}</h1>
            <div class="material-tag">
              {{ activeJob.material }} | {{ activeJob.thickness }}
            </div>
          </div>

          <!-- 3D Render Area -->
          <div class="render-viewport">
            <div class="overlay-hud">
              <div class="hud-item">
                <span>Z-Height</span>
                <strong>1.50mm</strong>
              </div>
              <div class="hud-item">
                <span>Assist Gas</span>
                <strong>O2 / N2</strong>
              </div>
            </div>
            <!-- We wrap ThreeDViewer in a fixed container for the HMI look -->
            <ThreeDViewer />
          </div>

          <div class="progress-section">
            <div class="progress-stats">
              <div class="stat">
                <span class="label">Completed</span>
                <span class="value">{{ activeJob.completed }} / {{ activeJob.qty }}</span>
              </div>
              <div class="stat text-right">
                <span class="label">Time Remaining</span>
                <span class="value">02:35:10</span>
              </div>
            </div>
            <div class="progress-bar-bg">
              <div class="progress-bar-fill" :style="{ width: `${telemetry.progress}%` }">
                <div class="progress-glow"></div>
              </div>
            </div>
          </div>
        </section>

        <!-- RIGHT: QUEUE & ACTIONS -->
        <aside class="right-panel">
          
          <section class="queue-section glass-panel">
            <div class="section-header">
              <h3>PRODUCTION QUEUE</h3>
              <span class="badge">{{ queue.length }} Jobs</span>
            </div>
            <div class="queue-list">
              <div v-for="(q, idx) in queue" :key="q.id" class="queue-item">
                <div class="q-index">{{ idx + 1 }}</div>
                <div class="q-info">
                  <h4>{{ q.part }}</h4>
                  <p>{{ q.material }} - {{ q.thickness }} ({{ q.qty }} pcs)</p>
                </div>
                <div class="q-time">{{ q.time }}</div>
                <ChevronRight class="text-gray-500" size="18" />
              </div>
            </div>
          </section>

          <!-- HMI HARD BUTTONS -->
          <section class="action-buttons glass-panel">
            <button @click="actionCycleStart" class="hmi-btn btn-start" :class="{ active: telemetry.status === 'RUNNING' }">
              <Play size="32" />
              <span>CYCLE START</span>
            </button>
            <button @click="actionFeedHold" class="hmi-btn btn-hold" :class="{ active: telemetry.status === 'IDLE' }">
              <Pause size="32" />
              <span>FEED HOLD</span>
            </button>
            <button @click="actionEStop" class="hmi-btn btn-estop">
              <AlertTriangle size="32" />
              <span>EMERGENCY STOP</span>
            </button>
          </section>

        </aside>
      </div>
    </main>
  </div>
</template>

<style scoped>
/* Core Layout */
.cnc-dashboard {
  min-height: 100vh;
  background: #f8fafc;
  color: #0f172a;
  font-family: -apple-system, BlinkMacSystemFont, 'Inter', 'Segoe UI', sans-serif;
}

.dashboard-grid {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  height: calc(100vh - 64px);
}

/* Glass Panels (Industrial HMI variant) */
.glass-panel {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
  border-radius: 12px 12px 0 0;
}

.section-header h3 {
  font-size: 0.85rem;
  font-weight: 700;
  color: #64748b;
  letter-spacing: 0.1em;
  margin: 0;
}

/* TOP STATUS BAR */
.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
}

.machine-id {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.machine-id h2 {
  font-size: 1.25rem;
  font-weight: 800;
  margin: 0;
  color: #0f172a;
  letter-spacing: 1px;
}

.machine-id .model {
  font-size: 0.75rem;
  color: #38bdf8;
  font-weight: 600;
  letter-spacing: 2px;
}

/* Status Indicator (RUNNING, IDLE, ERROR) */
.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 1.5rem;
  border-radius: 50px;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
}

.status-text {
  font-weight: 800;
  font-size: 1.25rem;
  letter-spacing: 2px;
}

.pulse-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.status-indicator.running .pulse-dot {
  background: #22c55e;
  box-shadow: 0 0 15px #22c55e;
  animation: pulse-green 1.5s infinite;
}
.status-indicator.running .status-text { color: #22c55e; }

.status-indicator.idle .pulse-dot { background: #eab308; }
.status-indicator.idle .status-text { color: #eab308; }

.status-indicator.error .pulse-dot { 
  background: #ef4444; 
  box-shadow: 0 0 20px #ef4444;
  animation: pulse-red 0.5s infinite alternate;
}
.status-indicator.error .status-text { color: #ef4444; }

@keyframes pulse-green {
  0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.7); }
  70% { transform: scale(1); box-shadow: 0 0 0 10px rgba(34, 197, 94, 0); }
  100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(34, 197, 94, 0); }
}

@keyframes pulse-red {
  from { opacity: 1; transform: scale(1); }
  to { opacity: 0.5; transform: scale(1.2); }
}

/* Telemetry Gauges */
.telemetry-gauges {
  display: flex;
  gap: 2rem;
}

.gauge {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
}

.gauge .val {
  font-family: 'Fira Code', monospace;
  font-size: 1.25rem;
  font-weight: 700;
  color: #0f172a;
}

.gauge .val span {
  font-size: 0.75rem;
  color: #64748b;
}

/* Main Content Layout */
.main-content {
  display: flex;
  gap: 1.5rem;
  flex: 1;
  overflow: hidden;
}

/* LEFT PANEL (ACTIVE JOB) */
.active-job-section {
  flex: 2;
  display: flex;
  flex-direction: column;
}

.job-meta {
  padding: 1.5rem;
  text-align: center;
}

.part-name {
  font-size: 2.5rem;
  font-weight: 800;
  margin: 0 0 0.5rem 0;
  color: #0f172a;
}

.material-tag {
  display: inline-block;
  background: rgba(56, 189, 248, 0.1);
  color: #38bdf8;
  border: 1px solid rgba(56, 189, 248, 0.3);
  padding: 0.25rem 1rem;
  border-radius: 50px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.875rem;
}

/* 3D Render Viewport */
.render-viewport {
  flex: 1;
  position: relative;
  margin: 0 1.5rem;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(255,255,255,0.05);
  background: #000;
  box-shadow: inset 0 0 50px rgba(0,0,0,0.8);
}

.overlay-hud {
  position: absolute;
  top: 1rem;
  left: 1rem;
  z-index: 10;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.hud-item {
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: 1px solid rgba(255,255,255,0.1);
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.75rem;
  display: flex;
  flex-direction: column;
}

.hud-item span { color: #94a3b8; }
.hud-item strong { color: #38bdf8; font-size: 0.9rem;}

/* Progress Bar */
.progress-section {
  padding: 1.5rem;
  background: #f8fafc;
  border-top: 1px solid #e2e8f0;
  margin-top: auto;
}

.progress-stats {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.stat { display: flex; flex-direction: column; }
.stat .label { font-size: 0.75rem; color: #64748b; text-transform: uppercase; letter-spacing: 1px;}
.stat .value { font-family: 'Fira Code', monospace; font-size: 1.25rem; font-weight: 700; color: #0f172a;}

.progress-bar-bg {
  height: 12px;
  background: #e2e8f0;
  border-radius: 50px;
  overflow: hidden;
  position: relative;
}

.progress-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #0284c7, #38bdf8);
  transition: width 0.3s ease;
  position: relative;
}

.progress-glow {
  position: absolute;
  top: 0; right: 0; bottom: 0; width: 20px;
  background: rgba(255,255,255,0.8);
  filter: blur(5px);
}


/* RIGHT PANEL */
.right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* Queue List */
.queue-section {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.queue-list {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  overflow-y: auto;
}

.queue-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  transition: background 0.2s;
}

.queue-item:hover {
  background: #f1f5f9;
}

.q-index {
  width: 24px; height: 24px;
  border-radius: 50%;
  background: #e2e8f0;
  color: #475569;
  display: flex; align-items: center; justify-content: center;
  font-size: 0.75rem; font-weight: bold;
}

.q-info { flex: 1; }
.q-info h4 { margin: 0; font-size: 0.95rem; font-weight: 700; color: #0f172a; }
.q-info p { margin: 0; font-size: 0.75rem; color: #64748b; }

.q-time {
  font-family: 'Fira Code', monospace;
  font-size: 0.85rem;
  color: #1d6bf3;
}

/* HMI Action Buttons (Giant Touch Buttons) */
.action-buttons {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background: #f8fafc;
  border-top: 1px solid #e2e8f0;
}

.hmi-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 1.5rem;
  border-radius: 8px;
  border: none;
  font-size: 1.25rem;
  font-weight: 800;
  letter-spacing: 2px;
  color: #fff;
  cursor: pointer;
  transition: all 0.1s;
  box-shadow: inset 0 2px 0 rgba(255,255,255,0.2), 0 4px 10px rgba(0,0,0,0.5);
}

.hmi-btn:active {
  transform: translateY(4px);
  box-shadow: inset 0 2px 0 rgba(255,255,255,0.0), 0 0px 0px rgba(0,0,0,0.5);
}

.btn-start {
  background: linear-gradient(180deg, #22c55e, #166534);
  border: 1px solid #14532d;
}
.btn-start:hover { background: linear-gradient(180deg, #4ade80, #15803d); }
.btn-start.active {
  box-shadow: inset 0 0 20px rgba(0,0,0,0.5), 0 0 20px rgba(34, 197, 94, 0.4);
  border-color: #4ade80;
}

.btn-hold {
  background: linear-gradient(180deg, #eab308, #854d0e);
  border: 1px solid #713f12;
}
.btn-hold:hover { background: linear-gradient(180deg, #fde047, #a16207); }
.btn-hold.active {
  box-shadow: inset 0 0 20px rgba(0,0,0,0.5), 0 0 20px rgba(234, 179, 8, 0.4);
  border-color: #fde047;
}

.btn-estop {
  background: linear-gradient(180deg, #ef4444, #7f1d1d);
  border: 1px solid #450a0a;
  margin-top: 1rem;
  padding: 2rem;
  font-size: 1.5rem;
}
.btn-estop:hover { background: linear-gradient(180deg, #f87171, #991b1b); }
</style>
