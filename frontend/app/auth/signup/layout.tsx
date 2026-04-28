'use client';

import { useMachine } from '@xstate/react';
import { onboardingMachine } from '@/lib/onboardingMachine';
import { useState } from 'react';

export default function SignupLayout({ children }: { children: React.ReactNode }) {
  const [state, send] = useMachine(onboardingMachine);
  const [otpTimer, setOtpTimer] = useState(0);

  const currentStep = state.value;
  const progress = {
    name: 25,
    credentials: 50,
    otp: 75,
    profile: 100,
  }[currentStep as string] || 0;

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-8">
        {/* Progress Bar */}
        <div className="max-w-md mx-auto mb-8">
          <div className="flex justify-between mb-2">
            <span className="text-sm font-medium text-gray-700">
              {currentStep === 'completed' ? 'Complete!' : `Step ${['name', 'credentials', 'otp', 'profile'].indexOf(currentStep as string) + 1} of 4`}
            </span>
            <span className="text-sm font-medium text-gray-600">{progress}%</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className="bg-indigo-600 h-2 rounded-full transition-all duration-300"
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>

        {/* Signup Form Card */}
        <div className="max-w-md mx-auto bg-white rounded-lg shadow-lg p-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Create Account</h1>
          <p className="text-gray-600 mb-6">
            {currentStep === 'name' && 'Tell us your name'}
            {currentStep === 'credentials' && 'Create your login credentials'}
            {currentStep === 'otp' && 'Verify your email'}
            {currentStep === 'profile' && 'Complete your profile'}
            {currentStep === 'completed' && 'Account created successfully!'}
          </p>

          {children}
        </div>

        {/* Footer */}
        <div className="max-w-md mx-auto mt-6 text-center text-sm text-gray-600">
          <p>Already have an account? <a href="/auth" className="text-indigo-600 hover:underline">Sign in</a></p>
        </div>
      </div>
    </div>
  );
}
