import { describe, expect, it, vi } from 'vitest';
import { restoreDashboardEvents } from './app-dashboard-events';

function eventOwner() {
  return { start: vi.fn(), stop: vi.fn(), clear: vi.fn() };
}

describe('App dashboard bfcache ownership', () => {
  it('restarts the current generation only for an authenticated persisted restore', () => {
    const events = eventOwner();

    restoreDashboardEvents(false, true, 4, events);
    restoreDashboardEvents(true, false, 4, events);
    expect(events.start).not.toHaveBeenCalled();

    restoreDashboardEvents(true, true, 4, events);
    expect(events.start).toHaveBeenCalledOnce();
    expect(events.start).toHaveBeenCalledWith(4);
  });
});
