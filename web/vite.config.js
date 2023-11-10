import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
    build: {
        rollupOptions: {
            output: {
                manualChunks(id, {getModuleInfo, getModuleIds}) {
                    if (id.includes('monaco-editor')) {
                        return 'monaco-editor';
                    }
                    if (id.includes('node_modules')) {
                        return 'vendor';
                    }
                },
            },
        },
    }
})
