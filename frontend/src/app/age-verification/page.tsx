export const metadata = { title: '18+ Age Verification' };
export default function AgeVerification() {
  return (
    <div className="max-w-3xl mx-auto px-4 py-8 prose prose-invert">
      <h1>18+ Age Verification</h1>
      <p>This website contains material intended only for adults aged 18 or older
        (or the age of majority in your jurisdiction).</p>
      <p>By accessing this site you confirm:</p>
      <ul>
        <li>You are at least 18 years of age.</li>
        <li>You are accessing this material of your own free will.</li>
        <li>Adult material is legal in your country and locality.</li>
        <li>You will not allow minors to access this site.</li>
      </ul>
    </div>
  );
}
