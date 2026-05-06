'use client';

import { useMachine } from '@xstate/react';
import { onboardingMachine } from '@/lib/onboardingMachine';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/AuthContext';

export default function SignupLayout({ children }: { children: React.ReactNode }) {
  const [state, send] = useMachine(onboardingMachine);
  const [otpTimer, setOtpTimer] = useState(0);
  const router = useRouter();
  const { user, loading } = useAuth();

  // Redirect if already completed onboarding
  useEffect(() => {
    if (!loading && user && user.onboarding_status === "COMPLETED") {
      router.replace("/dashboard/wallet");
    }
  }, [user, loading, router]);

  if (loading) return null;

  return (
    <>
      {children}
    </>
  );
}
