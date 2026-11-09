import { createRawSnippet } from 'svelte';
import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';
import LoadState from './LoadState.svelte';
import { createLoader, type LoadStatus } from './loader.svelte';

const content = createRawSnippet(() => ({ render: () => '<p>CONTENT</p>' }));

function html(status: LoadStatus): string {
  return render(LoadState, { props: { status, retry: () => {}, children: content } }).body;
}

describe('LoadState', () => {
  it('renders a skeleton while loading, not content', () => {
    const out = html({ state: 'loading' });
    expect(out).toContain('skeleton');
    expect(out).not.toContain('CONTENT');
  });

  it('renders the error and a retry button on failure, never content or nothing', () => {
    const out = html({ state: 'error', message: 'boom' });
    expect(out).toContain('role="alert"');
    expect(out).toContain('boom');
    expect(out).toContain('<button');
    expect(out).not.toContain('CONTENT');
  });

  it('renders content when ready', () => {
    expect(html({ state: 'ready' })).toContain('CONTENT');
  });

  it('turns a rejection into an error state and recovers on retry', async () => {
    const fn = vi.fn().mockRejectedValueOnce(new Error('down')).mockResolvedValueOnce('up');
    const loader = createLoader(fn);
    expect(loader.status.state).toBe('loading');
    await vi.waitFor(() => expect(loader.status).toEqual({ state: 'error', message: 'down' }));
    expect(html(loader.status)).toContain('down');
    await loader.retry();
    expect(loader.status.state).toBe('ready');
    expect(loader.data).toBe('up');
  });

  it('ignores a slow earlier run that resolves after a newer retry', async () => {
    let resolveA: (v: string) => void = () => {};
    const fn = vi
      .fn()
      .mockImplementationOnce(() => new Promise<string>((r) => (resolveA = r)))
      .mockResolvedValueOnce('B');
    const loader = createLoader(fn);
    await loader.retry();
    expect(loader.data).toBe('B');
    resolveA('A');
    await Promise.resolve();
    await Promise.resolve();
    expect(loader.data).toBe('B');
    expect(loader.status.state).toBe('ready');
  });

  it('ignores a slow earlier failure after a newer retry succeeded', async () => {
    let rejectA: (e: Error) => void = () => {};
    const fn = vi
      .fn()
      .mockImplementationOnce(() => new Promise<string>((_, r) => (rejectA = r)))
      .mockResolvedValueOnce('B');
    const loader = createLoader(fn);
    await loader.retry();
    rejectA(new Error('late'));
    await Promise.resolve();
    await Promise.resolve();
    expect(loader.status.state).toBe('ready');
    expect(loader.data).toBe('B');
  });
});
