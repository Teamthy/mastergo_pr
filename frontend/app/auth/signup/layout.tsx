'use client';

import { useMachine } from '@xstate/react';
import { onboardingMachine } from '@/lib/onboardingMachine';
import { useState } from 'react';

export default function SignupLayout({ children }: { children: React.ReactNode }) {
  const [state, send] = useMachine(onboardingMachine);
  const [otpTimer, setOtpTimer] = useState(0);

  const currentStep = state.value;
  const steps = ['name', 'credentials', 'otp', 'profile'];
  const stepIndex = steps.indexOf(currentStep as string);
  const stepNumber = stepIndex >= 0 ? stepIndex + 1 : 0;

  const progress = {
    name: 25,
    credentials: 50,
    otp: 75,
    profile: 100,
    completed: 100,
  }[currentStep as string] || 0;

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-8">


        {/* Signup Form Card */}
        <div className="max-w-md mx-auto bg-white rounded-lg shadow-lg p-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Create Account</h1>
          <p className="text-gray-600 mb-6">

            {currentStep === 'credentials' && 'Set up your password'}
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
