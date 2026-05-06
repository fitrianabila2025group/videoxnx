'use client';
import { useState } from 'react';

export default function ReportButton({ postId }: { postId: number }) {
  const [open, setOpen] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [reason, setReason] = useState('dmca');
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [busy, setBusy] = useState(false);

  const submit = async () => {
    setBusy(true);
    try {
      await fetch('/api/reports', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ post_id: postId, reason, email, message }),
      });
      setSubmitted(true);
    } finally { setBusy(false); }
  };

  return (
    <>
      <button className="btn-ghost" onClick={() => setOpen(true)}>Laporkan</button>
      {open && (
        <div className="fixed inset-0 z-40 bg-black/70 flex items-center justify-center p-4" onClick={() => setOpen(false)}>
          <div className="bg-panel rounded-xl p-6 max-w-md w-full" onClick={(e) => e.stopPropagation()}>
            {submitted ? (
              <>
                <h3 className="font-semibold">Thank you</h3>
                <p className="text-sub mt-2">Your report was received and will be reviewed.</p>
                <button className="btn mt-4" onClick={() => setOpen(false)}>Close</button>
              </>
            ) : (
              <>
                <h3 className="font-semibold mb-3">Report content</h3>
                <label className="block text-sm text-sub">Reason</label>
                <select className="input mt-1" value={reason} onChange={(e) => setReason(e.target.value)}>
                  <option value="dmca">DMCA / copyright</option>
                  <option value="illegal">Illegal content</option>
                  <option value="non_consensual">Non-consensual</option>
                  <option value="underage">Suspected underage</option>
                  <option value="other">Other</option>
                </select>
                <label className="block text-sm text-sub mt-3">Your email (optional)</label>
                <input className="input mt-1" value={email} onChange={(e) => setEmail(e.target.value)} />
                <label className="block text-sm text-sub mt-3">Message</label>
                <textarea className="input mt-1" rows={4} value={message} onChange={(e) => setMessage(e.target.value)} />
                <div className="flex gap-2 mt-4 justify-end">
                  <button className="btn-ghost" onClick={() => setOpen(false)}>Cancel</button>
                  <button className="btn" disabled={busy} onClick={submit}>{busy ? 'Sending…' : 'Submit'}</button>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </>
  );
}
