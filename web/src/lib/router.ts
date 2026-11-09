export type Route = 'dashboard' | 'login' | 'setup' | 'wifi' | 'system';

const ROUTES: Record<string, Route> = {
  '/': 'dashboard',
  '/login': 'login',
  '/setup': 'setup',
  '/wifi': 'wifi',
  '/system': 'system',
};

/** Maps a location hash (`#/wifi`) to a route; anything unknown is the dashboard. */
export function parseHash(hash: string): Route {
  const path = hash.replace(/^#/, '').split('?')[0] ?? '';
  return ROUTES[path === '' ? '/' : path] ?? 'dashboard';
}

export function href(route: Route): string {
  return route === 'dashboard' ? '#/' : `#/${route}`;
}
