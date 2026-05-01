'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/lib/AuthContext';
import { Check } from 'lucide-react';

interface PasswordStrengthMeterProps {
  password: string;
  onStrengthChange?: (strength: 'weak' | 'medium' | 'strong') => void;
}

export function PasswordStrengthMeter({
  password,
  onStrengthChange,
}: PasswordStrengthMeterProps) {
  const { getPasswordStrength } = useAuth();
  const [strength, setStrength] = useState<'weak' | 'medium' | 'strong' | null>(null);
  const [requirements, setRequirements] = useState({
    minLength: false,
    uppercase: false,
    lowercase: false,
    number: false,
    special: false,
  });

  // Debounce password strength check
  useEffect(() => {
    if (!password) {
      setStrength(null);
      return;
    }

    const timer = setTimeout(async () => {
      const result = await getPasswordStrength(password);
      setStrength(result);
      onStrengthChange?.(result);
    }, 300);

    return () => clearTimeout(timer);
  }, [password, getPasswordStrength, onStrengthChange]);

  // Check requirements in real-time
  useEffect(() => {
    setRequirements({
      minLength: password.length >= 8,
      uppercase: /[A-Z]/.test(password),
      lowercase: /[a-z]/.test(password),
      number: /[0-9]/.test(password),
      special: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password),
    });
  }, [password]);

  if (!password) return null;

  const strengthColor = {
    weak: 'bg-red-500',
    medium: 'bg-yellow-500',
    strong: 'bg-green-500',
  };

  const strengthLabel = {
    weak: 'Weak',
    medium: 'Medium',
    strong: 'Strong',
  };

  const strengthPercentage = {
    weak: 33,
    medium: 66,
    strong: 100,
  };

  return (
    <div className="space-y-3 mt-3">
      {/* Strength Meter */}
      <div className="space-y-1">
        <div className="flex items-center justify-between">
          <label className="block text-xs font-medium text-gray-600">
            Password Strength
          </label>
          {strength && (
            <span className="text-xs font-semibold">
              <span className={`${
                strength === 'weak'
                  ? 'text-red-600'
                  : strength === 'medium'
                  ? 'text-yellow-600'
                  : 'text-green-600'
              }`}>
                {strengthLabel[strength]}
              </span>
            </span>
          )}
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
          <div
            className={`h-full transition-all duration-300 ${
              strength ? strengthColor[strength] : 'bg-gray-300'
            }`}
            style={{
              width: strength ? `${strengthPercentage[strength]}%` : '0%',
            }}
          />
        </div>
      </div>

      {/* Requirements Checklist */}
      <div className="space-y-2">
        <p className="text-xs font-medium text-gray-600">Requirements:</p>
        <ul className="space-y-1">
          <li className="flex items-center gap-2">
            <span className={`w-4 h-4 rounded flex items-center justify-center ${
              requirements.minLength
                ? 'bg-green-100 text-green-600'
                : 'bg-gray-100 text-gray-400'
            }`}>
              {requirements.minLength && <Check size={14} />}
            </span>
            <span className={`text-xs ${
              requirements.minLength
                ? 'text-green-700'
                : 'text-gray-600'
            }`}>
              At least 8 characters
            </span>
          </li>
          <li className="flex items-center gap-2">
            <span className={`w-4 h-4 rounded flex items-center justify-center ${
              requirements.uppercase
                ? 'bg-green-100 text-green-600'
                : 'bg-gray-100 text-gray-400'
            }`}>
              {requirements.uppercase && <Check size={14} />}
            </span>
            <span className={`text-xs ${
              requirements.uppercase
                ? 'text-green-700'
                : 'text-gray-600'
            }`}>
              One uppercase letter
            </span>
          </li>
          <li className="flex items-center gap-2">
            <span className={`w-4 h-4 rounded flex items-center justify-center ${
              requirements.lowercase
                ? 'bg-green-100 text-green-600'
                : 'bg-gray-100 text-gray-400'
            }`}>
              {requirements.lowercase && <Check size={14} />}
            </span>
            <span className={`text-xs ${
              requirements.lowercase
                ? 'text-green-700'
                : 'text-gray-600'
            }`}>
              One lowercase letter
            </span>
          </li>
          <li className="flex items-center gap-2">
            <span className={`w-4 h-4 rounded flex items-center justify-center ${
              requirements.number
                ? 'bg-green-100 text-green-600'
                : 'bg-gray-100 text-gray-400'
            }`}>
              {requirements.number && <Check size={14} />}
            </span>
            <span className={`text-xs ${
              requirements.number
                ? 'text-green-700'
                : 'text-gray-600'
            }`}>
              One number
            </span>
          </li>
          <li className="flex items-center gap-2">
            <span className={`w-4 h-4 rounded flex items-center justify-center ${
              requirements.special
                ? 'bg-green-100 text-green-600'
                : 'bg-gray-100 text-gray-400'
            }`}>
              {requirements.special && <Check size={14} />}
            </span>
            <span className={`text-xs ${
              requirements.special
                ? 'text-green-700'
                : 'text-gray-600'
            }`}>
              One special character (!@#$%^&*)
            </span>
          </li>
        </ul>
      </div>
    </div>
  );
}
