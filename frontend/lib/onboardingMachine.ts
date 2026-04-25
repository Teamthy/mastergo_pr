import { createMachine, assign } from "xstate"
import {
    OnboardingState,
    OnboardingContext,
    OnboardingEvent,
} from "./types"

export const onboardingMachine = createMachine(
    {
        id: "onboarding",

        types: {} as {
            context: OnboardingContext
            events: OnboardingEvent
        },

        initial: OnboardingState.START,

        context: {
            email: undefined,
            otp: undefined,
            name: undefined,
            contact: undefined,
            address: undefined,
        },

        states: {
            [OnboardingState.START]: {
                on: {
                    SUBMIT_EMAIL: {
                        target: OnboardingState.EMAIL_ENTERED,
                        actions: "assignEmail",
                    },
                },
            },

            [OnboardingState.EMAIL_ENTERED]: {
                on: {
                    VERIFY_OTP: {
                        target: OnboardingState.EMAIL_VERIFIED,
                        actions: "assignOtp",
                    },
                    RESET: {
                        target: OnboardingState.START,
                    },
                },
            },

            [OnboardingState.EMAIL_VERIFIED]: {
                on: {
                    SUBMIT_NAME: {
                        target: OnboardingState.PROFILE_NAME,
                        actions: "assignName",
                    },
                },
            },

            [OnboardingState.PROFILE_NAME]: {
                on: {
                    SUBMIT_CONTACT: {
                        target: OnboardingState.PROFILE_CONTACT,
                        actions: "assignContact",
                    },
                },
            },

            [OnboardingState.PROFILE_CONTACT]: {
                on: {
                    SUBMIT_ADDRESS: {
                        target: OnboardingState.COMPLETED,
                        actions: "assignAddress",
                    },
                },
            },

            [OnboardingState.COMPLETED]: {
                type: "final",
            },
        },
    },
    {
        actions: {
            assignEmail: assign({
                email: ({ event }) =>
                    event.type === "SUBMIT_EMAIL" ? event.email : undefined,
            }),

            assignOtp: assign({
                otp: ({ event }) =>
                    event.type === "VERIFY_OTP" ? event.otp : undefined,
            }),

            assignName: assign({
                name: ({ event }) =>
                    event.type === "SUBMIT_NAME" ? event.name : undefined,
            }),

            assignContact: assign({
                contact: ({ event }) =>
                    event.type === "SUBMIT_CONTACT" ? event.contact : undefined,
            }),

            assignAddress: assign({
                address: ({ event }) =>
                    event.type === "SUBMIT_ADDRESS" ? event.address : undefined,
            }),
        },
    }
)
