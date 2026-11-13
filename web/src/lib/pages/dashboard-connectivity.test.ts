import { render } from 'svelte/server';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { fetchConnectivity } from '../api/connectivity';
import type { Status } from '../api/gen';
import type { LoadStatus } from '../ui/loader.svelte';
import Dashboard from './Dashboard.svelte';

const { createLoader } = vi.hoisted(() => ({ createLoader: vi.fn() }));
vi.mock('../ui/loader.svelte', () => ({ createLoader }));

function dashboard(status: LoadStatus, data: Status) {
  createLoader.mockImplementation((fetch: unknown) => ({
    status: fetch === fetchConnectivity ? status : { state: 'ready' },
    data: fetch === fetchConnectivity ? data : undefined,
    retry: vi.fn(),
  }));
  return render(Dashboard).body;
}

describe('Dashboard connectivity', () => {
  beforeEach(() => {
    vi.stubGlobal('document', { documentElement: { getAttribute: () => null } });
  });
  afterEach(() => vi.unstubAllGlobals());

  it('shows active uplink and link up even when the internet is offline', () => {
    const body = dashboard(
      { state: 'ready' },
      { active_interface: 'phy0-sta0', carrier: true, internet: 'offline' },
    );
    expect(body).toContain('phy0-sta0');
    expect(body).toContain('Link up');
    expect(body).toContain('Offline');
    expect(createLoader).toHaveBeenCalledWith(fetchConnectivity);
  });

  it('distinguishes unknown and down carrier', () => {
    const data = { active_interface: null, carrier: null, internet: 'online' } as const;
    const body = dashboard({ state: 'ready' }, data);
    expect(body).toContain('Unknown');
    expect(body).toContain('Online');
    expect(body).not.toContain('Link down');
    expect(dashboard({ state: 'ready' }, { ...data, carrier: false })).toContain('Link down');
  });

  it('retains the last connectivity snapshot on failed refresh', () => {
    const body = dashboard(
      { state: 'error', message: 'Connectivity unavailable' },
      { active_interface: 'phy0-sta0', carrier: true, internet: 'captive' },
    );
    expect(body).toContain('Last known value');
    expect(body).toContain('phy0-sta0');
    expect(body).toContain('Captive portal');
    expect(body).toContain('Connectivity unavailable');
  });
});
