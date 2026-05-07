import { api, Paginated, Post } from '@/lib/api';
import PostCard from '@/components/PostCard';
import Pagination from '@/components/Pagination';

export const metadata = { title: 'Lagi Trending' };

export default async function Trending({ searchParams }: { searchParams: { page?: string } }) {
  const page = Number(searchParams.page || '1');
  let posts: Paginated<Post> = { data: [], page, per_page: 24, total: 0 };
  try { posts = await api(`/api/trending?per_page=24&page=${page}`); } catch {}
  return (
    <div className="container-x py-4 sm:py-6">
      <h1 className="text-2xl font-bold mb-4">Lagi Trending</h1>
      <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-2.5 sm:gap-4">
        {posts.data.map((p) => <PostCard key={p.id} post={p} />)}
      </div>
      <Pagination page={posts.page} perPage={posts.per_page} total={posts.total} basePath="/trending" />
    </div>
  );
}
