import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'
import {resolve} from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      wailsjs: resolve('./wailsjs'),
      components: resolve('./src/components'),
      scripts: resolve('./src/scripts'),
      stores: resolve('./src/stores'),
      style: resolve('./src/style'),
      src: resolve('./src'),
    }
  }
})
