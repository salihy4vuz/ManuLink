<script setup>
import { ref, watch, onMounted } from 'vue'
import { TresCanvas } from '@tresjs/core'
import { OrbitControls, Stars } from '@tresjs/cientos'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader'

const props = defineProps({
  gltfData: {
    type: String,
    default: null
  }
})

const modelRef = ref(null)

watch(() => props.gltfData, (newData) => {
  if (newData) {
    const loader = new GLTFLoader()
    loader.parse(newData, '', (gltf) => {
      modelRef.value = gltf.scene
    }, (error) => {
      console.error("Error parsing GLTF:", error)
    })
  }
}, { immediate: true })
</script>

<template>
  <TresCanvas clear-color="#0f172a" shadows>
    <TresPerspectiveCamera
      :position="[5, 5, 5]"
      :look-at="[0, 0, 0]"
    />
    <OrbitControls />
    <Stars />
    <TresAmbientLight :intensity="0.5" />
    <TresDirectionalLight
      :position="[5, 10, 5]"
      :intensity="1"
      cast-shadow
      :shadow-mapSize-width="1024"
      :shadow-mapSize-height="1024"
    />
    
    <primitive v-if="modelRef" :object="modelRef" cast-shadow receive-shadow />
    <TresMesh v-else cast-shadow>
      <TresBoxGeometry :args="[2, 0.1, 2]" />
      <TresMeshStandardMaterial color="#38bdf8" metalness="0.8" roughness="0.2" wireframe />
    </TresMesh>

    <!-- Floor for shadow reception -->
    <TresMesh :position="[0, -1, 0]" :rotation="[-Math.PI / 2, 0, 0]" receive-shadow>
      <TresPlaneGeometry :args="[20, 20]" />
      <TresMeshStandardMaterial color="#0f172a" metalness="0.5" roughness="0.8" />
    </TresMesh>
  </TresCanvas>
</template>

