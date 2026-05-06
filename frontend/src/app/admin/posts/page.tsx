'use client';
import { useEffect, useState } from 'react';
import { bearer } from '../_auth';

type Post = {
  id: number; title: string; slug: string; status: string;
  safety_status: string; thumbnail_url: string; source_url: string;
};

export default function AdminPosts() {
  const [items, setItems] = useState<Post[]>([]);
  const [q, setQ] = useState('');
  const [status, setStatus] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const per = 20;

  const load = () => {
    const params = new URLSearchParams({ page: String(page), per_page: String(per) });
    if (q) params.set('q', q);
    if (status) params.set('status', status);
    fetch(`/api/admin/posts?${params}`, { credentials: 'include', headers: bearer() })
      .then((r) => r.json())
      .then((j) => { setItems(j.data || []); setTotal(j.total || 0); });
  };
  useEffect(() => { load(); }, [page, status]);

  const setPostStatus = async (id: number, s: string) => {
    await fetch(`/api/admin/posts/${id}/status`, {
      method: 'PATCH', credentials: 'include',
      headers: { 'Content-Type': 'application/json', ...bearer() },
      body: JSON.stringify({ status: s }),
    });
    load();
  };
  const del = async (id: number) => {
    if (!confirm('Delete this post permanently?')) return;
    await fetch(`/api/admin/posts/${id}`, { method: 'DELETE', credentials: 'include', headers: bearer() });
    load();
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Posts</h1>
      <div className="flex flex-wrap gap-2 mb-4">
        <input className="input max-w-xs" placeholder="Search…" value={q} onChange={(e) => setQ(e.target.value)} />
        <select className="input max-w-xs" value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">All status</option>
          <option value="published">Published</option>
          <option value="hidden">Hidden</option>
          <option value="blocked">Blocked</option>
          <option value="draft">Draft</option>
        </select>
        <button className="btn" onClick={() => { setPage(1); load(); }}>Search</button>
      </div>

      <div className="bg-panel border border-white/5 rounded-xl overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="text-sub border-b border-white/5">
            <tr><th className="text-left p-3">Title</th><th>Status</th><th>Safety</th><th>Actions</th></tr>
          </thead>
          <tbody>
            {items.map((p) => (
              <tr key={p.id} className="border-b border-white/5">
                <td className="p-3">
                  <a href={`/post/${p.slug}`} target="_blank" className="hover:text-brand">{p.title}</a>
                </td>
                <td className="text-center"><span className="chip">{p.status}</span></td>
                <td className="text-center"><span className="chip">{p.safety_status}</span></td>
                <td className="text-center">
                  <div className="flex gap-1 justify-center flex-wrap">
                    {p.status !== 'published' && <button className="btn-ghost" onClick={() => setPostStatus(p.id, 'published')}>Show</button>}
                    {p.status !== 'hidden' && <button className="btn-ghost" onClick={() => setPostStatus(p.id, 'hidden')}>Hide</button>}
                    {p.status !== 'blocked' && <button className="btn-ghost" onClick={() => setPostStatus(p.id, 'blocked')}>Block</button>}
                    <button className="btn-ghost text-red-400" onClick={() => del(p.id)}>Delete</button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Pager page={page} total={total} per={per} onChange={setPage} />
    </div>
  );
}

function Pager({ page, total, per, onChange }: { page: number; total: number; per: number; onChange: (p: number) => void }) {
  const tp = Math.max(1, Math.ceil(total / per));
  return (
    <div className="flex justify-center gap-2 mt-4">
      <button className="btn-ghost" disabled={page <= 1} onClick={() => onChange(page - 1)}>Prev</button>
      <span className="text-sub text-sm self-center">{page} / {tp}</span>
      <button className="btn-ghost" disabled={page >= tp} onClick={() => onChange(page + 1)}>Next</button>
    </div>
  );
}
