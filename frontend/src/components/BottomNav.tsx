'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const items = [
  { href: '/',           label: 'Home',     icon: HomeIcon },
  { href: '/latest',     label: 'Terbaru',  icon: ClockIcon },
  { href: '/trending',   label: 'Hot',      icon: FlameIcon },
  { href: '/categories', label: 'Kategori', icon: GridIcon },
  { href: '/search',     label: 'Cari',     icon: SearchIcon },
];

export default function BottomNav() {
  const path = usePathname() || '/';
  return (
    <nav
      aria-label="Navigasi bawah"
      className="md:hidden fixed bottom-0 inset-x-0 z-40 bg-panel/95 backdrop-blur border-t border-white/10"
      style={{ paddingBottom: 'var(--safe-bottom)' }}
    >
      <ul className="grid grid-cols-5 h-14">
        {items.map((it) => {
          const active =
            it.href === '/' ? path === '/' : path === it.href || path.startsWith(it.href + '/');
          const Icon = it.icon;
          return (
            <li key={it.href}>
              <Link
                href={it.href}
                aria-current={active ? 'page' : undefined}
                className={`flex flex-col items-center justify-center h-full gap-0.5 text-[11px] transition-colors ${
                  active ? 'text-brand' : 'text-sub hover:text-ink active:text-ink'
                }`}
              >
                <Icon className="w-5 h-5" />
                <span>{it.label}</span>
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}

function svg(props: React.SVGProps<SVGSVGElement>) {
  return { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', strokeWidth: 2, strokeLinecap: 'round' as const, strokeLinejoin: 'round' as const, ...props };
}
function HomeIcon(p: React.SVGProps<SVGSVGElement>) {
  return (<svg {...svg(p)}><path d="M3 12l9-9 9 9"/><path d="M5 10v10h14V10"/></svg>);
}
function ClockIcon(p: React.SVGProps<SVGSVGElement>) {
  return (<svg {...svg(p)}><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg>);
}
function FlameIcon(p: React.SVGProps<SVGSVGElement>) {
  return (<svg {...svg(p)}><path d="M12 2s4 4 4 8a4 4 0 11-8 0c0-2 1-3 1-3-2 1-4 4-4 7a7 7 0 0014 0c0-6-7-12-7-12z"/></svg>);
}
function GridIcon(p: React.SVGProps<SVGSVGElement>) {
  return (<svg {...svg(p)}><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>);
}
function SearchIcon(p: React.SVGProps<SVGSVGElement>) {
  return (<svg {...svg(p)}><circle cx="11" cy="11" r="7"/><path d="M21 21l-4.3-4.3"/></svg>);
}
