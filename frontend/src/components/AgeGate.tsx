'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';

const KEY = 'age_verified_v1';

export default function AgeGate() {
  const [shown, setShown] = useState(false);

  useEffect(() => {
    try {
      if (!localStorage.getItem(KEY)) setShown(true);
    } catch {}
  }, []);

  if (!shown) return null;

  const accept = () => {
    try { localStorage.setItem(KEY, '1'); } catch {}
    setShown(false);
  };
  const decline = () => {
    window.location.href = 'https://www.google.com';
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/95 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-panel border border-white/10 rounded-2xl p-6 text-center">
        <h1 className="text-2xl font-bold text-brand">Khusus Dewasa 18+</h1>
        <p className="mt-3 text-sub">
          Situs ini berisi konten dewasa. Dengan masuk, kamu menyatakan sudah
          berusia minimal 18 tahun (atau usia dewasa di wilayahmu) dan setuju
          dengan{' '}
          <Link href="/disclaimer" className="underline text-ink">disclaimer</Link>{' '}
          kami.
        </p>
        <div className="mt-6 flex gap-3 justify-center">
          <button onClick={accept} className="btn">Saya 18+ — Masuk</button>
          <button onClick={decline} className="btn-ghost">Keluar</button>
        </div>
        <p className="mt-4 text-xs text-sub">
          Konten hanya untuk hiburan dewasa.
        </p>
      </div>
    </div>
  );
}
