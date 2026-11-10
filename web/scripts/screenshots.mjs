// Writes PNG screenshots of the UI at mobile and desktop widths in light and dark scheme.
// Expects the daemon and `vite preview` already running (see `mise run screenshots`).
import { mkdirSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { launchChrome } from './cdp.mjs';

const BASE_URL = process.env.SCREENSHOT_URL ?? `http://localhost:${process.env.WEB_PORT ?? 18092}`;
const ROUTES = ['/', '/system'];
const VIEWPORTS = [
  { width: 390, height: 844, mobile: true },
  { width: 1280, height: 800, mobile: false },
];
const SCHEMES = ['light', 'dark'];

const OUT_DIR = fileURLToPath(new URL('../screenshots/', import.meta.url));

/** `/` becomes `dashboard`; other routes drop their slashes (`/system` -> `system`). */
export function routeName(route) {
  return route === '/' ? 'dashboard' : route.replace(/\//g, '');
}

export function screenshotPath(route, width, scheme) {
  return `${OUT_DIR}${routeName(route)}-${width}-${scheme}.png`;
}

async function main() {
  mkdirSync(OUT_DIR, { recursive: true });
  const chrome = await launchChrome();
  try {
    for (const route of ROUTES) {
      for (const vp of VIEWPORTS) {
        for (const scheme of SCHEMES) {
          await chrome.open(`${BASE_URL}/#${route}`, { ...vp, scheme });
          const out = screenshotPath(route, vp.width, scheme);
          await chrome.screenshot(out);
          console.log(`wrote ${out}`);
        }
      }
    }
  } finally {
    chrome.close();
  }
}

main().catch((err) => {
  console.error(err.message);
  process.exitCode = 1;
});
