// Helper to attach the JWT bearer header from localStorage in admin pages.
export function bearer(): HeadersInit {
  if (typeof window === 'undefined') return {};
  const t = localStorage.getItem('admin_token');
  return t ? { Authorization: `Bearer ${t}` } : {};
}
