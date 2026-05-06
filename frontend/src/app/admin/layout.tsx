'use client';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { bearer } from './_auth';

const NAV = [
  { href: '/admin', label: 'Dashboard' },
  { href: '/admin/posts', label: 'Posts' },
  { href: '/admin/categories', label: 'Categories' },
  { href: '/admin/tags', label: 'Tags' },
  { href: '/admin/reports', label: 'Reports' },
  { href: '/admin/scraper', label: 'Scraper' },
  { href: '/admin/settings', label: 'Settings' },
];

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const path = usePathname();
  const r = useRouter();
  const [ok, setOk] = useState<boolean | null>(null);

  useEffect(() => {
    if (path === '/admin/login') { setOk(true); return; }
    fetch('/api/admin/dashboard', { credentials: 'include', headers: bearer() })
      .then((res) => {
        if (res.ok) setOk(true);
        else { setOk(false); r.push('/admin/login'); }
      })
      .catch(() => { setOk(false); r.push('/admin/login'); });
  }, [path, r]);

  if (path === '/admin/login') return <>{children}</>;
  if (ok === null) return <div className="p-8 text-sub">Loading admin…</div>;
  if (!ok) return null;

  return (
    <div className="grid grid-cols-[220px_1fr] min-h-[calc(100vh-120px)]">
      <aside className="bg-panel border-r border-white/5 p-4">
        <div className="text-brand font-bold mb-4">Admin</div>
        <nav className="flex flex-col gap-1">
          {NAV.map((n) => (
            <Link
              key={n.href}
              href={n.href}
              className={`px-3 py-2 rounded hover:bg-white/5 ${path === n.href ? 'bg-white/10 text-ink' : 'text-sub'}`}
            >
              {n.label}
            </Link>
          ))}
          <button
            className="text-left px-3 py-2 rounded text-sub hover:bg-white/5 mt-4"
            onClick={async () => {
              await fetch('/api/admin/logout', { method: 'POST', credentials: 'include' });
              try { localStorage.removeItem('admin_token'); } catch {}
              r.push('/admin/login');
            }}
          >
            Logout
          </button>
        </nav>
      </aside>
      <section className="p-6">{children}</section>
    </div>
  );
}


