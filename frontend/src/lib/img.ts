/**
 * Wrap a remote image URL through our same-origin /api/img proxy so that
 * hotlink protection, cross-origin referrer rules, mixed-content blocking
 * (Android Chrome especially) and CORS issues all go away. Local/proxied
 * URLs are returned unchanged.
 */
export function proxyImg(url?: string | null): string {
  if (!url) return '';
  if (url.startsWith('/') || url.startsWith('data:') || url.startsWith('blob:')) return url;
  // Already proxied
  if (url.includes('/api/img?')) return url;
  return `/api/img?u=${encodeURIComponent(url)}`;
}
