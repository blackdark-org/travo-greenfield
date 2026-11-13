import './client';
import { m } from '../paraglide/messages.js';
import { getDashboardThroughput } from './gen';

export async function fetchThroughput(signal?: AbortSignal) {
  const { data } = await getDashboardThroughput({ signal });
  if (!data) throw new Error(m.dashboard_throughput_unavailable());
  return data;
}
