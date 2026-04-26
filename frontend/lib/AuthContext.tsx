'use client';

import { createContext, useContext, useEffect, useState } from "react";

type AuthContextType = {
    user: any;
    login: (email: string) => Promise<void>;
    logout: () => void;
    loading: boolean;
};

const AuthContext = createContext<AuthContextType>({
    user: null,
    login: async () => { },
    logout: () => { },
    loading: true,
});

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const [user, setUser] = useState<any>(null);
    const [loading, setLoading] = useState(true);


    useEffect(() => {
        const token = localStorage.getItem("token");

        if (token) {
            setUser({ id: "persisted", email: "persisted" });
        }

        setLoading(false);
    }, []);

    useEffect(() => {
        const token = localStorage.getItem("token");

        if (!token) {
            setLoading(false);
            return;
        }

        fetch("http://localhost:8080/auth/me", {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        })
            .then(res => {
                if (!res.ok) throw new Error();
                return res.json();
            })
            .then(data => {
                setUser(data);
            })
            .catch(() => {
                localStorage.removeItem("token");
                setUser(null);
            })
            .finally(() => setLoading(false));
    }, []);

    const login = (email: string, token?: string) => {
        const fakeUser = {
            email,
            id: email,
        };

        localStorage.setItem("token", token ?? "dev-token");
        setUser(fakeUser);
    };

    const logout = () => {
        localStorage.removeItem("token");
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, login, logout, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);