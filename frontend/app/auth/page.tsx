'use client';

import SignupFlow from "./signup/page";
import { useAuth } from "@/lib/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function AuthPage() {
    const { user, loading } = useAuth();
    const router = useRouter();

    useEffect(() => {
        if (!loading && user) {
            router.replace("/dashboard");
        }
    }, [user, loading, router]);

    if (loading) return null;

    return (
        <div className="min-h-screen flex items-center justify-center">
            <SignupFlow />
        </div>
    );
}