export type LoadStatus =
  | { state: 'loading' }
  | { state: 'error'; message: string }
  | { state: 'ready' };

/** Runs `fn`, tracking loading/error/data. A rejection is an error state, never "no data". */
export function createLoader<T>(fn: (signal?: AbortSignal) => Promise<T>) {
  let status = $state<LoadStatus>({ state: 'loading' });
  let data = $state<T>();

  // Only the latest run may write state: a slow earlier attempt must not overwrite a retry.
  let latest = 0;

  async function run(): Promise<void> {
    const mine = ++latest;
    status = { state: 'loading' };
    try {
      const result = await fn();
      if (mine !== latest) return;
      data = result;
      status = { state: 'ready' };
    } catch (e) {
      if (mine !== latest) return;
      status = { state: 'error', message: e instanceof Error ? e.message : String(e) };
    }
  }

  void run();

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
