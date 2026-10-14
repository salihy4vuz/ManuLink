<template>
  <div class="erp-status">
    <h3>ERP Sync Status</h3>
    <table>
      <thead>
        <tr><th>ERP</th><th>Status</th><th>Last Sync</th><th>Action</th></tr>
      </thead>
      <tbody>
        <tr v-for="erp in erps" :key="erp.name">
          <td>{{ erp.name }}</td>
          <td>{{ erp.status }}</td>
          <td>{{ erp.lastSync }}</td>
          <td><button @click="sync(erp)">Sync Now</button></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const erps = ref([
  { name: 'SAP', status: '✅ Başarılı', lastSync: '10 mins ago' },
  { name: 'Logo', status: '⏳ Bekliyor', lastSync: 'N/A' },
])
const sync = (erp) => { erp.status = '🔄 Syncing...' }
</script>
<style scoped>
table { width: 100%; border-collapse: collapse; }
th, td { padding: 8px; border: 1px solid #ddd; }
</style>
