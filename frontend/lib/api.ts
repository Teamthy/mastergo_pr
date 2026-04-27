const API = "http://localhost:8080";

export const authAPI = {
    signup: async (email: string, password: string) => {
        const res = await fetch(`${API}/auth/signup`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
        });

        if (!res.ok) throw new Error(await res.text());
        return res.json();
    },

    verifyOTP: async (email: string, otp: string) => {
        const res = await fetch(`${API}/auth/verify-email`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, otp }),
        });

        if (!res.ok) throw new Error(await res.text());
        return res.json();
    },

    login: async (email: string, password: string) => {
        const res = await fetch(`${API}/auth/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
        });

        if (!res.ok) throw new Error(await res.text());
        return res.json(); // { token, user }
    },

    me: async (token: string) => {
        const res = await fetch(`${API}/auth/me`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });

        if (!res.ok) throw new Error("unauthorized");
        return res.json();
    },
};