import { api, Paginated, Post } from '@/lib/api';
import PostCard from '@/components/PostCard';
import Pagination from '@/components/Pagination';
import SearchBar from '@/components/SearchBar';

export const metadata = { title: 'Cari Video' };

export default async function Search({ searchParams }: { searchParams: { q?: string; page?: string } }) {
  const q = (searchParams.q || '').trim();
  const page = Number(searchParams.page || '1');
  let posts: Paginated<Post> = { data: [], page, per_page: 24, total: 0 };
  if (q) {
    try { posts = await api(`/api/search?q=${encodeURIComponent(q)}&per_page=24&page=${page}`); } catch {}
  }
  return (
    <div className="container-x py-4 sm:py-6">
      <h1 className="text-xl sm:text-2xl font-bold mb-3">Cari Video</h1>

      {/* Always-visible search input (especially important on mobile) */}
      <div className="mb-4 sm:mb-6">
        <SearchBar />
      </div>

      {q && (
        <p className="text-sub mb-4 text-sm">
          Hasil untuk: <span className="text-ink font-semibold">&ldquo;{q}&rdquo;</span>
          {posts.total > 0 && <span className="ml-2 text-xs">({posts.total} video)</span>}
        </p>
      )}

      {q === '' ? (
        <div className="text-center py-12 sm:py-16 border border-dashed border-white/10 rounded-xl text-sub">
          <p className="text-base">Ketik kata kunci di kolom pencarian di atas.</p>
          <p className="text-xs mt-1 opacity-70">Contoh: jilbab, tante, abg, viral</p>
        </div>
      ) : posts.data.length === 0 ? (
        <div className="text-center py-12 sm:py-16 border border-dashed border-white/10 rounded-xl text-sub">
          <p>Tidak ada hasil untuk &ldquo;{q}&rdquo;.</p>
        </div>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-2.5 sm:gap-4">
          {posts.data.map((p) => <PostCard key={p.id} post={p} />)}
        </div>
      )}
      <Pagination page={posts.page} perPage={posts.per_page} total={posts.total} basePath="/search" query={{ q }} />
    </div>
  );
}
