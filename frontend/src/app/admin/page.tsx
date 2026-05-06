'use client';
import { useEffect, useState } from 'react';
import { bearer } from './_auth';

export default function AdminDashboard() {
  const [data, setData] = useState<any>(null);
  useEffect(() => {
    fetch('/api/admin/dashboard', { credentials: 'include', headers: bearer() })
      .then((r) => r.json()).then(setData).catch(() => {});
  }, []);
  if (!data) return <div className="text-sub">Loading…</div>;
  const c = data.counts || {};
  const cards = [
    ['Total posts', c.Total],
    ['Published', c.Published],
    ['Hidden', c.Hidden],
    ['Blocked', c.Blocked],
    ['Categories', c.Categories],
    ['Tags', c.Tags],
    ['Open reports', c.Reports],
  ];
  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Dashboard</h1>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {cards.map(([k, v]) => (
          <div key={k as string} className="bg-panel border border-white/5 rounded-xl p-4">
            <div className="text-sub text-sm">{k}</div>
            <div className="text-2xl font-bold mt-1">{v ?? 0}</div>
          </div>
        ))}
      </div>
      {data.last_scrape && (
        <div className="mt-8 bg-panel border border-white/5 rounded-xl p-4">
          <h2 className="font-semibold mb-2">Last scrape</h2>
          <pre className="text-xs text-sub overflow-x-auto">{JSON.stringify(data.last_scrape, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}
