'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import SearchBar from './SearchBar';

const links = [
  { href: '/latest',     label: 'Terbaru',  icon: '🆕' },
  { href: '/trending',   label: 'Trending', icon: '🔥' },
  { href: '/categories', label: 'Kategori', icon: '📁' },
  { href: '/tags',       label: 'Tag',      icon: '🏷️' },
  { href: '/search',     label: 'Cari',     icon: '🔍' },
];

const footerLinks = [
  { href: '/dmca',             label: 'DMCA' },
  { href: '/contact',          label: 'Kontak' },
  { href: '/disclaimer',       label: 'Disclaimer' },
  { href: '/privacy',          label: 'Privasi' },
  { href: '/age-verification', label: '18+' },
];

export default function MobileNav() {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (typeof document === 'undefined') return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = open ? 'hidden' : prev || '';
    return () => { document.body.style.overflow = prev; };
  }, [open]);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') setOpen(false); };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [open]);

  return (
    <>
      <button
        type="button"
        aria-label={open ? 'Tutup menu' : 'Buka menu'}
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
        className="md:hidden btn-icon"
      >
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
          <line x1="3" y1="6"  x2="21" y2="6" />
          <line x1="3" y1="12" x2="21" y2="12" />
          <line x1="3" y1="18" x2="21" y2="18" />
        </svg>
      </button>

      {open && (
        <div className="md:hidden fixed inset-0 z-50">
          <button
            aria-label="Tutup menu"
            onClick={() => setOpen(false)}
            className="absolute inset-0 bg-black/70 backdrop-blur-sm animate-fade"
          />
          <aside
            role="dialog"
            aria-modal="true"
            aria-label="Menu navigasi"
            className="absolute top-0 right-0 bottom-0 w-[82%] max-w-sm bg-panel border-l border-white/10 shadow-2xl flex flex-col animate-drawer"
          >
            <div className="flex items-center justify-between px-4 py-3 border-b border-white/5">
              <Link href="/" onClick={() => setOpen(false)} className="text-lg font-bold text-brand">
                Video<span className="text-ink">XNX</span>
              </Link>
              <button
                type="button"
                aria-label="Tutup"
                onClick={() => setOpen(false)}
                className="btn-icon"
              >
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden>
                  <line x1="18" y1="6"  x2="6"  y2="18" />
                  <line x1="6"  y1="6"  x2="18" y2="18" />
                </svg>
              </button>
            </div>

            <div className="px-4 pt-4">
              <SearchBar />
            </div>

            <nav className="px-2 py-3 flex-1 overflow-y-auto">
              <ul className="space-y-1">
                {links.map((l) => (
                  <li key={l.href}>
                    <Link
                      href={l.href}
                      onClick={() => setOpen(false)}
                      className="flex items-center gap-3 px-3 py-3 rounded-lg text-ink hover:bg-white/5 active:bg-white/10"
                    >
                      <span className="text-lg w-6 text-center" aria-hidden>{l.icon}</span>
                      <span className="text-base">{l.label}</span>
                    </Link>
                  </li>
                ))}
              </ul>
            </nav>

            <div className="px-4 py-3 border-t border-white/5 flex flex-wrap gap-x-4 gap-y-2 text-xs text-sub">
              {footerLinks.map((l) => (
                <Link
                  key={l.href}
                  href={l.href}
                  onClick={() => setOpen(false)}
                  className="hover:text-ink"
                >
                  {l.label}
                </Link>
              ))}
            </div>
          </aside>
        </div>
      )}
    </>
  );
}
