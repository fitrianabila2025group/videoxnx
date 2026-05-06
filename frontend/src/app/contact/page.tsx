export const metadata = { title: 'Contact' };
export default function Contact() {
  const email = process.env.NEXT_PUBLIC_DMCA_EMAIL || 'contact@example.com';
  return (
    <div className="max-w-3xl mx-auto px-4 py-8 prose prose-invert">
      <h1>Contact</h1>
      <p>For takedown requests, partnerships, or general questions, email:{' '}
        <a href={`mailto:${email}`}>{email}</a>.</p>
    </div>
  );
}
