<script setup>
import { useAppStore } from './stores/app';
import DirectoryConfig from './components/DirectoryConfig.vue';
import ProfileConfig from './components/ProfileConfig.vue';
import ModConfig from './components/ModConfig.vue';
import { onMounted } from 'vue'
const store = useAppStore();

onMounted(() => {
  // Setup Listeners
  window.runtime.EventsOn('app:modFolder', (modFolder) => {
    console.log("received event: modFolder", modFolder)
    store.setModFolder(modFolder)
  });
  window.runtime.EventsOn('app:appReady', () => {
    console.log("received event: appReady")
    store.setReady()
  });

  // Tell the backend we are ready
  window.runtime.EventsEmit('view:setupComplete')
})

</script>

<template>
  <div class="loading" v-if="!store.ready">
    <h1>🚀 Loading...</h1>
  </div>
  <div className="layout" v-else>
    <DirectoryConfig />
    <ProfileConfig v-if="store.modFolder != ''" />
    <ModConfig v-if="store.modFolder != ''" />
    <div v-else class="notConfigured">
      <h3>👆 Please provide the location for where your mods will be stored. 👆</h3>
    </div>
  </div>
</template>

<style>
.notConfigured {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}
.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.layout {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  height: 100%;
}
</style>
