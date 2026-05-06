import Link from 'next/link';

export default function Pagination({
  page,
  perPage,
  total,
  basePath,
  query = {},
}: {
  page: number;
  perPage: number;
  total: number;
  basePath: string;
  query?: Record<string, string>;
}) {
  const totalPages = Math.max(1, Math.ceil(total / perPage));
  const make = (p: number) => {
    const params = new URLSearchParams({ ...query, page: String(p) });
    return `${basePath}?${params.toString()}`;
  };
  if (totalPages <= 1) return null;
  return (
    <nav className="flex flex-wrap items-center justify-center gap-2 mt-8">
      {page > 1 && <Link href={make(page - 1)} className="btn-ghost">← Prev</Link>}
      <span className="text-sub text-sm">Page {page} of {totalPages}</span>
      {page < totalPages && <Link href={make(page + 1)} className="btn-ghost">Next →</Link>}
    </nav>
  );
}
