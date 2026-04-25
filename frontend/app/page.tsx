'use client';

import { motion } from "motion/react";
import { ArrowRight, Shield, Zap, Globe, LayoutDashboard } from "lucide-react";
import Link from "next/link";
import { Button } from "@/lib/components/ui";
import { useAuth } from "@/lib/AuthContext";

export default function Home() {
  const { user } = useAuth();

  return (
    <div className="flex flex-col gap-32 pt-20">
      {/* Hero Section */}
      <section className="relative overflow-hidden">
        <div className="flex flex-col items-center text-center gap-12 relative z-10">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8 }}
            className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-zinc-100 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-xs font-bold tracking-widest text-zinc-500 uppercase"
          >
            <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
            V0.4.0 Live on Mainnet
          </motion.div>

          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2, duration: 0.8 }}
            className="text-6xl md:text-9xl font-bold tracking-tighter leading-[0.85] uppercase max-w-4xl text-black dark:text-white"
          >
            The Future of <br />
            <span className="text-zinc-400 dark:text-zinc-500">Hybrid Finance.</span>
          </motion.h1>

          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.4, duration: 0.8 }}
            className="text-zinc-500 text-lg md:text-xl max-w-xl font-medium leading-relaxed"
          >
            Seamlessly bridge Web2 performance with Web3 sovereignty.
            Professional API keys meets institutional-grade wallet ledger.
          </motion.p>

          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
            className="flex flex-wrap justify-center gap-4"
          >
            {user ? (
              <Link href="/dashboard">
                <Button size="lg" className="h-14 px-8 text-lg">
                  Go to Dashboard <LayoutDashboard className="ml-2" size={20} />
                </Button>
              </Link>
            ) : (
              <Link href="/auth/signup">
                <Button size="lg" className="h-14 px-8 text-lg">
                  Enter Vanguard <ArrowRight className="ml-2" size={20} />
                </Button>
              </Link>
            )}
            <Button variant="outline" size="lg" className="h-14 px-8 text-lg">
              Read Docs
            </Button>
          </motion.div>
        </div>

        {/* Background Atmosphere */}
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-zinc-200 dark:bg-zinc-900/20 blur-[120px] rounded-full -z-10" />
      </section>

      {/* Features Grid */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {[
          {
            icon: Shield,
            title: "Identity Sovereignty",
            description: "State-machine driven onboarding with cryptographically secured user profiles."
          },
          {
            icon: Zap,
            title: "Developer First",
            description: "Enterprise API management with hashed secrets and granular permissioning."
          },
          {
            icon: Globe,
            title: "Global Ledger",
            description: "Per-user Ethereum wallets with internal synchronization for instant off-chain settlement."
          }
        ].map((feature, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.1 }}
            className="p-8 rounded-3xl bg-zinc-50 dark:bg-zinc-900/30 border border-zinc-200 dark:border-zinc-900 flex flex-col gap-6 group hover:border-zinc-300 dark:hover:border-zinc-700 transition-colors"
          >
            <div className="w-12 h-12 rounded-2xl bg-white dark:bg-zinc-900 flex items-center justify-center text-black dark:text-white border border-zinc-200 dark:border-zinc-800">
              <feature.icon size={24} />
            </div>
            <div className="flex flex-col gap-2">
              <h3 className="text-xl font-bold tracking-tight uppercase text-black dark:text-white">{feature.title}</h3>
              <p className="text-zinc-500 leading-relaxed">{feature.description}</p>
            </div>
          </motion.div>
        ))}
      </section>
    </div>
  );
}
