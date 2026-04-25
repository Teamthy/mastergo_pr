'use client';

import { createContext, useContext, useState, ReactNode, useEffect } from "react";
import { useRouter } from "next/navigation";

type User = {
    email: string;
};

type AuthContextType = {
    user: User | null;
    login: (email: string) => void;
    logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const router = useRouter();

    // restore session
    useEffect(() => {
        const stored = localStorage.getItem("vanguard_user");
        if (stored) {
            setUser(JSON.parse(stored));
        }
    }, []);

    const login = (email: string) => {
        const newUser = { email };
        setUser(newUser);
        localStorage.setItem("vanguard_user", JSON.stringify(newUser));
        router.push("/dashboard");
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem("vanguard_user");
        router.push("/");
    };

    return (
        <AuthContext.Provider value={{ user, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const ctx = useContext(AuthContext);
    if (!ctx) {
        throw new Error("useAuth must be used inside AuthProvider");
    }
    return ctx;
}