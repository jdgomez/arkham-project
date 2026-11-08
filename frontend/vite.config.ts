import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

const backend = 'http://localhost:8080';

export default defineConfig({
  plugins: [svelte()],
  build: { outDir: '../backend/internal/web/dist', emptyOutDir: true },
  server: {
    proxy: { '/api': backend },
  },
});
