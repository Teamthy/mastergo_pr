'use client';

import { useAuth } from "@/lib/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, ReactNode } from "react";

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

    useEffect(() => {
        if (loading) return;

        if (!user) {
            router.replace("/auth");
            return;
        }

        if (requiredRoles && requiredRoles.length > 0) {
            // Add role-based logic if needed in the future
            // For now, just ensure user is authenticated
        }
    }, [user, loading, router, requiredRoles]);

    // Show nothing while loading to prevent flash of unprotected content
    if (loading) return null;

    // Redirect will happen in useEffect, but show nothing until redirected
    if (!user) return null;

    return <>{children}</>;
}