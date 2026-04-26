export const apiFetch = (url: string, options: RequestInit = {}) => {
    const token = typeof window !== "undefined"
        ? localStorage.getItem("token")
        : null;

    return fetch(`http://localhost:8080${url}`, {
        ...options,
        headers: {
            "Content-Type": "application/json",
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
            ...(options.headers || {}),
        },
    });
};