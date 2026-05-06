import { api } from '@/lib/api';
import Link from 'next/link';

export const metadata = { title: 'Tags' };

export default async function Tags() {
  let r: { data: { id: number; name: string; slug: string }[] } = { data: [] };
  try { r = await api('/api/tags'); } catch {}
  return (
    <div className="max-w-7xl mx-auto px-4 py-6">
      <h1 className="text-2xl font-bold mb-4">All Tags</h1>
      <div className="flex flex-wrap gap-2">
        {r.data.map((t) => (
          <Link key={t.id} href={`/tag/${t.slug}`} className="chip">#{t.name}</Link>
        ))}
      </div>
    </div>
  );
}
