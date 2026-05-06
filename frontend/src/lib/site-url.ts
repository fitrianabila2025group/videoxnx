import { headers } from 'next/headers';

/**
 * Resolve the public site URL from the incoming request headers (works in Codespaces,
 * behind reverse proxies, custom domains, etc). Falls back to NEXT_PUBLIC_SITE_URL,
 * then to http://localhost:3000.
 */
export function getSiteUrl(): string {
  try {
    const h = headers();
    const host = h.get('x-forwarded-host') || h.get('host');
    if (host) {
      const proto =
        h.get('x-forwarded-proto') ||
        (host.includes('localhost') || host.startsWith('127.') ? 'http' : 'https');
      return `${proto}://${host}`;
    }
  } catch {
    // headers() may be unavailable during build
  }
  return process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3000';
}
