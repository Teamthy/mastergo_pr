'use client';

import { useState } from 'react';
import { ArrowRight, ArrowLeft } from 'lucide-react';

interface OTPStepProps {
  onSubmit: (otp: string) => void;
  onBack: () => void;
  onResendOTP: () => void;
  isLoading: boolean;
  countdown: number;
  email: string;
}

export default function OTPStep({
  onSubmit,
  onBack,
  onResendOTP,
  isLoading,
  countdown,
  email,
}: OTPStepProps) {
  const [otp, setOtp] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (otp.length !== 6) {
      setError('OTP must be 6 digits');
      return;
    }
    if (!/^\d+$/.test(otp)) {
      setError('OTP must contain only numbers');
      return;
    }
    setError('');
    onSubmit(otp);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <p className="text-sm text-blue-800">
          We've sent a 6-digit OTP to <span className="font-semibold">{email}</span>
        </p>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Enter OTP
        </label>
        <input
          type="text"
          value={otp}
          onChange={(e) => {
            const value = e.target.value.replace(/\D/g, '').slice(0, 6);
            setOtp(value);
            if (error) setError('');
          }}
          placeholder="000000"
          maxLength={6}
          className={`w-full px-4 py-2 text-2xl text-center tracking-widest font-mono border rounded-lg focus:outline-none focus:ring-2 ${
            error
              ? 'border-red-500 focus:ring-red-500'
              : 'border-gray-300 focus:ring-indigo-500'
          }`}
        />
        {error && (
          <p className="mt-1 text-sm text-red-600">{error}</p>
        )}
      </div>

      <div className="text-center">
        <p className="text-sm text-gray-600 mb-2">Didn't receive the code?</p>
        <button
          type="button"
          onClick={onResendOTP}
          disabled={countdown > 0 || isLoading}
          className="text-indigo-600 hover:text-indigo-700 font-semibold text-sm disabled:text-gray-400 disabled:cursor-not-allowed"
        >
          {countdown > 0 ? `Resend in ${countdown}s` : 'Resend OTP'}
        </button>
      </div>

      <div className="flex gap-3 pt-4">
        <button
          type="button"
          onClick={onBack}
          disabled={isLoading}
          className="flex-1 border border-gray-300 text-gray-700 font-semibold py-3 rounded-lg hover:bg-gray-50 disabled:bg-gray-100 transition flex items-center justify-center gap-2"
        >
          <ArrowLeft size={18} />
          Back
        </button>
        <button
          type="submit"
          disabled={isLoading || otp.length !== 6}
          className="flex-1 bg-indigo-600 text-white font-semibold py-3 rounded-lg hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition flex items-center justify-center gap-2"
        >
          {isLoading ? 'Verifying...' : 'Verify'}
          <ArrowRight size={18} />
        </button>
      </div>
    </form>
  );
}
