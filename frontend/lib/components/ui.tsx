'use client';

import * as React from "react";
import { cn } from "@/lib/utils";

export interface ButtonProps
    extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: "default" | "outline";
    size?: "sm" | "md" | "lg";
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
    ({ className, variant = "default", size = "md", ...props }, ref) => {
        return (
            <button
                ref={ref}
                className={cn(
                    "inline-flex items-center justify-center rounded-xl font-semibold transition-colors",
                    "focus:outline-none focus:ring-2 focus:ring-offset-2",
                    variant === "default" &&
                    "bg-black text-white hover:bg-zinc-800 dark:bg-white dark:text-black",
                    variant === "outline" &&
                    "border border-zinc-200 dark:border-zinc-800 hover:bg-zinc-100 dark:hover:bg-zinc-900",
                    size === "sm" && "h-8 px-3 text-sm",
                    size === "md" && "h-10 px-4",
                    size === "lg" && "h-14 px-8 text-lg",
                    className
                )}
                {...props}
            />
        );
    }
);

Button.displayName = "Button";

export interface InputProps
    extends React.InputHTMLAttributes<HTMLInputElement> { }

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
    ({ className, ...props }, ref) => {
        return (
            <input
                ref={ref}
                className={cn(
                    "w-full h-14 px-5 rounded-2xl",
                    "bg-white dark:bg-zinc-900",
                    "border border-zinc-200 dark:border-zinc-800",
                    "focus:outline-none focus:ring-2 focus:ring-black dark:focus:ring-white",
                    className
                )}
                {...props}
            />
        );
    }
);

Input.displayName = "Input";