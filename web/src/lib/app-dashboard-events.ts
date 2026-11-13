import type { dashboardEvents } from './dashboard-events.svelte';

export type DashboardEventOwner = Pick<typeof dashboardEvents, 'start' | 'stop' | 'clear'>;

/** A persisted page keeps its session and values, but its aborted stream needs a new owner. */
export function restoreDashboardEvents(
  persisted: boolean,
  authenticated: boolean,
  generation: number,
  events: DashboardEventOwner,
): void {
  if (persisted && authenticated) events.start(generation);
}
