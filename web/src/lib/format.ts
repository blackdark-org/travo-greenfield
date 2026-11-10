/** Two largest non-zero units of d/h/m: 90061 -> "1d 1h"; under a minute -> "<1m". */
export function formatUptime(seconds: number): string {
  const total = Math.max(0, Math.floor(seconds));
  if (total < 60) return '<1m';
  const parts = [
    [Math.floor(total / 86400), 'd'],
    [Math.floor((total % 86400) / 3600), 'h'],
    [Math.floor((total % 3600) / 60), 'm'],
  ] as const;
  return parts
    .filter(([n]) => n > 0)
    .slice(0, 2)
    .map(([n, unit]) => `${n}${unit}`)
    .join(' ');
}

/** Binary (1024) units with SI-style labels: 428646400 -> "408.8 MB". */
export function formatBytes(n: number): string {
  if (n < 1024) return `${n} B`;
  const units = ['KB', 'MB', 'GB'];
  let value = n / 1024;
  let i = 0;
  while (value >= 1024 && i < units.length - 1) {
    value /= 1024;
    i++;
  }
  return `${value.toFixed(1)} ${units[i]}`;
}

/** Load averages as shown to people: "0.39 · 0.24 · 0.10". */
export function formatLoad(load: readonly number[]): string {
  return load.map((l) => l.toFixed(2)).join(' · ');
}
