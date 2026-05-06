import Image from 'next/image';
import Link from 'next/link';
import type { Post } from '@/lib/api';
import { proxyImg } from '@/lib/img';

export default function PostCard({ post }: { post: Post }) {
  return (
    <Link href={`/post/${post.slug}`} className="card group block">
      <div className="relative aspect-video bg-black/40">
        {post.thumbnail_url ? (
          <Image
            src={proxyImg(post.thumbnail_url)}
            alt={post.title}
            fill
            sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 25vw"
            className="object-cover group-hover:scale-105 transition-transform"
            unoptimized
          />
        ) : (
          <div className="absolute inset-0 grid place-items-center text-sub">No preview</div>
        )}
        {post.video_embed_url && (
          <span className="absolute bottom-2 right-2 text-xs bg-black/70 text-ink px-2 py-0.5 rounded">
            VIDEO
          </span>
        )}
      </div>
      <div className="p-3">
        <h3 className="text-sm font-semibold line-clamp-2 group-hover:text-brand">
          {post.title}
        </h3>
        <div className="mt-2 flex flex-wrap gap-1">
          {(post.categories || []).slice(0, 2).map((c) => (
            <span key={c.id} className="chip">{c.name}</span>
          ))}
        </div>
      </div>
    </Link>
  );
}
