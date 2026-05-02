'use client';

import { CheckCircle, Home } from 'lucide-react';
import { useEffect } from 'react';
import Link from 'next/link';

interface CompletedStepProps {
  onDone: () => void;
}

export default function CompletedStep({ onDone }: CompletedStepProps) {
  useEffect(() => {
    // Auto-redirect after 2 seconds
    const timer = setTimeout(() => {
      onDone();
    }, 2000);

    return () => clearTimeout(timer);
  }, [onDone]);
  return (
    <div className="text-center py-8">
      <div className="mx-auto w-20 h-20 bg-gradient-to-br from-green-50 to-emerald-50 rounded-full flex items-center justify-center mb-6 animate-pulse">
        <CheckCircle className="w-12 h-12 text-green-600" />
      </div>

      <h2 className="text-3xl font-bold text-gray-900 mb-2">Account Created!</h2>

      <p className="text-gray-600 mb-2">
        Your account has been successfully created and verified.
      </p>
      <p className="text-sm text-gray-500 mb-8">
        You're all set to start using the platform.
      </p>

      <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
        <p className="text-sm text-green-800">
          ✓ Email verified<br />
          ✓ Profile completed<br />
          ✓ Ready to use
        </p>
      </div>

      <button
        onClick={onDone}
        className="w-full bg-gradient-to-r from-indigo-600 to-indigo-700 text-white font-semibold py-3 rounded-lg hover:from-indigo-700 hover:to-indigo-800 transition flex items-center justify-center gap-2 mb-4"
      >
        <Home size={20} />
        Go to Wallet
      </button>


    </div>
  );
}
