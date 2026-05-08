'use client';
import { useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

const links = [
  { href: '/',           label: 'Beranda',  icon: '🏠' },
  { href: '/latest',     label: 'Terbaru',  icon: '🆕' },
  { href: '/trending',   label: 'Trending', icon: '🔥' },
  { href: '/categories', label: 'Kategori', icon: '📁' },
  { href: '/tags',       label: 'Tag',      icon: '🏷️' },
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
  const [q, setQ] = useState('');
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement | null>(null);

  // Allow other components (e.g. BottomNav) to open the drawer with search focus
  useEffect(() => {
    const onOpen = (e: Event) => {
      setOpen(true);
      const detail = (e as CustomEvent).detail as { focusSearch?: boolean } | undefined;
      if (detail?.focusSearch) {
        setTimeout(() => inputRef.current?.focus(), 80);
      }
    };
    window.addEventListener('vxnx:open-menu', onOpen as EventListener);
    return () => window.removeEventListener('vxnx:open-menu', onOpen as EventListener);
  }, []);

  // Lock body scroll while drawer is open
  useEffect(() => {
    if (typeof document === 'undefined') return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = open ? 'hidden' : prev || '';
    return () => { document.body.style.overflow = prev; };
  }, [open]);

  // Close with ESC
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') setOpen(false); };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [open]);

  const submitSearch = (e: React.FormEvent) => {
    e.preventDefault();
    const v = q.trim();
    if (!v) return;
    setOpen(false);
    router.push(`/search?q=${encodeURIComponent(v)}`);
    setQ('');
  };

  return (
    <>
      <button
        type="button"
        aria-label="Buka menu"
        aria-expanded={open}
        onClick={() => setOpen(true)}
        className="md:hidden btn-icon"
      >
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
          <line x1="3" y1="6"  x2="21" y2="6" />
          <line x1="3" y1="12" x2="21" y2="12" />
          <line x1="3" y1="18" x2="21" y2="18" />
        </svg>
      </button>

      {open && (
        <div
          className="md:hidden fixed inset-0 z-[60]"
          style={{ height: '100dvh' }}
        >
          {/* Backdrop */}
          <button
            type="button"
            aria-label="Tutup menu"
            onClick={() => setOpen(false)}
            className="absolute inset-0 bg-black/80 animate-fade"
          />
          {/* Drawer */}
          <aside
            role="dialog"
            aria-modal="true"
            aria-label="Menu navigasi"
            className="absolute top-0 right-0 w-[85%] max-w-sm flex flex-col border-l border-white/10 shadow-2xl animate-drawer"
            style={{
              height: '100dvh',
              backgroundColor: '#13131c',
            }}
          >
            <div className="flex items-center justify-between px-4 h-14 border-b border-white/10 shrink-0">
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

            <form onSubmit={submitSearch} className="px-4 pt-4 shrink-0">
              <div className="relative">
                <span className="absolute inset-y-0 left-3 flex items-center text-sub pointer-events-none">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
                    <circle cx="11" cy="11" r="7" />
                    <path d="M21 21l-4.3-4.3" />
                  </svg>
                </span>
                <input
                  ref={inputRef}
                  type="search"
                  inputMode="search"
                  enterKeyHint="search"
                  aria-label="Cari video"
                  className="input pl-11 pr-3"
                  placeholder="Cari video bokep..."
                  value={q}
                  onChange={(e) => setQ(e.target.value)}
                />
              </div>
            </form>

            <nav className="px-2 py-3 flex-1 min-h-0 overflow-y-auto">
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

            <div
              className="px-4 py-3 border-t border-white/10 flex flex-wrap gap-x-4 gap-y-2 text-xs text-sub shrink-0"
              style={{ paddingBottom: 'calc(0.75rem + var(--safe-bottom))' }}
            >
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
