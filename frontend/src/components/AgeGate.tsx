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
    <div className="fixed inset-0 z-[60] bg-black/95 backdrop-blur-sm flex items-center justify-center p-4 animate-fade">
      <div className="max-w-md w-full bg-panel border border-white/10 rounded-2xl p-5 sm:p-6 text-center shadow-2xl">
        <div className="mx-auto w-14 h-14 rounded-full bg-brand/20 grid place-items-center mb-3">
          <span className="text-brand font-extrabold text-lg">18+</span>
        </div>
        <h1 className="text-xl sm:text-2xl font-bold text-brand">Khusus Dewasa 18+</h1>
        <p className="mt-3 text-sm sm:text-base text-sub leading-relaxed">
          Situs ini berisi konten dewasa. Dengan masuk, kamu menyatakan sudah
          berusia minimal 18 tahun (atau usia dewasa di wilayahmu) dan setuju
          dengan{' '}
          <Link href="/disclaimer" className="underline text-ink">disclaimer</Link>{' '}
          kami.
        </p>
        <div className="mt-5 sm:mt-6 grid grid-cols-2 gap-3">
          <button onClick={accept} className="btn">Saya 18+</button>
          <button onClick={decline} className="btn-ghost">Keluar</button>
        </div>
        <p className="mt-4 text-xs text-sub">
          Konten hanya untuk hiburan dewasa.
        </p>
      </div>
    </div>
  );
}
