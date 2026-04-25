'use client';

import { motion } from "framer-motion";
import { ArrowRight, Shield, Zap, Globe, LayoutDashboard } from "lucide-react";
import Link from "next/link";
import { useAuth } from "@/lib/AuthContext";

export default function Home() {
  const { user } = useAuth();

  return (
    <div className="flex flex-col gap-32 pt-20">
      {/* Hero Section */}
      <section className="relative overflow-hidden">
        <div className="flex flex-col items-center text-center gap-12 relative z-10">


          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2, duration: 0.8 }}
            className="text-6xl md:text-9xl font-bold tracking-tighter leading-[0.85] uppercase max-w-4xl"
          >
            AUTH SYSTEM  <br />
            <span className="text-zinc-400">API + ETH WALLET PR</span>
          </motion.h1>



          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
            className="flex flex-wrap justify-center gap-4"
          >
            {user ? (
              <Link href="/dashboard">
                <button className="h-14 px-8 text-lg rounded-xl bg-black text-white flex items-center">
                  Go to Dashboard
                  <LayoutDashboard className="ml-2" size={20} />
                </button>
              </Link>
            ) : (
              <Link href="/auth/signup">
                <button className="h-14 px-8 text-lg rounded-xl bg-black text-white flex items-center rounded-xl border border-zinc-300">
                  Enter Vanguard
                  <ArrowRight className="ml-2" size={20} />
                </button>
              </Link>
            )}


          </motion.div>
        </div>

        {/* Background */}
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-zinc-200 dark:bg-zinc-900/20 blur-[120px] rounded-full -z-10" />
      </section>

      {/* Features */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {[
          {
            icon: Shield,
            title: "Identity Sovereignty",
            description: "State-Flow Auth onboarding System with OTP service."
          },
          {
            icon: Zap,
            title: "Developer API Service",
            description: "API management with hashed secrets."
          },
          {
            icon: Globe,
            title: "Ethereum Wallet",
            description: "Per-user Ethereum wallets with internal synchronization."
          }
        ].map((feature, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: i * 0.1 }}
            className="p-8 rounded-3xl bg-zinc-50 dark:bg-zinc-900/30 border border-zinc-200 dark:border-zinc-900 flex flex-col gap-6"
          >
            <div className="w-12 h-12 rounded-2xl bg-white dark:bg-zinc-900 flex items-center justify-center">
              <feature.icon size={24} />
            </div>

            <div>
              <h3 className="text-xl font-bold uppercase">
                {feature.title}
              </h3>
              <p className="text-zinc-500">
                {feature.description}
              </p>
            </div>
          </motion.div>
        ))}
      </section>
    </div>
  );
}