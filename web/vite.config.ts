import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import { loadEnv } from 'vite';
import { defineConfig } from 'vitest/config';

export default defineConfig(({ mode }) => {
  // Daemon for local dev: `travo-greenfield --mock` (default :8081); override with DAEMON_URL.
  const target = loadEnv(mode, '.', '').DAEMON_URL || 'http://127.0.0.1:8081';
  const proxy = { '/api': target };
  return {
    plugins: [tailwindcss(), svelte()],
    build: {
      // Flat output: the device server serves dist/ as-is, with the *.gz written by scripts/gzip-dist.mjs.
      assetsDir: '',
      cssCodeSplit: false,
      modulePreload: { polyfill: false },
      rollupOptions: {
        output: {
          manualChunks(id) {
            // Keep the eager shell runtime together. Route chunks stay lazy, while
            // extra shared entry chunks cost over 1 KB in aggregate gzip overhead.
            if (
              id.endsWith('/src/main.ts') ||
              id.includes('/node_modules/svelte/') ||
              id.endsWith('/src/lib/session.svelte.ts')
            ) {
              return 'index';
            }
          },
        },
      },
    },
    server: { proxy },
    preview: { proxy },
    test: {
      include: ['src/**/*.test.ts', 'scripts/**/*.test.mjs'],
    },
  };
});
