import { describe, expect, it } from 'vitest';
import { href, parseHash } from './router';

describe('router', () => {
  it.each([
    ['', 'dashboard'],
    ['#/', 'dashboard'],
    ['#/login', 'login'],
    ['#/setup', 'setup'],
    ['#/wifi', 'wifi'],
    ['#/system', 'system'],
    ['#/system?tab=x', 'system'],
  ])('maps %j to %s', (hash, route) => {
    expect(parseHash(hash)).toBe(route);
  });

  it('falls back to the dashboard for unknown hashes', () => {
    expect(parseHash('#/nope')).toBe('dashboard');
    expect(parseHash('#/system/extra')).toBe('dashboard');
  });

  it('builds hrefs that parse back to the same route', () => {
    for (const r of ['dashboard', 'login', 'setup', 'wifi', 'system'] as const) {
      expect(parseHash(href(r))).toBe(r);
    }
  });
});
