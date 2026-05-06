import { api, Paginated, Post } from '@/lib/api';
import PostCard from '@/components/PostCard';
import Pagination from '@/components/Pagination';

export const metadata = { title: 'Search' };

export default async function Search({ searchParams }: { searchParams: { q?: string; page?: string } }) {
  const q = (searchParams.q || '').trim();
  const page = Number(searchParams.page || '1');
  let posts: Paginated<Post> = { data: [], page, per_page: 24, total: 0 };
  if (q) {
    try { posts = await api(`/api/search?q=${encodeURIComponent(q)}&per_page=24&page=${page}`); } catch {}
  }
  return (
    <div className="max-w-7xl mx-auto px-4 py-6">
      <h1 className="text-2xl font-bold mb-1">Search</h1>
      <p className="text-sub mb-4">Query: <span className="text-ink">{q || '—'}</span></p>
      {q === '' ? (
        <p className="text-sub">Type a query in the search bar above.</p>
      ) : posts.data.length === 0 ? (
        <p className="text-sub">No results.</p>
      ) : (
        <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-4">
          {posts.data.map((p) => <PostCard key={p.id} post={p} />)}
        </div>
      )}
      <Pagination page={posts.page} perPage={posts.per_page} total={posts.total} basePath="/search" query={{ q }} />
    </div>
  );
}
