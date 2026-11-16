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
    resolve: { alias: { $lib: new URL('./src/lib', import.meta.url).pathname } },
    build: {
      // Flat output: the device server serves dist/ as-is, with the *.gz written by scripts/gzip-dist.mjs.
      assetsDir: '',
      cssCodeSplit: false,
      modulePreload: { polyfill: false },
      rollupOptions: {
        output: {
          codeSplitting: {
            includeDependenciesRecursively: false,
            groups: [
              {
                name: 'index',
                includeDependenciesRecursively: true,
                priority: -1,
                test: /\/node_modules\/svelte\/|\/src\/main.ts$|\/src\/lib\/session.svelte.ts$/,
              },
              { name: 'dialog-primitives', test: /\/bits-ui\/.*\/bits\// },
              { name: 'dialog-internal', test: /\/bits-ui\/.*\/internal\// },
              { name: 'dialog-utils', test: /\/(svelte-toolbelt|runed)\// },
            ],
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
