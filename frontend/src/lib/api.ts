// Server-side API helper. Calls the Go backend.
// Configure NEXT_PUBLIC_API_URL (e.g. http://backend:8080).
const API = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

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
