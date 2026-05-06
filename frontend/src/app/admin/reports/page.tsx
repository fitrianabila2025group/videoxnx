'use client';
import { useEffect, useState } from 'react';
import { bearer } from '../_auth';

export default function Reports() {
  const [items, setItems] = useState<any[]>([]);
  const load = () => fetch('/api/admin/reports', { credentials: 'include', headers: bearer() })
    .then((r) => r.json()).then((j) => setItems(j.data || []));
  useEffect(() => { load(); }, []);
  const update = async (id: number, status: string) => {
    await fetch(`/api/admin/reports/${id}`, {
      method: 'PATCH', credentials: 'include',
      headers: { 'Content-Type': 'application/json', ...bearer() },
      body: JSON.stringify({ status }),
    });
    load();
  };
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Reports</h1>
      <div className="bg-panel border border-white/5 rounded-xl divide-y divide-white/5">
        {items.map((r) => (
          <div key={r.id} className="p-3">
            <div className="flex justify-between gap-2">
              <div>
                <div><span className="chip mr-2">{r.status}</span> Post #{r.post_id} — <b>{r.reason}</b></div>
                <div className="text-sub text-sm mt-1">{r.email} — {new Date(r.created_at).toLocaleString()}</div>
                <div className="mt-1">{r.message}</div>
              </div>
              <div className="flex gap-1">
                <button className="btn-ghost" onClick={() => update(r.id, 'reviewing')}>Review</button>
                <button className="btn-ghost" onClick={() => update(r.id, 'closed')}>Close</button>
              </div>
            </div>
          </div>
        ))}
        {items.length === 0 && <div className="p-3 text-sub">No reports.</div>}
      </div>
    </div>
  );
}
