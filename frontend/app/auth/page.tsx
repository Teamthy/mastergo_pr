'use client';

import SignupFlow from "./signup/page";
import { useAuth } from "@/lib/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import Link from "next/link";

export default function AuthPage() {
    const { user, loading } = useAuth();
    const router = useRouter();
    const [showChoice, setShowChoice] = useState(false);

    useEffect(() => {
        if (!loading && user) {
            router.replace("/dashboard/wallet");
        } else if (!loading && !user) {
            setShowChoice(true);
        }
    }, [user, loading, router]);

    if (loading) return null;

    if (showChoice && !user) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="w-full max-w-md mx-auto text-center space-y-6">
                    <div>
                        <h1 className="text-4xl font-bold tracking-tighter mb-2">Welcome</h1>

                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Link href="/auth/login">
                            <button
                                className="w-full bg-indigo-600 text-white font-semibold py-4 rounded-lg hover:bg-indigo-700 transition"
                            >
                                Sign In
                            </button>
                        </Link>
                        <Link href="/auth/signup">
                            <button
                                className="w-full bg-white border-2 border-indigo-600 text-indigo-600 font-semibold py-4 rounded-lg hover:bg-indigo-50 transition"
                            >
                                Create Account
                            </button>
                        </Link>
                    </div>

                    <div className="relative">
                        <div className="absolute inset-0 flex items-center">
                            <div className="w-full border-t border-gray-300"></div>
                        </div>

                    </div>

                    <Link href="/">
                        <button
                            className="w-full text-gray-600 font-semibold py-3 rounded-lg hover:bg-gray-100 transition"
                        >
                            back to landing page
                        </button>
                    </Link>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen flex items-center justify-center">
            <SignupFlow />
        </div>
    );
}