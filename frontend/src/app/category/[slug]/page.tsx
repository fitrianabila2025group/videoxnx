import { api, Paginated, Post } from '@/lib/api';
import PostCard from '@/components/PostCard';
import Pagination from '@/components/Pagination';

type Resp = Paginated<Post> & { category: { name: string; slug: string } };

export async function generateMetadata({ params }: { params: { slug: string } }) {
  return { title: `Category: ${params.slug}` };
}

export default async function CategoryPage({ params, searchParams }: { params: { slug: string }; searchParams: { page?: string } }) {
  const page = Number(searchParams.page || '1');
  let r: Resp = { data: [], page, per_page: 24, total: 0, category: { name: params.slug, slug: params.slug } };
  try { r = await api(`/api/categories/${params.slug}/posts?per_page=24&page=${page}`); } catch {}
  return (
    <div className="container-x py-4 sm:py-6">
      <h1 className="text-2xl font-bold mb-4">Category: {r.category?.name || params.slug}</h1>
      <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-2.5 sm:gap-4">
        {r.data.map((p) => <PostCard key={p.id} post={p} />)}
      </div>
      <Pagination page={r.page} perPage={r.per_page} total={r.total} basePath={`/category/${params.slug}`} />
    </div>
  );
}
