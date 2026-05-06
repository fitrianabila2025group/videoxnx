'use client';
import { useEffect, useState } from 'react';
import { bearer } from '../_auth';

const FIELDS = [
  ['site_name', 'Site name'],
  ['logo_url', 'Logo URL'],
  ['meta_title', 'Meta title'],
  ['meta_description', 'Meta description'],
  ['source_website_url', 'Source website URL'],
  ['age_gate_enabled', 'Age gate enabled (true/false)'],
  ['dmca_email', 'DMCA email'],
];

export default function Settings() {
  const [vals, setVals] = useState<Record<string, string>>({});
  useEffect(() => {
    fetch('/api/admin/settings', { credentials: 'include', headers: bearer() })
      .then((r) => r.json()).then((j) => setVals(j.data || {}));
  }, []);
  const save = async () => {
    await fetch('/api/admin/settings', {
      method: 'PUT', credentials: 'include',
      headers: { 'Content-Type': 'application/json', ...bearer() },
      body: JSON.stringify(vals),
    });
    alert('Saved');
  };
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Settings</h1>
      <div className="space-y-3 max-w-xl">
        {FIELDS.map(([k, label]) => (
          <div key={k}>
            <label className="block text-sm text-sub mb-1">{label}</label>
            <input className="input" value={vals[k] || ''} onChange={(e) => setVals({ ...vals, [k]: e.target.value })} />
          </div>
        ))}
        <button className="btn" onClick={save}>Save</button>
      </div>
    </div>
  );
}
