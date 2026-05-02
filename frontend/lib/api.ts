const API = "http://localhost:8080";

// Helper function to handle API responses
const handleResponse = async (res: Response) => {
  const contentType = res.headers.get("content-type");
  let data: any = null;

  if (contentType && contentType.includes("application/json")) {
    data = await res.json();
  } else {
    data = await res.text();
  }

  if (!res.ok) {
    throw {
      status: res.status,
      response: { data },
    };
  }
  return data;
};

// Helper function to make requests
const apiRequest = async (
  method: string,
  url: string,
  body?: any,
  token?: string
) => {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  console.log("API CALL →", url);

  const res = await fetch(url, {
    method,
    headers,
    body: method !== "GET" && body ? JSON.stringify(body) : undefined,
  });
  return handleResponse(res);
};

export const authAPI = {
  // Signup with name, email, and password
  signup: async (data: {
    first_name: string;
    last_name: string;
    email: string;
    password: string;
    confirm_password: string;
  }) => {
    return apiRequest("POST", `${API}/auth/signup`, data);
  },

  // Login with email and password
  login: async (email: string, password: string) => {
    return apiRequest("POST", `${API}/auth/login`, { email, password });
  },

  // Verify email with OTP
  verifyEmail: async (email: string, otp: string) => {
    return apiRequest("POST", `${API}/auth/verify-email`, { email, otp });
  },

  // Resend OTP
  resendOTP: async (email: string) => {
    return apiRequest("POST", `${API}/auth/resend-otp`, { email });
  },

  // Get current user profile (requires authentication)
  me: async (token: string) => {
    return apiRequest("GET", `${API}/auth/me`, undefined, token);
  },

  // Update profile with phone and address (requires authentication)
  updateProfile: async (
    token: string,
    data: { phone: string; address: string }
  ) => {
    return apiRequest("PATCH", `${API}/auth/profile`, data, token);
  },

  // Check if email is available (using query parameter)
  checkEmailAvailability: async (email: string) => {
    try {
      return await apiRequest("GET", `${API}/auth/email-available?email=${encodeURIComponent(email)}`);
    } catch (err: any) {
      // If query parameter fails, return available=true (assume available)
      console.warn("Email availability check failed:", err);
      return { available: true };
    }
  },

  // Get password strength (using query parameter)
  getPasswordStrength: async (password: string) => {
    try {
      const response = await apiRequest("GET", `${API}/auth/password-strength?password=${encodeURIComponent(password)}`);
      return response;
    } catch (err: any) {
      // If query parameter fails, evaluate locally
      console.warn("Password strength check failed:", err);
      return { strength: evaluateLocalPasswordStrength(password) };
    }
  },

  // Logout (optional backend call)
  logout: async (token: string) => {
    try {
      return await apiRequest("POST", `${API}/auth/logout`, {}, token);
    } catch (err) {
      // Logout is not critical if backend fails
      console.warn("Logout failed:", err);
      return { message: "Logged out" };
    }
  },
};

// Local password strength evaluation (fallback)
function evaluateLocalPasswordStrength(password: string): "weak" | "medium" | "strong" {
  let score = 0;

  if (password.length >= 12) score++;
  if (password.length >= 16) score++;
  if (/[A-Z]/.test(password)) score++;
  if (/[a-z]/.test(password)) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) score++;

  if (score <= 2) return "weak";
  if (score <= 4) return "medium";
  return "strong";
}

// Fetch wrapper with authentication token support
export const apiFetch = async (
  url: string,
  options?: RequestInit & { token?: string }
): Promise<Response> => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
  const headers: Record<string, string> = {};

  // Copy existing headers if they're a plain object
  if (options?.headers && typeof options.headers === 'object' && !Array.isArray(options.headers)) {
    Object.entries(options.headers as Record<string, string>).forEach(([key, value]) => {
      headers[key] = String(value);
    });
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const fullUrl = url.startsWith('http') ? url : `${API}${url}`;

  const res = await fetch(fullUrl, {
    ...options,
    headers,
  });

  if (!res.ok) {
    throw new Error(`API error: ${res.status}`);
  }

  return res;
};

// Wallet API methods
export const walletAPI = {
  createWallet: async (token: string) => {
    return apiRequest("POST", `${API}/api/v1/wallet/create`, {}, token);
  },

  getBalance: async (token: string) => {
    return apiRequest("GET", `${API}/api/v1/wallet/balance`, undefined, token);
  },

  getTransactions: async (token: string) => {
    return apiRequest("GET", `${API}/api/v1/wallet/transactions`, undefined, token);
  },

  withdraw: async (token: string, data: { amount_wei: string; to: string }) => {
    return apiRequest("POST", `${API}/api/v1/wallet/withdraw`, data, token);
  },
};

// API Key API methods
export const apiKeyAPI = {
  create: async (token: string, data: { name: string }) => {
    return apiRequest("POST", `${API}/api/v1/apikeys`, data, token);
  },

  list: async (token: string) => {
    return apiRequest("GET", `${API}/api/v1/apikeys`, undefined, token);
  },

  delete: async (token: string, id: string) => {
    return apiRequest("DELETE", `${API}/api/v1/apikeys/${id}`, {}, token);
  },

  regenerate: async (token: string, id: string) => {
    return apiRequest("POST", `${API}/api/v1/apikeys/${id}/regenerate`, {}, token);
  },
};