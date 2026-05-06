'use client';
import { useEffect, useState } from 'react';
import { bearer } from './_auth';

type Item = { id: number; name: string; slug: string };

function CrudPage({ resource, label }: { resource: 'categories' | 'tags'; label: string }) {
  const [items, setItems] = useState<Item[]>([]);
  const [name, setName] = useState('');
  const load = () => fetch(`/api/admin/${resource}`, { credentials: 'include', headers: bearer() })
    .then((r) => r.json()).then((j) => setItems(j.data || []));
  useEffect(() => { load(); }, []);
  const add = async () => {
    if (!name.trim()) return;
    await fetch(`/api/admin/${resource}`, {
      method: 'POST', credentials: 'include',
      headers: { 'Content-Type': 'application/json', ...bearer() },
      body: JSON.stringify({ name }),
    });
    setName(''); load();
  };
  const del = async (id: number) => {
    if (!confirm('Delete?')) return;
    await fetch(`/api/admin/${resource}/${id}`, { method: 'DELETE', credentials: 'include', headers: bearer() });
    load();
  };
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">{label}</h1>
      <div className="flex gap-2 mb-4">
        <input className="input max-w-xs" placeholder="New name" value={name} onChange={(e) => setName(e.target.value)} />
        <button className="btn" onClick={add}>Add</button>
      </div>
      <div className="bg-panel border border-white/5 rounded-xl divide-y divide-white/5">
        {items.map((i) => (
          <div key={i.id} className="flex items-center justify-between p-3">
            <div><span className="font-medium">{i.name}</span> <span className="text-sub text-xs ml-2">/{i.slug}</span></div>
            <button className="btn-ghost text-red-400" onClick={() => del(i.id)}>Delete</button>
          </div>
        ))}
        {items.length === 0 && <div className="p-3 text-sub">Empty.</div>}
      </div>
    </div>
  );
}

export { CrudPage };
