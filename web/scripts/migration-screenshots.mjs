import { mkdirSync } from 'node:fs';
import { launchChrome } from './cdp.mjs';

const base = `http://localhost:${process.env.WEB_PORT ?? 18192}`;
const out = `web/screenshots/${process.env.PHASE ?? 'after'}`;
mkdirSync(out, { recursive: true });
const chrome = await launchChrome();
const pause = () => new Promise((r) => setTimeout(r, 300));
const click = async (text) => {
  await chrome.evaluate(
    `Array.from(document.querySelectorAll('button,a')).find(el => el.textContent.trim() === ${JSON.stringify(text)})?.click()`,
  );
  await pause();
};
try {
  await chrome.open(`${base}/#/login`, { width: 1280, height: 800 });
  await chrome.evaluate(
    `(() => {const e=document.querySelector('input');e.value='mock-password';e.dispatchEvent(new Event('input',{bubbles:true}));document.querySelector('form').requestSubmit()})()`,
  );
  await pause();
  for (const width of [390, 1280])
    for (const scheme of ['light', 'dark']) {
      await chrome.open(`${base}/#/`, { width, height: 844, scheme });
      await chrome.evaluate(`document.documentElement.setAttribute('data-theme','${scheme}')`);
      await chrome.screenshot(`${out}/dashboard-${width}-${scheme}.png`);
      if (width === 390) {
        await chrome.evaluate(`document.querySelector('[aria-controls="sidebar"]').click()`);
        await pause();
        await chrome.screenshot(`${out}/drawer-${width}-${scheme}.png`);
        await chrome.send('Input.dispatchKeyEvent', {
          type: 'keyDown',
          key: 'Escape',
          code: 'Escape',
        });
        await pause();
      }
      await click('Connect to WiFi');
      await pause();
      await click('Scan');
      for (let i = 0; i < 40; i++) {
        if (await chrome.evaluate(`!!document.querySelector('[aria-label="Join Neighbor2"]')`))
          break;
        await pause();
      }
      await chrome.screenshot(`${out}/wifi-${width}-${scheme}.png`);
      await chrome.evaluate(`document.querySelector('[aria-label="Join Neighbor2"]').click()`);
      await pause();
      await chrome.screenshot(`${out}/join-${width}-${scheme}.png`);
      await chrome.send('Input.dispatchKeyEvent', {
        type: 'keyDown',
        key: 'Escape',
        code: 'Escape',
      });
      await pause();
      await chrome.open(`${base}/#/system`, { width, height: 844, scheme });
      await click('Restart router');
      await chrome.evaluate(
        `Array.from(document.querySelectorAll('button')).find(el=>el.textContent.trim()==='Restart')?.scrollIntoView({block:'center'})`,
      );
      await chrome.screenshot(`${out}/reboot-${width}-${scheme}.png`);
    }
} finally {
  chrome.close();
}
