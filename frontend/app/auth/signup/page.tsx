'use client';

import { useMachine } from "@xstate/react";
import { onboardingMachine } from "@/lib/onboardingMachine";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/AuthContext";
import { PasswordStrengthMeter } from "@/components/PasswordStrengthMeter";
import NameStep from "./steps/NameStep";
import CredentialsStep from "./steps/CredentialsStep";
import OTPStep from "./steps/OTPStep";
import ProfileStep from "./steps/ProfileStep";
import CompletedStep from "./steps/CompletedStep";

export default function SignupFlow() {
    const [state, send] = useMachine(onboardingMachine);
    const { signup, verifyEmail, updateProfile, login } = useAuth();
    const router = useRouter();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [otpCountdown, setOtpCountdown] = useState(0);

    useEffect(() => {
        if (otpCountdown <= 0) return;
        const timer = setTimeout(() => setOtpCountdown(otpCountdown - 1), 1000);
        return () => clearTimeout(timer);
    }, [otpCountdown]);

    const handleNameSubmit = (firstName: string, lastName: string) => {
        if (!firstName.trim() || !lastName.trim()) {
            setError("First name and last name are required");
            return;
        }
        if (!/^[a-zA-Z\s'-]+$/.test(firstName) || !/^[a-zA-Z\s'-]+$/.test(lastName)) {
            setError("Names must contain only letters");
            return;
        }
        setError(null);
        send({
            type: "SUBMIT_NAME",
            firstName: firstName.trim(),
            lastName: lastName.trim(),
        });
    };

    const handleCredentialsSubmit = async (
        email: string,
        password: string,
        confirmPassword: string
    ) => {
        setError(null);
        setLoading(true);

        try {
            await signup(state.context.firstName, state.context.lastName, email, password, confirmPassword);
            send({
                type: "SUBMIT_CREDENTIALS",
                email,
                password,
                confirmPassword,
            });
            setOtpCountdown(60);
        } catch (err: any) {
            const errorMsg = err.response?.data?.message || err.message || "Signup failed";
            setError(errorMsg);
        } finally {
            setLoading(false);
        }
    };

    const handleOTPSubmit = async (otp: string) => {
        setError(null);
        setLoading(true);

        try {
            await verifyEmail(state.context.email, otp);
            send({
                type: "VERIFY_OTP",
                otp,
            });
        } catch (err: any) {
            const errorMsg = err.response?.data?.message || err.message || "Verification failed";
            setError(errorMsg);
        } finally {
            setLoading(false);
        }
    };

    const handleResendOTP = async () => {
        setError(null);
        setLoading(true);

        try {
            const res = await fetch("http://127.0.0.1:8080/auth/resend-otp", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email: state.context.email }),
            });
            if (!res.ok) throw new Error("Failed to resend OTP");
            setOtpCountdown(60);
        } catch (err: any) {
            setError(err.message || "Failed to resend OTP");
        } finally {
            setLoading(false);
        }
    };

    const handleProfileSubmit = async (phone: string, address: string) => {
        setError(null);
        setLoading(true);

        try {
            // Login first to get JWT token
            await login(state.context.email, state.context.password);
            // Then update profile with authenticated session
            await updateProfile(phone, address);
            send({
                type: "SUBMIT_PROFILE",
                phone,
                address,
            });
            setTimeout(() => router.push("/dashboard"), 1000);
        } catch (err: any) {
            const errorMsg = err.response?.data?.message || err.message || "Profile update failed";
            setError(errorMsg);
        } finally {
            setLoading(false);
        }
    };

    const handleBack = () => {
        if (state.value === "credentials") {
            send({ type: "BACK" });
        } else if (state.value === "otp") {
            send({ type: "BACK" });
        } else if (state.value === "profile") {
            send({ type: "BACK" });
        }
    };

    return (
        <div className="space-y-6">
            {/* Error Alert */}
            {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                    <p className="text-sm text-red-800 font-medium">Error</p>
                    <p className="text-sm text-red-700">{error}</p>
                </div>
            )}

            {/* Name Step */}
            {state.value === "name" && (
                <NameStep onSubmit={handleNameSubmit} isLoading={loading} />
            )}

            {/* Credentials Step */}
            {state.value === "credentials" && (
                <CredentialsStep
                    onSubmit={handleCredentialsSubmit}
                    onBack={handleBack}
                    isLoading={loading}
                    initialEmail={state.context.email}
                />
            )}

            {/* OTP Step */}
            {state.value === "otp" && (
                <OTPStep
                    onSubmit={handleOTPSubmit}
                    onBack={handleBack}
                    onResendOTP={handleResendOTP}
                    isLoading={loading}
                    countdown={otpCountdown}
                    email={state.context.email}
                />
            )}

            {/* Profile Step */}
            {state.value === "profile" && (
                <ProfileStep
                    onSubmit={handleProfileSubmit}
                    onBack={handleBack}
                    isLoading={loading}
                />
            )}

            {/* Completed */}
            {state.value === "completed" && (
                <CompletedStep onDone={() => router.push("/dashboard")} />
            )}
        </div>
    );
}