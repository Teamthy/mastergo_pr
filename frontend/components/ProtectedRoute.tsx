'use client'

import { useAuth } from "@/lib/AuthContext"
import { useRouter } from "next/navigation"
import { useEffect } from "react"

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
    const { user } = useAuth()
    const router = useRouter()

    useEffect(() => {
        if (!user) router.push("/auth/signup")
    }, [user])

    if (!user) return null

    return <>{children}</>
}