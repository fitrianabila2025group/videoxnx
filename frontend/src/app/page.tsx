import { api, Paginated, Post } from '@/lib/api';
import PostCard from '@/components/PostCard';
import Pagination from '@/components/Pagination';
import Link from 'next/link';

export const revalidate = 60;

export default async function Home({ searchParams }: { searchParams: { page?: string } }) {
  const page = Number(searchParams.page || '1');
  let posts: Paginated<Post> = { data: [], page, per_page: 24, total: 0 };
  let cats: { data: { id: number; name: string; slug: string }[] } = { data: [] };
  let tags: { data: { id: number; name: string; slug: string }[] } = { data: [] };

  try {
    [posts, cats, tags] = await Promise.all([
      api<Paginated<Post>>(`/api/posts?per_page=24&page=${page}`),
      api(`/api/categories`),
      api(`/api/tags`),
    ]);
  } catch {}

  return (
    <div className="container-x py-4 sm:py-6">
      <section className="mb-6 sm:mb-8 rounded-2xl bg-gradient-to-br from-brand/40 via-panel to-panel p-5 sm:p-8 border border-white/5 relative overflow-hidden">
        <div className="absolute -right-10 -top-10 w-48 h-48 rounded-full bg-brand/20 blur-3xl pointer-events-none" />
        <h1 className="relative text-2xl sm:text-3xl md:text-4xl font-bold leading-tight">
          Nonton Bokep Indo Terbaru — VideoXNX
        </h1>
        <p className="relative mt-2 text-sm sm:text-base text-sub max-w-2xl">
          Koleksi video bokep Indonesia, jilbab, tante, abg, dan janda viral
          terbaru. Update setiap hari, full HD, streaming cepat tanpa ribet.
        </p>
        <div className="relative mt-4 flex flex-wrap gap-2 sm:gap-3">
          <Link href="/latest" className="btn">Tonton Terbaru</Link>
          <Link href="/trending" className="btn-ghost">Lagi Trending</Link>
        </div>
      </section>

      <div className="grid lg:grid-cols-[1fr_320px] gap-4 sm:gap-6">
        <div>
          <div className="flex items-end justify-between mb-3 sm:mb-4">
            <h2 className="text-lg sm:text-xl font-semibold">Video Terbaru</h2>
            <Link href="/latest" className="text-xs sm:text-sm text-sub hover:text-brand">Lihat semua →</Link>
          </div>
          {posts.data.length === 0 ? (
            <EmptyState />
          ) : (
            <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-2.5 sm:gap-4">
              {posts.data.map((p) => <PostCard key={p.id} post={p} />)}
            </div>
          )}
          <Pagination page={posts.page} perPage={posts.per_page} total={posts.total} basePath="/" />
        </div>

        <aside className="space-y-4 sm:space-y-6">
          <Section title="Kategori">
            <div className="flex flex-wrap gap-2">
              {cats.data.slice(0, 30).map((c) => (
                <Link key={c.id} href={`/category/${c.slug}`} className="chip">{c.name}</Link>
              ))}
            </div>
          </Section>
          <Section title="Tag Populer">
            <div className="flex flex-wrap gap-2">
              {tags.data.slice(0, 40).map((t) => (
                <Link key={t.id} href={`/tag/${t.slug}`} className="chip">#{t.name}</Link>
              ))}
            </div>
          </Section>
        </aside>
      </div>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-panel border border-white/5 rounded-xl p-4">
      <h3 className="font-semibold mb-3">{title}</h3>
      {children}
    </div>
  );
}

function EmptyState() {
  return (
    <div className="text-center text-sub py-16 border border-dashed border-white/10 rounded-xl">
      <p>Belum ada video. Silakan kembali sebentar lagi.</p>
    </div>
  );
}
