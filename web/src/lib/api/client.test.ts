import { describe, expect, it } from 'vitest';
import { isAuthEndpoint } from './client';

describe('isAuthEndpoint', () => {
  it.each([
    ['http://x.test/api/v1/auth/login', true],
    ['http://x.test/api/v1/setup/password', true],
    ['/api/v1/auth/login', true],
    ['http://x.test/api/v1/auth/session', false],
    ['http://x.test/api/v1/auth/password', false],
    ['http://x.test/api/v1/auth/logout', false],
    ['http://x.test/api/v1/auth/login/extra', false],
  ])('%s -> %s', (url, want) => {
    expect(isAuthEndpoint(url)).toBe(want);
  });
});
