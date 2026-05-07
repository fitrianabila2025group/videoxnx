// Server-side API helper. Calls the Go backend.
// In the all-in-one container the backend listens on the same host:port as the
// public URL, so an empty NEXT_PUBLIC_API_URL means "same origin" — but during
// SSR the Node process must reach it via 127.0.0.1:8080. Configure with
// INTERNAL_API_URL (server-only) or NEXT_PUBLIC_API_URL (client+server).
const API =
  process.env.INTERNAL_API_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  'http://127.0.0.1:8080';

export async function api<T = any>(
  path: string,
  init: RequestInit = {},
  revalidate = 60
): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...(init.headers || {}) },
    next: { revalidate },
  });
  if (!res.ok) {
    throw new Error(`API ${path} failed: ${res.status}`);
  }
  return res.json();
}

export type Post = {
  id: number;
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  thumbnail_url: string;
  video_embed_url: string;
  source_url: string;
  source_domain: string;
  status: string;
  is_adult: boolean;
  view_count: number;
  published_at?: string;
  scraped_at?: string;
  categories?: { id: number; name: string; slug: string }[];
  tags?: { id: number; name: string; slug: string }[];
};

export type Paginated<T> = {
  data: T[];
  page: number;
  per_page: number;
  total: number;
};
