import { defineConfig } from 'vitest/config';
import viteConfig from './vite.config.ts';

export default defineConfig((env) => {
  const config = typeof viteConfig === 'function' ? viteConfig(env) : viteConfig;
  return {
    ...config,
    test: {
      ...config.test,
      include: ['src/**/*.test.ts', 'scripts/**/*.test.mjs'],
    },
  };
});
