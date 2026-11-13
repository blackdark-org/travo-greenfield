import { onMount } from 'svelte';

export type LoadStatus =
  | { state: 'loading' }
  | { state: 'error'; message: string }
  | { state: 'ready' };

/** Runs `fn`, tracking loading/error/data. A rejection is an error state, never "no data". */
export function createLoader<T>(fn: (signal?: AbortSignal) => Promise<T>, interval?: number) {
  let status = $state<LoadStatus>({ state: 'loading' });
  let data = $state<T>();

  // Only the latest run may write state: a slow earlier attempt must not overwrite a retry.
  let latest = 0;

  let pending = false;
  let stopped = false;
  const controller = interval === undefined ? undefined : new AbortController();

  async function run(background = false): Promise<void> {
    if (stopped || (interval !== undefined && pending)) return;
    pending = true;
    const mine = ++latest;
    if (!background) status = { state: 'loading' };
    try {
      const result = await fn(controller?.signal);
      if (mine !== latest) return;
      data = result;
      status = { state: 'ready' };
    } catch (e) {
      if (mine !== latest) return;
      status = { state: 'error', message: e instanceof Error ? e.message : String(e) };
    } finally {
      if (mine === latest) pending = false;
    }
  }

  if (interval === undefined) {
    void run();
  } else {
    onMount(() => {
      void run();
      const timer = setInterval(() => void run(true), interval);
      return () => {
        clearInterval(timer);
        stopped = true;
        latest++;
        controller?.abort();
      };
    });
  }

  return {
    get status() {
      return status;
    },
    get data() {
      return data;
    },
    retry: () => run(),
  };
}
