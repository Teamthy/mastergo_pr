export async function apiFetch(url: string, options: RequestInit = {}) {
    const token = localStorage.getItem("token");

    const res = await fetch(`http://localhost:8080${url}`, {
        ...options,
        headers: {
            "Content-Type": "application/json",
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
            ...(options.headers || {}),
        },
    });

    return res;
}