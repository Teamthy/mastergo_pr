import { createMachine, assign } from "xstate";

export interface SignUpContext {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  confirmPassword: string;
  otp: string;
  phone: string;
  address: string;
  error: string;
  passwordStrength: "weak" | "medium" | "strong";
}

export const onboardingMachine = createMachine({
  id: "signup",
  initial: "name",
  context: {
    firstName: "",
    lastName: "",
    email: "",
    password: "",
    confirmPassword: "",
    otp: "",
    phone: "",
    address: "",
    error: "",
    passwordStrength: "weak",
  } as SignUpContext,

  states: {
    name: {
      on: {
        SUBMIT_NAME: {
          target: "credentials",
          actions: assign(({ event }: any) => ({
            firstName: event.firstName?.trim() ?? "",
            lastName: event.lastName?.trim() ?? "",
            error: "",
          })),
        },
        SET_ERROR: {
          actions: assign(({ event }: any) => ({
            error: event.error ?? "",
          })),
        },
      },
    },

    credentials: {
      on: {
        BACK: "name",
        SUBMIT_CREDENTIALS: {
          target: "otp",
          actions: assign(({ event }: any) => ({
            email: event.email ?? "",
            password: event.password ?? "",
            confirmPassword: event.confirmPassword ?? "",
            error: "",
          })),
        },
        UPDATE_PASSWORD_STRENGTH: {
          actions: assign(({ event }: any) => ({
            passwordStrength: event.strength ?? "weak",
          })),
        },
        SET_ERROR: {
          actions: assign(({ event }: any) => ({
            error: event.error ?? "",
          })),
        },
      },
    },

    otp: {
      on: {
        BACK: "credentials",
        VERIFY_OTP: {
          target: "profile",
          actions: assign(({ event }: any) => ({
            otp: event.otp ?? "",
            error: "",
          })),
        },
        RESEND_OTP: {
          actions: assign(() => ({
            error: "",
          })),
        },
        SET_ERROR: {
          actions: assign(({ event }: any) => ({
            error: event.error ?? "",
          })),
        },
      },
    },

    profile: {
      on: {
        BACK: "otp",
        SUBMIT_PROFILE: {
          target: "completed",
          actions: assign(({ event }: any) => ({
            phone: event.phone ?? "",
            address: event.address ?? "",
            error: "",
          })),
        },
        SET_ERROR: {
          actions: assign(({ event }: any) => ({
            error: event.error ?? "",
          })),
        },
      },
    },

    completed: {
      type: "final",
    },
  },
});