import { api, Post } from '@/lib/api';
import Link from 'next/link';
import Image from 'next/image';
import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import PostCard from '@/components/PostCard';
import ReportButton from '@/components/ReportButton';
import VideoPlayer from '@/components/VideoPlayer';
import { proxyImg } from '@/lib/img';

type PostResp = { data: Post; related: Post[]; prev: Post | null; next: Post | null };

export const revalidate = 120;

export async function generateMetadata({ params }: { params: { slug: string } }): Promise<Metadata> {
  try {
    const r = await api<PostResp>(`/api/posts/${params.slug}`);
    return {
      title: r.data.title,
      description: r.data.excerpt?.slice(0, 160),
      openGraph: {
        title: r.data.title,
        description: r.data.excerpt?.slice(0, 160),
        images: r.data.thumbnail_url ? [r.data.thumbnail_url] : [],
        type: 'video.other',
      },
      alternates: { canonical: `/post/${r.data.slug}` },
    };
  } catch {
    return { title: 'Not found' };
  }
}

export default async function PostPage({ params }: { params: { slug: string } }) {
  let resp: PostResp;
  try {
    resp = await api(`/api/posts/${params.slug}`);
  } catch {
    return notFound();
  }
  const p = resp.data;

  const ld = {
    '@context': 'https://schema.org',
    '@type': 'VideoObject',
    name: p.title,
    description: p.excerpt,
    thumbnailUrl: p.thumbnail_url,
    uploadDate: p.published_at || p.scraped_at,
    embedUrl: p.video_embed_url || undefined,
  };

  return (
    <article className="max-w-5xl mx-auto px-4 py-6">
      <h1 className="text-2xl md:text-3xl font-bold">{p.title}</h1>
      <div className="mt-2 text-sub text-sm flex flex-wrap gap-2">
        {p.published_at && <time>{new Date(p.published_at).toLocaleDateString('id-ID')}</time>}
        <span>· {p.view_count}x ditonton</span>
      </div>

      <div className="mt-4 rounded-xl overflow-hidden bg-black aspect-video">
        {p.video_embed_url ? (
          <VideoPlayer src={p.video_embed_url} poster={proxyImg(p.thumbnail_url) || undefined} />
        ) : p.thumbnail_url ? (
          <Image src={proxyImg(p.thumbnail_url)} alt={p.title} width={1280} height={720} className="w-full h-full object-cover" unoptimized />
        ) : null}
      </div>

      {p.excerpt && <p className="mt-4 text-sub">{p.excerpt}</p>}

      {p.content && (
        <div
          className="prose prose-invert max-w-none mt-6 [&_iframe]:w-full [&_iframe]:aspect-video"
          dangerouslySetInnerHTML={{ __html: p.content }}
        />
      )}

      <div className="mt-6 flex flex-wrap gap-2">
        {(p.categories || []).map((c) => (
          <Link key={c.id} href={`/category/${c.slug}`} className="chip">{c.name}</Link>
        ))}
        {(p.tags || []).map((t) => (
          <Link key={t.id} href={`/tag/${t.slug}`} className="chip">#{t.name}</Link>
        ))}
      </div>

      <div className="mt-6 flex flex-wrap gap-3 items-center">
        <ReportButton postId={p.id} />
      </div>

      <div className="mt-8 flex justify-between gap-3">
        {resp.prev ? <Link href={`/post/${resp.prev.slug}`} className="btn-ghost">← {resp.prev.title.slice(0, 40)}…</Link> : <span />}
        {resp.next ? <Link href={`/post/${resp.next.slug}`} className="btn-ghost">{resp.next.title.slice(0, 40)}… →</Link> : <span />}
      </div>

      {resp.related?.length > 0 && (
        <section className="mt-12">
          <h2 className="text-lg font-semibold mb-3">Video Terkait</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {resp.related.map((r) => <PostCard key={r.id} post={r} />)}
          </div>
        </section>
      )}

      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(ld) }} />
    </article>
  );
}
