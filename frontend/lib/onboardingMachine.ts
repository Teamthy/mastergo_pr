import { createMachine, assign } from "xstate";

export const onboardingMachine = createMachine({
    id: "signup",
    initial: "email",
    context: {
        email: "",
        password: "",
        otp: "",
    },

    states: {
        email: {
            on: {
                SUBMIT_EMAIL: {
                    target: "otp",
                    actions: assign({
                        email: (_, e: any) => e.email,
                        password: (_, e: any) => e.password,
                    }),
                },
            },
        },

        otp: {
            on: {
                VERIFY_OTP: {
                    target: "completed",
                    actions: assign({
                        otp: (_, e: any) => e.otp,
                    }),
                },
            },
        },

        completed: {
            type: "final",
        },
    },
});