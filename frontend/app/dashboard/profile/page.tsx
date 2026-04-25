'use client';

import { useState } from "react";
import { motion } from "motion/react";
import { User, Mail, Phone, MapPin, Shield, Save, Check, ShieldAlert } from "lucide-react";
import { Button, Input } from "@/lib/components/ui";

export default function ProfilePage() {
    const [loading, setLoading] = useState(false);
    const [saved, setSaved] = useState(false);

    const [profile, setProfile] = useState({
        name: "Johnathan Doe",
        email: "john.doe@vanguard.crypto",
        phone: "+1 (555) 000-0000",
        address: "7th Block, Silicon Valley, CA, USA",
        tier: "Pro Merchant",
        id: "USER_84144_992"
    });

    const handleSave = async () => {
        setLoading(true);
        await new Promise(r => setTimeout(r, 1200));
        setLoading(false);
        setSaved(true);
        setTimeout(() => setSaved(false), 3000);
    };

    return (
        <div className="flex flex-col gap-10 max-w-4xl">
            <header className="flex flex-col gap-1">
                <h1 className="text-4xl font-bold tracking-tighter uppercase text-black dark:text-white">Identity Console</h1>
                <p className="text-zinc-500 font-medium tracking-tight">Manage your verified global profile and security credentials.</p>
            </header>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                {/* Profile Sidebar */}
                <div className="flex flex-col gap-6">
                    <div className="p-8 rounded-3xl bg-zinc-50 dark:bg-zinc-900/40 border border-zinc-200 dark:border-zinc-900 flex flex-col items-center text-center gap-4">
                        <div className="w-24 h-24 rounded-full bg-gradient-to-tr from-zinc-200 to-zinc-100 dark:from-zinc-800 dark:to-zinc-700 p-1">
                            <div className="w-full h-full rounded-full bg-white dark:bg-black flex items-center justify-center text-black dark:text-white border border-zinc-200 dark:border-zinc-800">
                                <User size={40} />
                            </div>
                        </div>
                        <div className="flex flex-col">
                            <span className="font-bold text-lg uppercase tracking-tight text-black dark:text-white">{profile.name}</span>
                            <span className="text-[10px] uppercase font-bold tracking-widest text-emerald-600 dark:text-emerald-500 bg-emerald-500/10 px-3 py-1 rounded-full mt-2 self-center">
                                {profile.tier}
                            </span>
                        </div>
                        <div className="w-full h-px bg-zinc-200 dark:bg-zinc-800 my-2" />
                        <div className="flex flex-col gap-1">
                            <span className="text-[10px] uppercase font-bold text-zinc-400 dark:text-zinc-600">Account ID</span>
                            <code className="text-xs font-mono text-zinc-500">{profile.id}</code>
                        </div>
                    </div>

                    <div className="p-6 rounded-3xl border border-zinc-200 dark:border-zinc-900 flex flex-col gap-4 bg-zinc-50 dark:bg-transparent">
                        <h3 className="text-xs font-bold uppercase tracking-widest text-zinc-500">Security Pulse</h3>
                        <div className="flex flex-col gap-3">
                            <div className="flex items-center justify-between border-b border-zinc-100 dark:border-zinc-900 pb-2">
                                <span className="text-xs text-zinc-500 dark:text-zinc-400">2FA Status</span>
                                <span className="text-[10px] font-bold text-emerald-500 uppercase">Active</span>
                            </div>
                            <div className="flex items-center justify-between">
                                <span className="text-xs text-zinc-500 dark:text-zinc-400">KYC Level</span>
                                <span className="text-[10px] font-bold text-emerald-500 uppercase">Tier 3</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Edit Form */}
                <div className="md:col-span-2 flex flex-col gap-8">
                    <div className="p-8 rounded-3xl bg-zinc-50 dark:bg-zinc-900/20 border border-zinc-200 dark:border-zinc-900 flex flex-col gap-8">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div className="flex flex-col gap-3">
                                <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4 flex items-center gap-2">
                                    <User size={12} /> Full Name
                                </label>
                                <Input
                                    value={profile.name}
                                    onChange={(e) => setProfile({ ...profile, name: e.target.value })}
                                />
                            </div>
                            <div className="flex flex-col gap-3">
                                <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4 flex items-center gap-2">
                                    <Mail size={12} /> Email Address
                                </label>
                                <Input
                                    value={profile.email}
                                    onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                                />
                            </div>
                            <div className="flex flex-col gap-3">
                                <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4 flex items-center gap-2">
                                    <Phone size={12} /> Contact Number
                                </label>
                                <Input
                                    value={profile.phone}
                                    onChange={(e) => setProfile({ ...profile, phone: e.target.value })}
                                />
                            </div>
                            <div className="flex flex-col gap-3">
                                <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4 flex items-center gap-2">
                                    <MapPin size={12} /> Physical Address
                                </label>
                                <Input
                                    value={profile.address}
                                    onChange={(e) => setProfile({ ...profile, address: e.target.value })}
                                />
                            </div>
                        </div>

                        <div className="flex justify-end pt-4 border-t border-zinc-200 dark:border-zinc-800">
                            <Button
                                size="lg"
                                className="min-w-[160px]"
                                disabled={loading}
                                onClick={handleSave}
                            >
                                {loading ? (
                                    "Syncing..."
                                ) : saved ? (
                                    <><Check size={18} className="mr-2" /> Saved</>
                                ) : (
                                    <><Save size={18} className="mr-2" /> Update Profile</>
                                )}
                            </Button>
                        </div>
                    </div>

                    {/* Security Settings Area */}
                    <div className="grid grid-cols-1 gap-4">
                        <div className="p-6 rounded-3xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/10 flex items-center justify-between group cursor-pointer hover:border-zinc-300 dark:hover:border-zinc-700 transition-colors">
                            <div className="flex items-center gap-4">
                                <div className="w-10 h-10 rounded-xl bg-zinc-900 flex items-center justify-center text-zinc-500">
                                    <Shield size={20} />
                                </div>
                                <div className="flex flex-col">
                                    <span className="font-bold text-sm uppercase text-black dark:text-white">Change Password</span>
                                    <span className="text-[10px] text-zinc-500 dark:text-zinc-600 font-medium">Last changed 40 days ago</span>
                                </div>
                            </div>
                            <Button variant="ghost" size="sm">Manage</Button>
                        </div>

                        <div className="p-6 rounded-3xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/10 flex items-center justify-between group cursor-pointer hover:border-zinc-300 dark:hover:border-zinc-700 transition-colors">
                            <div className="flex items-center gap-4">
                                <div className="w-10 h-10 rounded-xl bg-orange-500/10 flex items-center justify-center text-orange-500">
                                    <ShieldAlert size={20} />
                                </div>
                                <div className="flex flex-col">
                                    <span className="font-bold text-sm uppercase text-black dark:text-white">Revoke Sessions</span>
                                    <span className="text-[10px] text-zinc-500 dark:text-zinc-600 font-medium">Kill all active web and mobile logins</span>
                                </div>
                            </div>
                            <Button variant="ghost" size="sm" className="text-orange-500">Revoke</Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
