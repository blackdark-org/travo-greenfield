import { describe, expect, it } from 'vitest';
import { formatBytes, formatLoad, formatUptime } from './format';

describe('formatUptime', () => {
  it.each([
    [0, '<1m'],
    [59, '<1m'],
    [60, '1m'],
    [3600, '1h'],
    [3660, '1h 1m'],
    [90061, '1d 1h'],
    [86400, '1d'],
    [86400 + 60, '1d 1m'],
  ])('%i s -> %s', (s, out) => expect(formatUptime(s)).toBe(out));
});

describe('formatBytes', () => {
  it.each([
    [0, '0 B'],
    [1023, '1023 B'],
    [1024, '1.0 KB'],
    [428646400, '408.8 MB'],
    [81616896, '77.8 MB'],
    [3 * 1024 ** 3, '3.0 GB'],
  ])('%i -> %s', (n, out) => expect(formatBytes(n)).toBe(out));
});

describe('formatLoad', () => {
  it('pads to two decimals', () =>
    expect(formatLoad([0.39, 0.24, 0.1])).toBe('0.39 · 0.24 · 0.10'));
});
