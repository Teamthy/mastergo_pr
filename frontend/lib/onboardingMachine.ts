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
          actions: assign({
            firstName: (_, e: any) => e?.firstName ?? "",
            lastName: (_, e: any) => e?.lastName ?? "",
            error: () => "",
          }),
        },
        SET_ERROR: {
          actions: assign({
            error: (_, e: any) => e?.error ?? "",
          }),
        },
      },
    },

    credentials: {
      on: {
        BACK: "name",
        SUBMIT_CREDENTIALS: {
          target: "otp",
          actions: assign({
            email: (_, e: any) => e?.email ?? "",
            password: (_, e: any) => e?.password ?? "",
            confirmPassword: (_, e: any) => e?.confirmPassword ?? "",
            error: () => "",
          }),
        },
        UPDATE_PASSWORD_STRENGTH: {
          actions: assign({
            passwordStrength: (_, e: any) => e?.strength ?? "weak",
          }),
        },
        SET_ERROR: {
          actions: assign({
            error: (_, e: any) => e?.error ?? "",
          }),
        },
      },
    },

    otp: {
      on: {
        BACK: "credentials",
        VERIFY_OTP: {
          target: "profile",
          actions: assign({
            otp: (_, e: any) => e?.otp ?? "",
            error: () => "",
          }),
        },
        RESEND_OTP: {
          actions: assign({
            error: () => "",
          }),
        },
        SET_ERROR: {
          actions: assign({
            error: (_, e: any) => e?.error ?? "",
          }),
        },
      },
    },

    profile: {
      on: {
        BACK: "otp",
        SUBMIT_PROFILE: {
          target: "completed",
          actions: assign({
            phone: (_, e: any) => e?.phone ?? "",
            address: (_, e: any) => e?.address ?? "",
            error: () => "",
          }),
        },
        SET_ERROR: {
          actions: assign({
            error: (_, e: any) => e?.error ?? "",
          }),
        },
      },
    },

    completed: {
      type: "final",
    },
  },
});