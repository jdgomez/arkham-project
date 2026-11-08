import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';

const backend = 'http://localhost:8080';

export default defineConfig({
  plugins: [svelte()],
  build: { outDir: '../backend/internal/web/dist', emptyOutDir: false },
  test: { environment: 'jsdom', include: ['src/**/*.test.ts'] },
  resolve: process.env.VITEST ? { conditions: ['browser'] } : undefined,
  server: {
    proxy: { '/api': backend },
  },
});
