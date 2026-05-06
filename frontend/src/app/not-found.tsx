import Link from 'next/link';
import { api, Paginated, Post } from '@/lib/api';
import PostCard from '@/components/PostCard';

export const metadata = {
  title: 'Halaman tidak ditemukan',
  robots: { index: false, follow: false },
};

export default async function NotFound() {
  let posts: Paginated<Post> | null = null;
  try { posts = await api<Paginated<Post>>(`/api/posts?per_page=8&page=1`); } catch {}

  return (
    <div className="max-w-5xl mx-auto px-4 py-16">
      <div className="text-center">
        <p className="text-brand text-sm tracking-widest font-semibold">ERROR 404</p>
        <h1 className="mt-2 text-4xl md:text-5xl font-extrabold">
          Halaman tidak ditemukan
        </h1>
        <p className="mt-3 text-sub max-w-xl mx-auto">
          Maaf, video atau halaman yang kamu cari sudah dihapus, dipindahkan, atau memang tidak pernah ada.
          Coba cari lagi atau lihat video terbaru kami di bawah.
        </p>
        <div className="mt-6 flex gap-3 justify-center">
          <Link href="/" className="btn">Ke Beranda</Link>
          <Link href="/latest" className="btn-ghost">Video Terbaru</Link>
          <Link href="/trending" className="btn-ghost">Lagi Trending</Link>
        </div>
      </div>

      {posts?.data?.length ? (
        <section className="mt-12">
          <h2 className="text-lg font-semibold mb-4">Mungkin kamu suka</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {posts.data.map((p) => <PostCard key={p.id} post={p} />)}
          </div>
        </section>
      ) : null}
    </div>
  );
}
