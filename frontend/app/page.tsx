'use client';

import { motion } from "framer-motion";
import { ArrowRight, LayoutDashboard } from "lucide-react";
import Link from "next/link";
import { useAuth } from "@/lib/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";


export default function HomePage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  // Redirect authenticated users with completed onboarding to dashboard
  useEffect(() => {
    if (!loading && user && user.onboarding_status === "COMPLETED") {
      router.replace("/dashboard/wallet");
    }
  }, [user, loading, router]);

  if (loading) return null;

  return (
    <div className="flex flex-col gap-32 pt-20">



      {/* Landing Page (only if logged in or before auth) */}
      <section className="relative overflow-hidden">
        <div className="flex flex-col items-center text-center gap-12 relative z-10">

          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2, duration: 0.8 }}
            className="text-6xl md:text-9xl font-bold tracking-tighter leading-[0.85] uppercase max-w-4xl"
          >
            AUTH SYSTEM <br />
            <span className="text-zinc-400">
              API + ETH<br />
            </span>

            <span className="text-zinc-600">

              WALLET PR
            </span>
          </motion.h1>

          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
            className="flex flex-wrap justify-center gap-4"
          >
            {user ? (
              <Link href="/dashboard/wallet">
                <button className="h-14 px-8 text-lg rounded-xl bg-black text-white flex items-center">
                  Go to Dashboard
                  <LayoutDashboard className="ml-2" size={20} />
                </button>
              </Link>
            ) : (
              <Link href="/auth/login">
                <button className="h-14 px-8 text-lg rounded-xl bg-white text-black flex items-center">
                  ENTER MY APP
                  <ArrowRight className="ml-2" size={20} />
                </button>
              </Link>
            )}
          </motion.div>
        </div>

        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-zinc-200 dark:bg-zinc-900/20 blur-[120px] rounded-full -z-10" />
      </section>

      {/* Features */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {[
          {

            title: "STATE AUTH FLOW",
            description: "State-Flow Auth onboarding System with OTP service."
          },
          {
            title: "Developer API Service",
            description: "API management with hashed secrets."
          },
          {

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