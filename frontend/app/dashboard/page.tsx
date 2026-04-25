'use client';

import { motion } from "motion/react";
import { Wallet, Key, Activity } from "lucide-react";
import { Button } from "@/lib/components/ui";
import { cn } from "@/lib/utils";

export default function DashboardPage() {
    const stats = [
        { label: "Total Asset Value", value: "$428,192.00", change: "+12.4%", icon: Wallet },
        { label: "Active API Keys", value: "08", change: "Stable", icon: Key },
        { label: "Network Activity", value: "High", change: "99.9% Uptime", icon: Activity },
    ];

    return (
        <div className="flex flex-col gap-12">
            <header className="flex flex-col gap-1">
                <h1 className="text-4xl font-bold tracking-tighter uppercase">Systems Overview</h1>
                <p className="text-zinc-500 font-medium tracking-tight">Welcome back. All systems are operational.</p>
            </header>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {stats.map((stat, i) => (
                    <motion.div
                        key={i}
                        initial={{ opacity: 0, scale: 0.95 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ delay: i * 0.1 }}
                        className="p-8 rounded-3xl bg-zinc-50 dark:bg-zinc-900/40 border border-zinc-200 dark:border-zinc-900 flex flex-col gap-4"
                    >
                        <div className="flex justify-between items-start">
                            <div className="w-10 h-10 rounded-xl bg-white dark:bg-black border border-zinc-200 dark:border-zinc-800 flex items-center justify-center text-zinc-500 dark:text-zinc-400">
                                <stat.icon size={20} />
                            </div>
                            <span className={cn(
                                "text-[10px] uppercase font-bold tracking-widest px-2 py-1 rounded-md",
                                stat.change.startsWith('+') ? "bg-emerald-500/10 text-emerald-500" : "bg-zinc-200 dark:bg-zinc-800 text-zinc-500 dark:text-zinc-400"
                            )}>
                                {stat.change}
                            </span>
                        </div>
                        <div>
                            <p className="text-xs uppercase tracking-widest font-bold text-zinc-500">{stat.label}</p>
                            <p className="text-3xl font-bold tracking-tight mt-1 text-black dark:text-white">{stat.value}</p>
                        </div>
                    </motion.div>
                ))}
            </div>

            {/* Charts / Activity Split */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                <div className="md:col-span-2 rounded-3xl border border-zinc-200 dark:border-zinc-900 bg-zinc-50 dark:bg-zinc-900/20 p-8 flex flex-col gap-8 h-[400px]">
                    <div className="flex justify-between items-center">
                        <h2 className="text-sm font-bold uppercase tracking-widest text-black dark:text-white">Market Performance</h2>
                        <Button variant="ghost" size="sm">Last 30 Days</Button>
                    </div>

                    <div className="flex-1 flex items-end gap-2 px-4">
                        {[40, 20, 60, 45, 90, 70, 85, 30, 40, 60, 50, 75, 95].map((h, i) => (
                            <motion.div
                                key={i}
                                initial={{ height: 0 }}
                                animate={{ height: `${h}%` }}
                                transition={{ delay: i * 0.05, duration: 1 }}
                                className="flex-1 bg-gradient-to-t from-zinc-300 to-zinc-200 dark:from-zinc-800 dark:to-zinc-700 rounded-t-sm"
                            />
                        ))}
                    </div>
                </div>

                <div className="rounded-3xl border border-zinc-200 dark:border-zinc-900 bg-zinc-50 dark:bg-zinc-900/20 p-8 flex flex-col gap-6">
                    <h2 className="text-sm font-bold uppercase tracking-widest text-black dark:text-white">Recent Security Logs</h2>
                    <div className="flex flex-col gap-4">
                        {[
                            { event: "New API Key Created", time: "2m ago", status: "Success" },
                            { event: "Login from New IP", time: "1h ago", status: "Warning" },
                            { event: "Wallet Sync Successful", time: "4h ago", status: "Success" },
                            { event: "Withdrawal Completed", time: "1d ago", status: "Success" },
                        ].map((log, i) => (
                            <div key={i} className="flex justify-between items-center border-b border-zinc-100 dark:border-zinc-900 pb-4 last:border-0 last:pb-0">
                                <div className="flex flex-col gap-0.5">
                                    <span className="text-sm font-medium text-black dark:text-white">{log.event}</span>
                                    <span className="text-[10px] uppercase font-bold text-zinc-400 dark:text-zinc-600">{log.time}</span>
                                </div>
                                <span className={cn(
                                    "w-2 h-2 rounded-full",
                                    log.status === "Success" ? "bg-emerald-500" : "bg-amber-500"
                                )} />
                            </div>
                        ))}
                    </div>
                    <Button variant="outline" size="sm" className="mt-4 w-full">View Audit Log</Button>
                </div>
            </div>
        </div>
    );
}
