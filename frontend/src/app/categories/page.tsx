import { api } from '@/lib/api';
import Link from 'next/link';

export const metadata = { title: 'Categories' };

export default async function Categories() {
  let r: { data: { id: number; name: string; slug: string }[] } = { data: [] };
  try { r = await api('/api/categories'); } catch {}
  return (
    <div className="container-x py-4 sm:py-6">
      <h1 className="text-2xl font-bold mb-4">All Categories</h1>
      <div className="flex flex-wrap gap-2">
        {r.data.map((c) => (
          <Link key={c.id} href={`/category/${c.slug}`} className="chip">{c.name}</Link>
        ))}
      </div>
    </div>
  );
}
