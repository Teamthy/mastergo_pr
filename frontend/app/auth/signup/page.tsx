'use client';

import { useMachine } from "@xstate/react";
import { onboardingMachine } from "../@/lib/onboardingMachine";
import { OnboardingState } from "../@/lib/types";
import { motion, AnimatePresence } from "motion/react";
import { Button, Input } from "../@/lib/components/ui";
import { useState } from "react";
import { ArrowRight, CheckCircle2, ShieldCheck } from "lucide-react";
import { useRouter } from "next/navigation";
import { cn } from "@/lib/utils";
import { useAuth } from "@/lib/AuthContext";

export default function SignupFlow() {
    const [state, send] = useMachine(onboardingMachine);
    const [inputValue, setInputValue] = useState("");
    const [loading, setLoading] = useState(false);
    const router = useRouter();
    const { login } = useAuth();

    const handleNext = async () => {
        setLoading(true);
        // Simulate API calls
        await new Promise((r) => setTimeout(r, 800));
        setLoading(false);

        const currentState = state.value as string;

        switch (currentState) {
            case OnboardingState.START:
                send({ type: "SUBMIT_EMAIL", email: inputValue });
                break;
            case OnboardingState.EMAIL_ENTERED:
                send({ type: "VERIFY_OTP", otp: inputValue });
                break;
            case OnboardingState.EMAIL_VERIFIED:
                send({ type: "SUBMIT_NAME", name: inputValue });
                break;
            case OnboardingState.PROFILE_NAME:
                send({ type: "SUBMIT_CONTACT", contact: inputValue });
                break;
            case OnboardingState.PROFILE_CONTACT:
                send({ type: "SUBMIT_ADDRESS", address: inputValue });
                break;
        }
        setInputValue("");
    };

    if (state.matches(OnboardingState.COMPLETED)) {
        return (
            <div className="max-w-md w-full bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-12 text-center flex flex-col gap-8 items-center shadow-2xl">
                <div className="w-20 h-20 rounded-full bg-emerald-500/10 flex items-center justify-center text-emerald-500">
                    <CheckCircle2 size={48} />
                </div>
                <div>
                    <h1 className="text-3xl font-bold uppercase tracking-tight text-black dark:text-white">Access Granted</h1>
                    <p className="text-zinc-500 mt-2">Your Vanguard identity has been verified.</p>
                </div>
                <Button size="lg" className="w-full" onClick={async () => {
                    await login(state.context.email || "user@example.com");
                    router.push('/dashboard');
                }}>
                    Enter Dashboard
                </Button>
            </div>
        );
    }

    const getStepContent = () => {
        const currentState = state.value as string;
        switch (currentState) {
            case OnboardingState.START:
                return {
                    title: "Verify Identity",
                    description: "Enter your official email to begin verification.",
                    placeholder: "name@company.com",
                    label: "Email Address",
                };
            case OnboardingState.EMAIL_ENTERED:
                return {
                    title: "OTP Verification",
                    description: "A secure code was sent to your email. Expires in 5:00.",
                    placeholder: "000-000",
                    label: "One-Time Password",
                };
            case OnboardingState.EMAIL_VERIFIED:
                return {
                    title: "Legal Identity",
                    description: "Enter your full name as it appears on official documents.",
                    placeholder: "Johnathan Doe",
                    label: "Full Name",
                };
            case OnboardingState.PROFILE_NAME:
                return {
                    title: "Secure Contact",
                    description: "Used for emergency identity recovery.",
                    placeholder: "+1 (555) 000-0000",
                    label: "Phone Number",
                };
            case OnboardingState.PROFILE_CONTACT:
                return {
                    title: "Registry Address",
                    description: "Your primary operating jurisdiction.",
                    placeholder: "Street, City, Country",
                    label: "Residential Address",
                };
            default:
                return { title: "", description: "", placeholder: "", label: "" };
        }
    };

    const content = getStepContent();

    return (
        <div className="max-w-xl w-full flex flex-col gap-12 p-6">
            {/* Header */}
            <div className="flex flex-col gap-4 text-center">
                <div className="flex justify-center">
                    <div className="w-12 h-12 bg-black dark:bg-white rounded-xl flex items-center justify-center">
                        <ShieldCheck className="text-white dark:text-black" size={28} />
                    </div>
                </div>
                <div className="flex flex-col gap-1">
                    <h1 className="text-4xl font-bold tracking-tighter uppercase text-black dark:text-white">{content.title}</h1>
                    <p className="text-zinc-500 font-medium">{content.description}</p>
                </div>
            </div>

            {/* Form */}
            <motion.div
                key={state.value as string}
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="flex flex-col gap-6"
            >
                <div className="flex flex-col gap-3">
                    <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4">
                        {content.label}
                    </label>
                    <Input
                        value={inputValue}
                        onChange={(e) => setInputValue(e.target.value)}
                        placeholder={content.placeholder}
                        disabled={loading}
                        onKeyDown={(e) => e.key === 'Enter' && handleNext()}
                    />
                </div>
                <Button
                    size="lg"
                    className="w-full h-14 text-base font-bold"
                    onClick={handleNext}
                    disabled={!inputValue || loading}
                >
                    {loading ? "Verifying..." : "Continue"} <ArrowRight className="ml-2" size={18} />
                </Button>
            </motion.div>

            {/* Progress Dots */}
            <div className="flex justify-center gap-2">
                {[0, 1, 2, 3, 4].map((i) => (
                    <div
                        key={i}
                        className={cn(
                            "w-1.5 h-1.5 rounded-full transition-all",
                            Object.values(OnboardingState).indexOf(state.value as any) >= i
                                ? "bg-black dark:bg-white scale-125"
                                : "bg-zinc-200 dark:bg-zinc-800"
                        )}
                    />
                ))}
            </div>
        </div>
    );
}
