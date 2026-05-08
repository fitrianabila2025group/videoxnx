'use client';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function SearchBar() {
  const r = useRouter();
  const [q, setQ] = useState('');
  return (
    <form
      role="search"
      onSubmit={(e) => {
        e.preventDefault();
        if (q.trim()) r.push(`/search?q=${encodeURIComponent(q.trim())}`);
      }}
      className="relative"
    >
      <span className="pointer-events-none absolute inset-y-0 left-3 flex items-center text-sub">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
          <circle cx="11" cy="11" r="7" />
          <path d="M21 21l-4.3-4.3" />
        </svg>
      </span>
      <input
        type="search"
        inputMode="search"
        aria-label="Cari video"
        className="input pl-11 pr-3"
        placeholder="Cari video bokep..."
        value={q}
        onChange={(e) => setQ(e.target.value)}
      />
    </form>
  );
}
