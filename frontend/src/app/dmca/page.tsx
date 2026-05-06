export const metadata = { title: 'DMCA / Takedown' };

export default function DMCA() {
  const email = process.env.NEXT_PUBLIC_DMCA_EMAIL || 'dmca@example.com';
  return (
    <div className="max-w-3xl mx-auto px-4 py-8 prose prose-invert">
      <h1>DMCA / Content Takedown</h1>
      <p>
        We respect intellectual property rights. This site aggregates publicly available
        metadata from third-party sources with permission. We do not host video files.
      </p>
      <h2>How to file a takedown notice</h2>
      <p>Send an email to <a href={`mailto:${email}`}>{email}</a> including:</p>
      <ul>
        <li>Identification of the copyrighted work claimed to be infringed.</li>
        <li>The exact URL(s) on this site you want removed.</li>
        <li>Your contact information (name, address, phone, email).</li>
        <li>A statement that you have a good-faith belief that the use is not authorized.</li>
        <li>A statement under penalty of perjury that the information is accurate and you are
          the rights owner or authorized agent.</li>
        <li>Your physical or electronic signature.</li>
      </ul>
      <p>We will respond to valid notices within a reasonable time.</p>
      <h2>Content removal</h2>
      <p>You may also use the “Report” button on any post to request removal.</p>
    </div>
  );
}
