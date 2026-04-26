import { createMachine, assign } from "xstate";
import { OnboardingState } from "./types";

export const onboardingMachine = createMachine({
    id: "onboarding",
    initial: OnboardingState.START,

    context: {
        email: "",
        name: "",
        contact: "",
        address: "",
    },

    states: {
        [OnboardingState.START]: {
            on: {
                SUBMIT_EMAIL: {
                    target: OnboardingState.EMAIL_ENTERED,
                    actions: assign({
                        email: (_, e: any) => e.email,
                    }),
                },
            },
        },

        [OnboardingState.EMAIL_ENTERED]: {
            on: {
                VERIFY_OTP: {
                    target: OnboardingState.EMAIL_VERIFIED,
                },
            },
        },

        [OnboardingState.EMAIL_VERIFIED]: {
            on: {
                SUBMIT_NAME: {
                    target: OnboardingState.PROFILE_NAME,
                    actions: assign({
                        name: (_, e: any) => e.name,
                    }),
                },
            },
        },

        [OnboardingState.PROFILE_NAME]: {
            on: {
                SUBMIT_CONTACT: {
                    target: OnboardingState.PROFILE_CONTACT,
                    actions: assign({
                        contact: (_, e: any) => e.contact,
                    }),
                },
            },
        },

        [OnboardingState.PROFILE_CONTACT]: {
            on: {
                SUBMIT_ADDRESS: {
                    target: OnboardingState.COMPLETED,
                    actions: assign({
                        address: (_, e: any) => e.address,
                    }),
                },
            },
        },

        [OnboardingState.COMPLETED]: {
            type: "final",
        },
    },
});