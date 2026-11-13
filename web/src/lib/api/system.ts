import { getSystemBackup, getSystemInfo } from './gen';
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

export async function downloadSystemBackup() {
  const { data, response } = await getSystemBackup({ parseAs: 'blob' });
  if (!response?.ok || !(data instanceof Blob) || data.type !== 'application/gzip' || !data.size) {
    throw new Error('Configuration backup failed');
  }

  const url = URL.createObjectURL(data);
  const link = document.createElement('a');
  link.href = url;
  link.download = 'travo-config.tar.gz';
  try {
    link.click();
  } finally {
    setTimeout(() => URL.revokeObjectURL(url), 0);
  }
}
