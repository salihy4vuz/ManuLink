<template>
  <div class="page-container">
    <NavBar />
    
    <div class="content-wrapper">
      <div class="header">
        <h1 class="title-gradient">System Telemetry</h1>
        <p class="text-muted">Global operations, LME materials, and IIoT machine status monitor.</p>
      </div>

      <div class="dashboard-grid">
        
        <!-- Live LME Data Panel -->
        <div class="glass-panel widget lme-panel">
          <div class="widget-header">
            <h3><span class="icon">📈</span> LME Metal Rates</h3>
            <button class="btn-icon" title="Force Sync">🔄</button>
          </div>
          
          <div class="lme-list">
            <div class="lme-item">
              <div class="lme-info">
                <h4>Mild Steel (S235)</h4>
                <span class="lme-currency">EUR / kg</span>
              </div>
              <div class="lme-price">
                <input type="number" value="1.05" class="neon-input" />
                <span class="trend up">+2.1%</span>
              </div>
            </div>
            
            <div class="lme-item">
              <div class="lme-info">
                <h4>Stainless 304</h4>
                <span class="lme-currency">EUR / kg</span>
              </div>
              <div class="lme-price">
                <input type="number" value="3.50" class="neon-input" />
                <span class="trend down">-0.5%</span>
              </div>
            </div>

            <div class="lme-item">
              <div class="lme-info">
                <h4>Aluminum 5052</h4>
                <span class="lme-currency">EUR / kg</span>
              </div>
              <div class="lme-price">
                <input type="number" value="2.80" class="neon-input" />
                <span class="trend up">+1.2%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Global Metrics -->
        <div class="glass-panel widget metrics-panel">
          <div class="widget-header">
            <h3><span class="icon">🌍</span> API & Engine Load</h3>
          </div>
          <div class="metrics-grid">
            <div class="metric-box">
              <span class="metric-label">Rust Geometry Queue</span>
              <span class="metric-value warning">3 Jobs</span>
            </div>
            <div class="metric-box">
              <span class="metric-label">Total Quotes Processed</span>
              <span class="metric-value primary">1,402</span>
            </div>
            <div class="metric-box">
              <span class="metric-label">Gateway WebSocket</span>
              <span class="metric-value success">Healthy (34 ms)</span>
            </div>
            <div class="metric-box">
              <span class="metric-label">MinIO Storage</span>
              <span class="metric-value primary">42.5 GB</span>
            </div>
          </div>
        </div>

        <!-- IIoT MTConnect Machine Status -->
        <div class="glass-panel widget full-width machine-panel">
          <div class="widget-header">
            <h3><span class="icon">⚙️</span> IIoT Live Machine Feed</h3>
            <span class="live-indicator">LIVE</span>
          </div>
          
          <div class="machine-list">
            <!-- Machine 1 -->
            <div class="machine-card active">
              <div class="machine-icon-wrap laser">⚡</div>
              <div class="machine-details">
                <h4>Trumpf TruLaser 5030</h4>
                <div class="machine-stats">
                  <div class="stat"><span class="label">Status</span><span class="value success">Cutting</span></div>
                  <div class="stat"><span class="label">OEE</span><span class="value">85%</span></div>
                  <div class="stat"><span class="label">Gas</span><span class="value">N2 (14 Bar)</span></div>
                </div>
              </div>
              <div class="pulse-ring green"></div>
            </div>

            <!-- Machine 2 -->
            <div class="machine-card warning">
              <div class="machine-icon-wrap laser">⚡</div>
              <div class="machine-details">
                <h4>Amada ENSIS 3015</h4>
                <div class="machine-stats">
                  <div class="stat"><span class="label">Status</span><span class="value warn">Maint. Req</span></div>
                  <div class="stat"><span class="label">OEE</span><span class="value">0%</span></div>
                  <div class="stat"><span class="label">Lens</span><span class="value">88% Wear</span></div>
                </div>
              </div>
              <div class="pulse-ring orange"></div>
            </div>

            <!-- Machine 3 -->
            <div class="machine-card idle">
              <div class="machine-icon-wrap bend">📐</div>
              <div class="machine-details">
                <h4>Bystronic Xpert 150</h4>
                <div class="machine-stats">
                  <div class="stat"><span class="label">Status</span><span class="value info">Idle (Setup)</span></div>
                  <div class="stat"><span class="label">OEE</span><span class="value">--</span></div>
                  <div class="stat"><span class="label">Tool</span><span class="value">V12 Die</span></div>
                </div>
              </div>
              <div class="pulse-ring blue"></div>
            </div>
          </div>
        </div>

        <!-- White-Label Theme Settings -->
        <ThemeCustomizer class="full-width" />

      </div>
    </div>
  </div>
</template>

<script setup>
import NavBar from '../components/NavBar.vue'
import ThemeCustomizer from '../components/ThemeCustomizer.vue'
</script>

