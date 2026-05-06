'use client';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function SearchBar() {
  const r = useRouter();
  const [q, setQ] = useState('');
  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        if (q.trim()) r.push(`/search?q=${encodeURIComponent(q.trim())}`);
      }}
    >
      <input
        className="input"
        placeholder="Cari video bokep..."
        value={q}
        onChange={(e) => setQ(e.target.value)}
      />
    </form>
  );
}
