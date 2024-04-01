<script setup>
import DirectoryConfig from './components/DirectoryConfig.vue';
import ProfileConfig from './components/ProfileConfig.vue';
import ModConfig from './components/ModConfig.vue';

</script>

<template>
  <div className="layout">
    <DirectoryConfig />
    <ProfileConfig />
    <ModConfig />
  </div>
</template>

<style>
.layout {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  height: 100%;
}
</style>