<style scoped>
.page-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.content-wrapper {
  padding: 40px;
  flex: 1;
  overflow-y: auto;
  background: radial-gradient(circle at top right, rgba(15, 23, 42, 1) 0%, rgba(2, 6, 23, 1) 100%);
}

.header {
  margin-bottom: 30px;
}

.title-gradient {
  background: linear-gradient(to right, #38bdf8, #818cf8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 2.5rem;
  margin: 0 0 0.5rem 0;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  max-width: 1400px;
}

.glass-panel {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

.widget {
  padding: 2rem;
  display: flex;
  flex-direction: column;
}

.full-width {
  grid-column: 1 / -1;
}

.widget-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 1rem;
}

.widget-header h3 {
  margin: 0;
  color: #0f172a;
  font-size: 1.2rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-icon {
  background: transparent;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  opacity: 0.7;
  transition: all 0.2s ease;
  color: #475569;
}

.btn-icon:hover {
  opacity: 1;
  transform: rotate(180deg);
}

/* LME Panel */
.lme-list { display: flex; flex-direction: column; gap: 1.5rem; }
.lme-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8fafc;
  padding: 1rem 1.5rem;
  border-radius: 8px;
}
.lme-info h4 { margin: 0 0 0.25rem 0; color: #0f172a; font-size: 1.05rem; }
.lme-currency { font-size: 0.8rem; color: #64748b; font-weight: 500; }
.lme-price { display: flex; align-items: center; gap: 1rem; }

.neon-input {
  background: #fff;
  border: 1px solid #cbd5e1;
  color: #0f172a;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  width: 90px;
  font-size: 1.1rem;
  font-weight: bold;
  font-family: 'Fira Code', monospace;
  text-align: right;
  transition: all 0.2s ease;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}
.neon-input:focus { outline: none; border-color: #1d6bf3; box-shadow: 0 0 0 3px rgba(29, 107, 243, 0.2); }

.trend { font-size: 0.9rem; font-weight: bold; padding: 4px 8px; border-radius: 4px; }
.trend.up { background: #ecfdf5; color: #10b981; }
.trend.down { background: #fef2f2; color: #ef4444; }

/* Metrics Panel */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}
.metric-box {
  background: #f8fafc;
  padding: 1.5rem;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.metric-label { font-size: 0.85rem; color: #64748b; font-weight: 600; text-transform: uppercase; }
.metric-value { font-size: 1.8rem; font-weight: 800; font-family: -apple-system, BlinkMacSystemFont, 'Inter', sans-serif; }
.metric-value.primary { color: #1d6bf3; }
.metric-value.warning { color: #f59e0b; }
.metric-value.success { color: #10b981; }

/* Machine Panel */
.live-indicator {
  background: #fef2f2;
  color: #ef4444;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: bold;
  border: 1px solid #fca5a5;
  animation: pulse-red 2s infinite;
}

@keyframes pulse-red {
  0% { box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.4); }
  70% { box-shadow: 0 0 0 6px rgba(239, 68, 68, 0); }
  100% { box-shadow: 0 0 0 0 rgba(239, 68, 68, 0); }
}

.machine-list { display: flex; gap: 1.5rem; }
.machine-card {
  flex: 1;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: flex-start;
  gap: 1.25rem;
  position: relative;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.machine-card.active { border-color: #a7f3d0; background: #ecfdf5; }
.machine-card.warning { border-color: #fde68a; background: #fffbeb; }

.machine-icon-wrap {
  width: 50px; height: 50px;
  border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  font-size: 1.5rem;
  background: #f1f5f9;
}

.machine-details { flex: 1; z-index: 2; }
.machine-details h4 { margin: 0 0 1rem 0; color: #0f172a; font-size: 1.1rem; font-weight: 700; }

.machine-stats {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.stat { display: flex; justify-content: space-between; font-size: 0.85rem; }
.stat .label { color: #64748b; }
.stat .value { color: #0f172a; font-weight: 600; }
.stat .value.success { color: #10b981; }
.stat .value.warn { color: #f59e0b; }
.stat .value.info { color: #1d6bf3; }

/* Status Pulses */
.pulse-ring {
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  width: 12px; height: 12px;
  border-radius: 50%;
}
.pulse-ring.green { background: #10b981; box-shadow: 0 0 5px rgba(16,185,129,0.5); animation: glow-green 2s infinite alternate; }
.pulse-ring.orange { background: #f59e0b; box-shadow: 0 0 5px rgba(245,158,11,0.5); animation: glow-orange 1s infinite alternate; }
.pulse-ring.blue { background: #1d6bf3; box-shadow: 0 0 5px rgba(29,107,243,0.5); }

@keyframes glow-green { from { opacity: 0.6; } to { opacity: 1; box-shadow: 0 0 10px #10b981; } }
@keyframes glow-orange { from { opacity: 0.5; } to { opacity: 1; box-shadow: 0 0 10px #f59e0b; } }

</style>
