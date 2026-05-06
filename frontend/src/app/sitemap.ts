import type { MetadataRoute } from 'next';
import { api } from '@/lib/api';
import { getSiteUrl } from '@/lib/site-url';

type ListResp = { data: { slug: string; published_at?: string; scraped_at?: string }[] };

export const revalidate = 600;
export const dynamic = 'force-dynamic';

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const SITE_URL = getSiteUrl();
  const out: MetadataRoute.Sitemap = [
    { url: `${SITE_URL}/`, changeFrequency: 'hourly', priority: 1 },
    { url: `${SITE_URL}/latest`, changeFrequency: 'hourly', priority: 0.9 },
    { url: `${SITE_URL}/trending`, changeFrequency: 'daily', priority: 0.8 },
    { url: `${SITE_URL}/categories`, changeFrequency: 'weekly', priority: 0.6 },
    { url: `${SITE_URL}/tags`, changeFrequency: 'weekly', priority: 0.5 },
    { url: `${SITE_URL}/dmca`, priority: 0.2 },
    { url: `${SITE_URL}/contact`, priority: 0.2 },
    { url: `${SITE_URL}/disclaimer`, priority: 0.2 },
    { url: `${SITE_URL}/privacy`, priority: 0.2 },
    { url: `${SITE_URL}/age-verification`, priority: 0.2 },
  ];

  try {
    // Pull up to 5000 posts in pages of 100.
    let page = 1;
    while (page <= 50) {
      const r = await api<ListResp>(`/api/posts?per_page=100&page=${page}`, {}, 600);
      if (!r.data?.length) break;
      for (const p of r.data) {
        out.push({
          url: `${SITE_URL}/post/${p.slug}`,
          lastModified: p.published_at || p.scraped_at,
          changeFrequency: 'weekly',
          priority: 0.7,
        });
      }
      if (r.data.length < 100) break;
      page++;
    }
  } catch {}

  try {
    const cats = await api<{ data: { slug: string }[] }>(`/api/categories`, {}, 3600);
    for (const c of cats.data || []) {
      out.push({ url: `${SITE_URL}/category/${c.slug}`, changeFrequency: 'daily', priority: 0.6 });
    }
    const tags = await api<{ data: { slug: string }[] }>(`/api/tags`, {}, 3600);
    for (const t of tags.data || []) {
      out.push({ url: `${SITE_URL}/tag/${t.slug}`, changeFrequency: 'weekly', priority: 0.4 });
    }
  } catch {}

  return out;
}
