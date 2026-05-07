import './../styles/globals.css';
import type { Metadata, Viewport } from 'next';
import Link from 'next/link';
import AgeGate from '@/components/AgeGate';
import SearchBar from '@/components/SearchBar';
import MobileNav from '@/components/MobileNav';
import BottomNav from '@/components/BottomNav';
import { getSiteUrl } from '@/lib/site-url';

const SITE_NAME = 'VideoXNX';
const SITE_DESC =
  'Nonton video bokep Indonesia terbaru: jilbab, tante, abg, janda, viral. Update tiap hari, streaming HD, gratis tanpa ribet.';

export async function generateMetadata(): Promise<Metadata> {
  const SITE_URL = getSiteUrl();
  return {
    metadataBase: new URL(SITE_URL),
    title: {
      default: 'VideoXNX — Nonton Bokep Indo Terbaru, Jilbab, Tante & Abg Viral',
      template: '%s | VideoXNX',
    },
    description: SITE_DESC,
    applicationName: SITE_NAME,
    generator: SITE_NAME,
    keywords: [
      'bokep indo', 'bokep indonesia', 'bokep terbaru', 'bokep viral',
      'video bokep', 'bokep jilbab', 'bokep tante', 'bokep abg',
      'nonton bokep', 'streaming bokep', 'videoxnx',
    ],
    alternates: { canonical: '/' },
    manifest: '/site.webmanifest',
    robots: {
      index: true,
      follow: true,
      googleBot: { index: true, follow: true, 'max-image-preview': 'large', 'max-snippet': -1, 'max-video-preview': -1 },
    },
    openGraph: {
      type: 'website',
      siteName: SITE_NAME,
      locale: 'id_ID',
      url: SITE_URL,
      title: 'VideoXNX — Nonton Bokep Indo Terbaru',
      description: SITE_DESC,
      images: [{ url: '/icon-512.png', width: 512, height: 512, alt: SITE_NAME }],
    },
    twitter: {
      card: 'summary_large_image',
      title: 'VideoXNX — Nonton Bokep Indo Terbaru',
      description: SITE_DESC,
      images: ['/icon-512.png'],
    },
    category: 'entertainment',
    other: { rating: 'adult', 'rating-content': 'RTA-5042-1996-1400-1577-RTA' },
  };
}

export const viewport: Viewport = {
  themeColor: '#0b0b10',
  width: 'device-width',
  initialScale: 1,
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const SITE_URL = getSiteUrl();
  const ld = {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    name: SITE_NAME,
    url: SITE_URL,
    inLanguage: 'id-ID',
    description: SITE_DESC,
    potentialAction: {
      '@type': 'SearchAction',
      target: `${SITE_URL}/search?q={search_term_string}`,
      'query-input': 'required name=search_term_string',
    },
  };
  return (
    <html lang="id">
      <head>
        <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(ld) }} />
      </head>
      <body className="min-h-screen flex flex-col">
        <AgeGate />
        <header className="sticky top-0 z-30 bg-bg/90 backdrop-blur-md border-b border-white/5">
          <div className="container-x py-2.5 sm:py-3 flex items-center gap-3">
            <Link
              href="/"
              className="text-lg sm:text-xl font-extrabold text-brand tracking-tight whitespace-nowrap shrink-0"
            >
              Video<span className="text-ink">XNX</span>
            </Link>
            <nav className="hidden md:flex items-center gap-1 text-sm">
              {[
                { href: '/latest', label: 'Terbaru' },
                { href: '/trending', label: 'Trending' },
                { href: '/categories', label: 'Kategori' },
                { href: '/tags', label: 'Tag' },
              ].map((l) => (
                <Link
                  key={l.href}
                  href={l.href}
                  className="px-3 py-1.5 rounded-lg text-sub hover:text-ink hover:bg-white/5 transition-colors"
                >
                  {l.label}
                </Link>
              ))}
            </nav>
            <div className="hidden sm:block flex-1 max-w-md ml-auto">
              <SearchBar />
            </div>
            <div className="ml-auto sm:ml-0">
              <MobileNav />
            </div>
          </div>
        </header>

        <main className="flex-1">{children}</main>

        <footer className="border-t border-white/5 mt-12 hidden md:block">
          <div className="container-x py-8 text-sm text-sub flex flex-wrap gap-4 justify-between">
            <div>
              <div className="text-ink font-semibold">VideoXNX</div>
              <p className="max-w-md">
                Situs nonton video bokep Indonesia terbaru. Khusus pengunjung 18 tahun ke atas.
              </p>
            </div>
            <div className="flex gap-4 flex-wrap">
              <Link href="/dmca" className="hover:text-ink">DMCA</Link>
              <Link href="/contact" className="hover:text-ink">Kontak</Link>
              <Link href="/disclaimer" className="hover:text-ink">Disclaimer</Link>
              <Link href="/privacy" className="hover:text-ink">Privasi</Link>
              <Link href="/age-verification" className="hover:text-ink">18+</Link>
            </div>
          </div>
        </footer>

        <BottomNav />
      </body>
    </html>
  );
}
