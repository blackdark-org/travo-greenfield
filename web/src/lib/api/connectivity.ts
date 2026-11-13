import './client';
import { m } from '../paraglide/messages.js';
import { getDashboardConnectivity } from './gen';

export async function fetchConnectivity(signal?: AbortSignal) {
  const { data } = await getDashboardConnectivity({ signal });
  if (!data) throw new Error(m.dashboard_connectivity_unavailable());
  return data;
}
