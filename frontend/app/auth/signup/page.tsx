'use client';

import { useMachine } from "@xstate/react";
import { onboardingMachine } from "@/lib/onboardingMachine";
import { OnboardingState } from "@/lib/types";
import { motion } from "framer-motion";
import { useState, useEffect } from "react";
import { ArrowRight, CheckCircle2, ShieldCheck } from "lucide-react";
import { useRouter } from "next/navigation";
import { cn } from "@/lib/utils";
import { useAuth } from "@/lib/AuthContext";

export default function SignupFlow() {
    const [state, send] = useMachine(onboardingMachine);
    const [inputValue, setInputValue] = useState("");
    const [loading, setLoading] = useState(false);
    const router = useRouter();
    const { login, user } = useAuth();

    // auto redirect after login
    useEffect(() => {
        if (user) {
            router.push("/dashboard");
        }
    }, [user, router]);

    const handleNext = async () => {
        setLoading(true);
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

    // when onboarding completes → login
    useEffect(() => {
        if (state.matches(OnboardingState.COMPLETED)) {
            login(state.context.email);

        }
    }, [state, login, router]);

    if (state.matches(OnboardingState.COMPLETED)) {
        return (
            <div className="max-w-md w-full bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-12 text-center flex flex-col gap-8 items-center shadow-2xl">
                <div className="w-20 h-20 rounded-full bg-emerald-500/10 flex items-center justify-center text-emerald-500">
                    <CheckCircle2 size={48} />
                </div>

                <div>
                    <h1 className="text-3xl font-bold uppercase tracking-tight">
                        Access Granted
                    </h1>
                    <p className="text-zinc-500 mt-2">
                        Your  identity has been verified.
                    </p>
                </div>

                <div className="w-full h-12 rounded-xl bg-black text-white flex items-center justify-center">
                    Redirecting...
                </div>
            </div>
        );
    }

    const getStepContent = () => {
        const currentState = state.value as string;

        switch (currentState) {
            case OnboardingState.START:
                return {
                    title: "EMAIL VERIFICATION",
                    description: "Enter your email to begin verification.",
                    placeholder: "name@email.com",
                };

            case OnboardingState.EMAIL_ENTERED:
                return {
                    title: "OTP Verification",
                    description: "A secure code was sent to your email.",
                    placeholder: "000-000",
                };

            case OnboardingState.EMAIL_VERIFIED:
                return {
                    title: "NAME DETAILS",
                    description: "Enter your full name.",
                    placeholder: "Teamthy Teethy",
                };

            case OnboardingState.PROFILE_NAME:
                return {
                    title: "CONTACT NUMBER",
                    description: "Enter your phone number.",
                    placeholder: "+2341234567",
                };

            case OnboardingState.PROFILE_CONTACT:
                return {
                    title: "ADDRESS DETAILS",
                    description: "Enter your address.",
                    placeholder: "Street, City, Country",
                };

            default:
                return { title: "", description: "", placeholder: "" };
        }
    };

    const content = getStepContent();

    return (
        <div className="max-w-xl w-full flex flex-col gap-12 p-6">
            <div className="text-center">
                <ShieldCheck className="mx-auto mb-4" size={28} />
                <h1 className="text-3xl font-bold">{content.title}</h1>
                <p className="text-zinc-500">{content.description}</p>
            </div>

            <motion.div
                key={state.value as string}
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                className="flex flex-col gap-4"
            >
                <input
                    className="border p-3 rounded-xl"
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    placeholder={content.placeholder}
                    disabled={loading}
                />

                <button
                    className="h-12 rounded-xl bg-black text-white flex items-center justify-center"
                    onClick={handleNext}
                >
                    Continue <ArrowRight className="ml-2" size={18} />
                </button>
            </motion.div>

            <div className="flex justify-center gap-2">
                {[0, 1, 2, 3, 4].map((i) => (
                    <div
                        key={i}
                        className={cn(
                            "w-2 h-2 rounded-full",
                            Object.values(OnboardingState).indexOf(state.value as any) >= i
                                ? "bg-black"
                                : "bg-gray-300"
                        )}
                    />
                ))}
            </div>
        </div>
    );
}