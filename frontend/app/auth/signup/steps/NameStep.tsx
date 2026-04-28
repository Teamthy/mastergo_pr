'use client';

import { useState } from 'react';
import { ArrowRight } from 'lucide-react';

interface NameStepProps {
  onSubmit: (firstName: string, lastName: string) => void;
  isLoading: boolean;
}

export default function NameStep({ onSubmit, isLoading }: NameStepProps) {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [errors, setErrors] = useState<{ firstName?: string; lastName?: string }>({});

  const validate = () => {
    const newErrors: typeof errors = {};

    if (!firstName.trim()) {
      newErrors.firstName = 'First name is required';
    } else if (firstName.length < 2 || firstName.length > 50) {
      newErrors.firstName = 'First name must be 2-50 characters';
    } else if (!/^[a-zA-Z\s'-]+$/.test(firstName)) {
      newErrors.firstName = 'First name must contain only letters';
    }

    if (!lastName.trim()) {
      newErrors.lastName = 'Last name is required';
    } else if (lastName.length < 2 || lastName.length > 50) {
      newErrors.lastName = 'Last name must be 2-50 characters';
    } else if (!/^[a-zA-Z\s'-]+$/.test(lastName)) {
      newErrors.lastName = 'Last name must contain only letters';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validate()) {
      onSubmit(firstName.trim(), lastName.trim());
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          First Name
        </label>
        <input
          type="text"
          value={firstName}
          onChange={(e) => {
            setFirstName(e.target.value);
            if (errors.firstName) setErrors({ ...errors, firstName: '' });
          }}
          placeholder="John"
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
            errors.firstName
              ? 'border-red-500 focus:ring-red-500'
              : 'border-gray-300 focus:ring-indigo-500'
          }`}
        />
        {errors.firstName && (
          <p className="mt-1 text-sm text-red-600">{errors.firstName}</p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Last Name
        </label>
        <input
          type="text"
          value={lastName}
          onChange={(e) => {
            setLastName(e.target.value);
            if (errors.lastName) setErrors({ ...errors, lastName: '' });
          }}
          placeholder="Doe"
          className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
            errors.lastName
              ? 'border-red-500 focus:ring-red-500'
              : 'border-gray-300 focus:ring-indigo-500'
          }`}
        />
        {errors.lastName && (
          <p className="mt-1 text-sm text-red-600">{errors.lastName}</p>
        )}
      </div>

      <button
        type="submit"
        disabled={isLoading || !firstName.trim() || !lastName.trim()}
        className="w-full bg-indigo-600 text-white font-semibold py-3 rounded-lg hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition flex items-center justify-center gap-2"
      >
        {isLoading ? 'Loading...' : 'Continue'}
        <ArrowRight size={18} />
      </button>
    </form>
  );
}
