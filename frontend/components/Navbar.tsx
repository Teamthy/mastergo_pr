'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { motion } from 'motion/react';
import { LayoutDashboard, Wallet, Key, User, LogOut, Sun, Moon } from "lucide-react";
import { cn } from '@/lib/utils';
import { useTheme } from '@/lib/ThemeContext';
import { useAuth } from '@/lib/AuthContext';

export function Navbar() {
    const pathname = usePathname();
    const router = useRouter();
    const { theme, toggleTheme } = useTheme();
    const { user, logout } = useAuth();

    const isAuthPage = pathname.startsWith('/auth');
    if (isAuthPage) return null;

    const navItems = [
        { label: "Dashboard", icon: LayoutDashboard, path: "/dashboard" },
        { label: "Wallet", icon: Wallet, path: "/dashboard/wallet" },
        { label: "API Keys", icon: Key, path: "/dashboard/apikeys" },
        { label: "Profile", icon: User, path: "/dashboard/profile" },
    ];

    const handleLogout = () => {
        logout();
        router.push('/');
    };

    return (
        <>
            <nav className="fixed top-0 left-0 right-0 z-50 border-b border-zinc-200 dark:border-zinc-900 bg-white/50 dark:bg-black/50 backdrop-blur-xl">
                <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
                    <Link href="/" className="flex items-center gap-2 group">
                        <div className="w-8 h-8 bg-black dark:bg-white rounded-lg flex items-center justify-center transform group-hover:rotate-12 transition-transform">
                            <div className="w-4 h-4 bg-white dark:bg-black rounded-sm" />
                        </div>
                        <span className="font-bold tracking-tight text-xl text-black dark:text-white">VANGUARD</span>
                    </Link>

                    <div className="hidden md:flex items-center gap-1">
                        {navItems.map((item) => (
                            <Link
                                key={item.path}
                                href={item.path}
                                className={cn(
                                    "px-4 py-2 rounded-full text-sm font-medium transition-all flex items-center gap-2",
                                    pathname === item.path
                                        ? "bg-zinc-100 dark:bg-zinc-900 text-black dark:text-white"
                                        : "text-zinc-500 hover:text-black dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-900/50"
                                )}
                            >
                                <item.icon size={16} />
                                {item.label}
                            </Link>
                        ))}
                    </div>

                    <div className="flex items-center gap-2">
                        <button
                            onClick={toggleTheme}
                            className="p-2 text-zinc-500 hover:text-black dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-900/50 rounded-full transition-all"
                        >
                            {theme === "dark" ? <Sun size={20} /> : <Moon size={20} />}
                        </button>
                        {user && (
                            <button
                                onClick={handleLogout}
                                className="p-2 text-zinc-500 hover:text-red-500 hover:bg-red-500/10 rounded-full transition-all"
                            >
                                <LogOut size={20} />
                            </button>
                        )}
                    </div>
                </div>
            </nav>

            <footer className="fixed bottom-0 left-0 right-0 md:hidden border-t border-zinc-200 dark:border-zinc-900 bg-white/80 dark:bg-black/80 backdrop-blur-xl z-50">
                <div className="flex justify-around items-center h-16">
                    {navItems.map((item) => (
                        <Link
                            key={item.path}
                            href={item.path}
                            className={cn(
                                "p-2 flex flex-col items-center gap-1",
                                pathname === item.path ? "text-black dark:text-white" : "text-zinc-400 dark:text-zinc-600"
                            )}
                        >
                            <item.icon size={20} />
                            <span className="text-[10px] uppercase tracking-widest font-bold">{item.label.split(' ')[0]}</span>
                        </Link>
                    ))}
                </div>
            </footer>
        </>
    );
}
