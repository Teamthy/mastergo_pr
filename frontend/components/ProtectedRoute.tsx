'use client';

import { useAuth } from "@/lib/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, ReactNode, useState } from "react";

interface ProtectedRouteProps {
    children: ReactNode;
    requiredRoles?: string[];
}

export default function ProtectedRoute({
    children,
    requiredRoles,
}: ProtectedRouteProps) {
    const { user, loading } = useAuth();
    const router = useRouter();
    const [isRedirecting, setIsRedirecting] = useState(false);

    useEffect(() => {
        if (loading) return;

        if (!user) {
            setIsRedirecting(true);
            router.replace("/");
            return;
        }

        // Check if user has completed onboarding
        if (user.onboarding_status !== "COMPLETED") {
            setIsRedirecting(true);
            router.replace("/auth/signup");
            return;
        }

        if (requiredRoles && requiredRoles.length > 0) {
            // Add role-based logic if needed in the future
            // For now, just ensure user is authenticated
        }
    }, [user, loading, router, requiredRoles]);

    // Show loading state while checking authentication
    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-zinc-900 dark:border-white"></div>
            </div>
        );
    }

    // Show nothing while redirecting
    if (isRedirecting || !user) {
        return null;
    }

    return <>{children}</>;
}