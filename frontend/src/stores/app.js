import { defineStore } from 'pinia'

export const useAppStore = defineStore({
  id: 'app',
  state: () => ({
    ready: false,
    modFolder: '',
  }),
  getters: {
    
  },
  actions: {
    setReady() {
        this.ready = true
    },
    setModFolder(location) {
      this.modFolder = location
    }
  }
})