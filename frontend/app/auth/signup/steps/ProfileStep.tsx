'use client';

import { useState } from 'react';
import { ArrowRight, ArrowLeft } from 'lucide-react';

interface ProfileStepProps {
  onSubmit: (phone: string, address: string) => void;
  onBack: () => void;
  isLoading: boolean;
}

export default function ProfileStep({
  onSubmit,
  onBack,
  isLoading,
}: ProfileStepProps) {
  const [phone, setPhone] = useState('');
  const [address, setAddress] = useState('');
  const [errors, setErrors] = useState<{ phone?: string; address?: string }>({});

  const validate = () => {
    const newErrors: typeof errors = {};

    if (!phone.trim()) {
      newErrors.phone = 'Phone number is required';
    } else if (phone.length < 10 || phone.length > 20) {
      newErrors.phone = 'Phone number must be 10-20 characters';
    } else if (!/^[\d\s\-\+()]+$/.test(phone)) {
      newErrors.phone = 'Phone number is invalid';
    }

    if (!address.trim()) {
      newErrors.address = 'Address is required';
    } else if (address.length < 10) {
      newErrors.address = 'Address must be at least 10 characters';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validate()) {
      onSubmit(phone.trim(), address.trim());
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="bg-green-50 border border-green-200 rounded-lg p-4">
        <p className="text-sm text-green-800">
          Almost done! Let's complete your profile.
        </p>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Phone Number
        </label>
        <input
          type="tel"
          value={phone}
          onChange={(e) => {
            setPhone(e.target.value);
            if (errors.phone) setErrors({ ...errors, phone: '' });
          }}
          placeholder="+1 (555) 123-4567"
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
            errors.phone
              ? 'border-red-500 focus:ring-red-500'
              : 'border-gray-300 focus:ring-indigo-500'
          }`}
        />
        {errors.phone && (
          <p className="mt-1 text-sm text-red-600">{errors.phone}</p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Address
        </label>
        <textarea
          value={address}
          onChange={(e) => {
            setAddress(e.target.value);
            if (errors.address) setErrors({ ...errors, address: '' });
          }}
          placeholder="123 Main Street, City, State 12345"
          rows={3}
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
            errors.address
              ? 'border-red-500 focus:ring-red-500'
              : 'border-gray-300 focus:ring-indigo-500'
          }`}
        />
        {errors.address && (
          <p className="mt-1 text-sm text-red-600">{errors.address}</p>
        )}
      </div>

      <div className="text-sm text-gray-600 bg-gray-50 p-3 rounded-lg">
        <p>This information will be used for:</p>
        <ul className="list-disc list-inside mt-2 space-y-1">
          <li>Communication and support</li>
          <li>Verification purposes</li>
          <li>Service delivery</li>
        </ul>
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
          disabled={isLoading || !phone || !address}
          className="flex-1 bg-indigo-600 text-white font-semibold py-3 rounded-lg hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition flex items-center justify-center gap-2"
        >
          {isLoading ? 'Creating Account...' : 'Complete Profile'}
          <ArrowRight size={18} />
        </button>
      </div>
    </form>
  );
}
