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