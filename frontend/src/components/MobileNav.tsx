'use client';
import { useState } from 'react';
import Link from 'next/link';
import SearchBar from './SearchBar';

const links = [
  { href: '/latest', label: 'Terbaru' },
  { href: '/trending', label: 'Trending' },
  { href: '/categories', label: 'Kategori' },
  { href: '/tags', label: 'Tag' },
];

export default function MobileNav() {
  const [open, setOpen] = useState(false);

  return (
    <>
      <button
        type="button"
        aria-label="Menu"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
        className="md:hidden inline-flex items-center justify-center w-10 h-10 rounded-lg bg-muted hover:bg-white/10"
      >
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
          {open ? (
            <>
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </>
          ) : (
            <>
              <line x1="3" y1="6" x2="21" y2="6" />
              <line x1="3" y1="12" x2="21" y2="12" />
              <line x1="3" y1="18" x2="21" y2="18" />
            </>
          )}
        </svg>
      </button>

      {open && (
        <div className="md:hidden absolute top-full left-0 right-0 bg-bg/98 backdrop-blur border-b border-white/5 px-4 py-4 space-y-3 shadow-lg">
          <SearchBar />
          <nav className="grid grid-cols-2 gap-2">
            {links.map((l) => (
              <Link
                key={l.href}
                href={l.href}
                onClick={() => setOpen(false)}
                className="px-3 py-2 rounded-lg bg-muted text-ink text-sm hover:bg-brand/20"
              >
                {l.label}
              </Link>
            ))}
          </nav>
        </div>
      )}
    </>
  );
}
