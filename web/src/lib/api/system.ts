import { getSystemInfo } from './gen';
import './client';

/** Resolves with the system body; rejects on transport, HTTP or shape errors so callers never see "absent". */
export async function fetchSystemInfo() {
  const { data, error, response } = await getSystemInfo();
  if (data) return data;
  const detail = (error as { error?: string; detail?: string; title?: string } | undefined) ?? {};
  throw new Error(
    detail.error ?? detail.detail ?? detail.title ?? `HTTP ${response?.status ?? 'error'}`,
  );
}
