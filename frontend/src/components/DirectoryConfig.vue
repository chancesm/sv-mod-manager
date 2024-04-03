<script setup>
import { useAppStore } from '../stores/app';
const store = useAppStore();
const openFolder = () => {
    console.log("openFolder")
    window.runtime.EventsEmit('openFolder')
}
const changeFolder = () => {
    console.log("changeFolder")
    window.runtime.EventsEmit('changeFolder')
}
</script>

<template>
    <div class="directoryConfig" >
        <div class="wrapper" v-if="store.modFolder == ''">
            <div class="folder">
                Please select the Stardew Valley Mods Directory
            </div>
            <div class="buttons">
                <button @click="changeFolder">Set Mods Folder</button>
            </div>
        </div>
        <div class="wrapper" v-else>
            <div class="folder">
                Stardew Valley Mods Directory: "{{ store.modFolder }}"
            </div>
            <div class="buttons">
                <button @click="openFolder">Open Folder</button>
                <button @click="changeFolder">Change</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.directoryConfig {
    margin-top: 1em;
    height: 100%;
    flex: 0.25;
}
.wrapper {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    gap: 1em;   
}
.buttons {
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    gap: 1em;   
}
</style>
