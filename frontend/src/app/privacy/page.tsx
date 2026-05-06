export const metadata = { title: 'Privacy Policy' };
export default function Privacy() {
  return (
    <div className="max-w-3xl mx-auto px-4 py-8 prose prose-invert">
      <h1>Privacy Policy</h1>
      <p>We store minimal data: server access logs and basic analytics. Age-gate
        confirmation is stored in your browser localStorage only.</p>
      <p>We do not knowingly collect personal information from anyone under 18.</p>
      <p>Third-party embedded video players may set their own cookies.</p>
    </div>
  );
}
