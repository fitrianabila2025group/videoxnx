import Image from 'next/image';
import Link from 'next/link';
import type { Post } from '@/lib/api';
import { proxyImg } from '@/lib/img';

export default function PostCard({ post }: { post: Post }) {
  return (
    <Link href={`/post/${post.slug}`} className="card group block focus:outline-none focus:ring-2 focus:ring-brand/40">
      <div className="relative aspect-video bg-black/40 overflow-hidden">
        {post.thumbnail_url ? (
          <Image
            src={proxyImg(post.thumbnail_url)}
            alt={post.title}
            fill
            sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 25vw"
            className="object-cover group-hover:scale-105 transition-transform duration-300"
            unoptimized
          />
        ) : (
          <div className="absolute inset-0 grid place-items-center text-sub text-xs">No preview</div>
        )}
        {/* Bottom gradient + meta badges */}
        <div className="absolute inset-x-0 bottom-0 h-14 bg-gradient-to-t from-black/80 to-transparent pointer-events-none" />
        {post.video_embed_url && (
          <span className="absolute top-1.5 left-1.5 text-[10px] bg-brand/90 text-white px-1.5 py-0.5 rounded font-bold tracking-wide">
            HD
          </span>
        )}
        {/* Play overlay on hover (desktop only) */}
        <div className="hidden md:flex absolute inset-0 items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
          <span className="w-12 h-12 rounded-full bg-brand/90 grid place-items-center shadow-lg">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="white" aria-hidden>
              <path d="M8 5v14l11-7z"/>
            </svg>
          </span>
        </div>
      </div>
      <div className="p-2.5 sm:p-3">
        <h3 className="text-[13px] sm:text-sm font-semibold leading-snug line-clamp-2 group-hover:text-brand transition-colors min-h-[2.5rem]">
          {post.title}
        </h3>
        <div className="mt-1.5 flex items-center gap-2 text-[11px] text-sub">
          {post.view_count > 0 && (
            <span className="inline-flex items-center gap-1">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden>
                <path d="M1 12s4-7 11-7 11 7 11 7-4 7-11 7S1 12 1 12z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
              {formatCount(post.view_count)}
            </span>
          )}
          {post.categories?.[0] && (
            <span className="truncate text-sub/80">· {post.categories[0].name}</span>
          )}
        </div>
      </div>
    </Link>
  );
}

function formatCount(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1).replace(/\.0$/, '') + 'M';
  if (n >= 1_000)     return (n / 1_000).toFixed(1).replace(/\.0$/, '') + 'K';
  return String(n);
}
