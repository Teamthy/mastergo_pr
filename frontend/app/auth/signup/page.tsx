'use client';

import { useMachine } from "@xstate/react";
import { onboardingMachine } from "@/lib/onboardingMachine";
import { motion } from "framer-motion";
import { useState, useEffect } from "react";
import { ArrowRight, CheckCircle2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/AuthContext";
import { authAPI } from "@/lib/api";

export default function SignupFlow() {
    const [state, send] = useMachine(onboardingMachine);
    const [input, setInput] = useState("");
    const [password, setPassword] = useState("");
    const router = useRouter();
    const { login } = useAuth();

    // handle submit
    const handleNext = async () => {
        if (state.matches("email")) {
            await authAPI.signup(input, password);

            send({
                type: "SUBMIT_EMAIL",
                email: input,
                password,
            });

            setInput("");
            setPassword("");
            return;
        }

        if (state.matches("otp")) {
            await authAPI.verifyOTP(state.context.email, input);

            await login(state.context.email, state.context.password);

            send({ type: "VERIFY_OTP", otp: input });

            router.replace("/dashboard");
        }
    };

    if (state.matches("completed")) {
        return (
            <div className="text-center">
                <CheckCircle2 className="mx-auto" size={48} />
                <p>Account verified</p>
            </div>
        );
    }

    return (
        <div className="max-w-md mx-auto space-y-4">

            {state.matches("email") && (
                <>
                    <input
                        placeholder="Email"
                        className="border p-3 w-full"
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                    />

                    <input
                        type="password"
                        placeholder="Password"
                        className="border p-3 w-full"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </>
            )}

            {state.matches("otp") && (
                <input
                    placeholder="Enter OTP"
                    className="border p-3 w-full"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                />
            )}

            <motion.button
                whileTap={{ scale: 0.95 }}
                className="bg-black text-white px-4 py-3 w-full"
                onClick={handleNext}
            >
                Continue <ArrowRight className="inline ml-2" size={16} />
            </motion.button>
        </div>
    );
}