'use client';
import { useEffect, useState } from 'react';
import { bearer } from '../_auth';

export default function Scraper() {
  const [logs, setLogs] = useState<any[]>([]);
  const [busy, setBusy] = useState(false);
  const load = () => fetch('/api/admin/scraper/logs', { credentials: 'include', headers: bearer() })
    .then((r) => r.json()).then((j) => setLogs(j.data || []));
  useEffect(() => { load(); }, []);
  const run = async () => {
    setBusy(true);
    await fetch('/api/admin/scraper/run', { method: 'POST', credentials: 'include', headers: bearer() });
    setBusy(false);
    setTimeout(load, 1500);
  };
  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Scraper</h1>
        <button className="btn" disabled={busy} onClick={run}>{busy ? 'Starting…' : 'Run scraper now'}</button>
      </div>
      <div className="bg-panel border border-white/5 rounded-xl overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="text-sub border-b border-white/5">
            <tr><th className="text-left p-3">When</th><th>Status</th><th>Saved</th><th>Failed</th><th>Duration</th><th className="text-left">Message</th></tr>
          </thead>
          <tbody>
            {logs.map((l) => (
              <tr key={l.id} className="border-b border-white/5">
                <td className="p-3">{new Date(l.created_at).toLocaleString()}</td>
                <td className="text-center"><span className="chip">{l.status}</span></td>
                <td className="text-center">{l.scraped_count}</td>
                <td className="text-center">{l.failed_count}</td>
                <td className="text-center">{l.duration_ms} ms</td>
                <td className="p-3 text-sub">{l.message}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
