'use client';

import { createContext, useContext, useEffect, useState } from "react";
import { authAPI } from "./api";

export interface User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone?: string;
  address?: string;
  email_verified: boolean;
  onboarding_status: string;
  created_at: string;
  updated_at: string;
}

type AuthContextType = {
  user: User | null;
  token: string | null;
  loading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  signup: (
    firstName: string,
    lastName: string,
    email: string,
    password: string,
    confirmPassword: string
  ) => Promise<void>;
  verifyEmail: (email: string, otp: string) => Promise<void>;
  resendOTP: (email: string) => Promise<void>;
  updateProfile: (phone: string, address: string) => Promise<void>;
  logout: () => void;
  checkEmailAvailable: (email: string) => Promise<boolean>;
  getPasswordStrength: (password: string) => Promise<"weak" | "medium" | "strong">;
};

const AuthContext = createContext<AuthContextType>({
  user: null,
  token: null,
  loading: true,
  error: null,
  login: async () => {},
  signup: async () => {},
  verifyEmail: async () => {},
  resendOTP: async () => {},
  updateProfile: async () => {},
  logout: () => {},
  checkEmailAvailable: async () => false,
  getPasswordStrength: async () => "weak",
});

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Restore session on mount
  useEffect(() => {
    const storedToken = localStorage.getItem("token");

    if (!storedToken) {
      setLoading(false);
      return;
    }

    setToken(storedToken);

    authAPI
      .me(storedToken)
      .then((userData) => {
        setUser(userData);
        setError(null);
      })
      .catch((err) => {
        console.error("Failed to restore session:", err);
        localStorage.removeItem("token");
        setToken(null);
        setUser(null);
      })
      .finally(() => setLoading(false));
  }, []);

  const login = async (email: string, password: string) => {
    try {
      setError(null);
      const response = await authAPI.login(email, password);
      localStorage.setItem("token", response.token);
      setToken(response.token);
      setUser(response.user);
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.message || "Login failed";
      setError(errorMsg);
      throw err;
    }
  };

  const signup = async (
    firstName: string,
    lastName: string,
    email: string,
    password: string,
    confirmPassword: string
  ) => {
    try {
      setError(null);
      await authAPI.signup({
        first_name: firstName,
        last_name: lastName,
        email,
        password,
        confirm_password: confirmPassword,
      });
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.message || "Signup failed";
      setError(errorMsg);
      throw err;
    }
  };

  const verifyEmail = async (email: string, otp: string) => {
    try {
      // Don't set user/token yet - they need to verify email first
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.message || "Signup failed";
      setError(errorMsg);
      throw err;
    }
  };

  const verifyEmail = async (email: string, otp: string) => {
    try {
      setError(null);
      await authAPI.verifyEmail(email, otp);
      // After verification, auto-login
      // This will be handled by the frontend flow
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.message || "Verification failed";
      setError(errorMsg);
      throw err;
    }
  };

  const resendOTP = async (email: string) => {
    try {
      setError(null);
      await authAPI.resendOTP(email);
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.message || "Resend failed";
      setError(errorMsg);
      throw err;
    }
  };

  const updateProfile = async (phone: string, address: string) => {
    if (!token) {
      throw new Error("Not authenticated");
    }

    try {
      setError(null);
      const userData = await authAPI.updateProfile(token, {
        phone,
        address,
      });
      setUser(userData);
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.message || "Profile update failed";
      setError(errorMsg);
      throw err;
    }
  };

  const logout = () => {
    localStorage.removeItem("token");
    setToken(null);
    setUser(null);
    setError(null);
  };

  const checkEmailAvailable = async (email: string): Promise<boolean> => {
    try {
      const response = await authAPI.checkEmailAvailability(email);
      return response.available;
    } catch {
      return false;
    }
  };

 

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        loading,
        error,
        login,
        signup,
        verifyEmail,
        resendOTP,
        updateProfile,
        logout,
        checkEmailAvailable,
       
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);