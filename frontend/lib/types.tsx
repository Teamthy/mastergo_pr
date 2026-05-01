export enum OnboardingState {
    START = "START",
    EMAIL_ENTERED = "EMAIL_ENTERED",
    EMAIL_VERIFIED = "EMAIL_VERIFIED",
    PROFILE_NAME = "PROFILE_NAME",
    PROFILE_CONTACT = "PROFILE_CONTACT",
    COMPLETED = "COMPLETED",
}

export interface OnboardingContext {
    email?: string
    otp?: string
    name?: string
    contact?: string
    address?: string
}

export type OnboardingEvent =
    | { type: "SUBMIT_EMAIL"; email: string }
    | { type: "VERIFY_OTP"; otp: string }
    | { type: "SUBMIT_NAME"; name: string }
    | { type: "SUBMIT_CONTACT"; contact: string }
    | { type: "SUBMIT_ADDRESS"; address: string }
    | { type: "RESET" }

// User types
export interface User {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    phone?: string;
    address?: string;
    email_verified: boolean;
    onboarding_status: string;
    created_at: string;
    updated_at: string;
    last_login_at?: string;
}

// API Key types
export interface ApiKey {
    id: string;
    userId: string;
    name: string;
    publicKey: string;
    createdAt: string;
    revokedAt?: string;
}

// Wallet types
export interface Wallet {
    id: string;
    userId: string;
    address?: string;
    publicKey?: string;
    balance?: string;
    createdAt: string;
    updatedAt?: string;
}

export interface Transaction {
    id: string;
    user_id?: string;
    walletId?: string;
    tx_hash?: string;
    txHash?: string;
    type: 'deposit' | 'withdrawal';
    amount_wei?: string;
    amount?: string;
    to?: string;
    status: 'pending' | 'confirmed' | 'failed';
    created_at?: string;
    createdAt?: string;
    updated_at?: string;
}