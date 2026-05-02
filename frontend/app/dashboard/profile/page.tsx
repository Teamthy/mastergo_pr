'use client';

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "motion/react";
import { User, Mail, Phone, MapPin, Shield, Save, Check, ShieldAlert, LogOut, Lock, Eye, EyeOff, AlertCircle } from "lucide-react";
import { Button, Input } from "@/lib/components/ui";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/AuthContext";

export default function ProfilePage() {
    const router = useRouter();
    const { user, logout } = useAuth();
    const [loading, setLoading] = useState(false);
    const [saved, setSaved] = useState(false);
    const [showPasswordModal, setShowPasswordModal] = useState(false);
    const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);
    const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
    const [passwordForm, setPasswordForm] = useState({
        current: '',
        new: '',
        confirm: ''
    });
    const [showPasswords, setShowPasswords] = useState({
        current: false,
        new: false,
        confirm: false
    });
    const [passwordError, setPasswordError] = useState('');

    const [profile, setProfile] = useState({
        name: user ? `${user.first_name} ${user.last_name}` : "User",
        email: user?.email || "user@example.com",
        phone: user?.phone || "",
        address: user?.address || "",
        tier: "Pro Merchant",
        id: user?.id || "USER_UNKNOWN"
    });

    useEffect(() => {
        if (user) {
            setProfile({
                name: `${user.first_name} ${user.last_name}`,
                email: user.email,
                phone: user.phone || "",
                address: user.address || "",
                tier: "Pro Merchant",
                id: user.id
            });
        }
    }, [user]);

    const handleSave = async () => {
        setLoading(true);
        await new Promise(r => setTimeout(r, 1200));
        setLoading(false);
        setSaved(true);
        setTimeout(() => setSaved(false), 3000);
    };

    const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setAvatarPreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        }
    };

    const handlePasswordSubmit = async () => {
        setPasswordError('');

        if (!passwordForm.current || !passwordForm.new || !passwordForm.confirm) {
            setPasswordError('All password fields are required');
            return;
        }

        if (passwordForm.new.length < 8) {
            setPasswordError('New password must be at least 8 characters');
            return;
        }

        if (passwordForm.new !== passwordForm.confirm) {
            setPasswordError('New passwords do not match');
            return;
        }

        setLoading(true);
        try {
            // Simulate password change
            await new Promise(r => setTimeout(r, 1200));
            setPasswordForm({ current: '', new: '', confirm: '' });
            setShowPasswordModal(false);
            setSaved(true);
            setTimeout(() => setSaved(false), 3000);
        } catch (err) {
            setPasswordError('Failed to change password. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    const handleLogout = () => {
        setShowLogoutConfirm(false);
        setLoading(true);
        // Call the logout function from AuthContext
        logout();
        // Redirect will be handled by ProtectedRoute after user is cleared
        setTimeout(() => {
            setLoading(false);
        }, 500);
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
                        <div className="relative w-24 h-24 rounded-full bg-gradient-to-tr from-zinc-200 to-zinc-100 dark:from-zinc-800 dark:to-zinc-700 p-1 group cursor-pointer">
                            <img
                                src={avatarPreview || "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2'%3E%3Cpath d='M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2'/%3E%3Ccircle cx='12' cy='7' r='4'/%3E%3C/svg%3E"}
                                alt="Avatar"
                                className="w-full h-full rounded-full bg-white dark:bg-black flex items-center justify-center text-black dark:text-white border border-zinc-200 dark:border-zinc-800 object-cover"
                            />
                            <label className="absolute inset-0 flex items-center justify-center bg-black/50 rounded-full opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer">
                                <User size={20} className="text-white" />
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={handleAvatarChange}
                                    className="hidden"
                                />
                            </label>
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

                    {/* Logout Button */}
                    <Button
                        variant="outline"
                        className="w-full text-red-500 hover:text-red-600 border-red-500/20 hover:border-red-500/40"
                        onClick={() => setShowLogoutConfirm(true)}
                    >
                        <LogOut size={16} className="mr-2" />
                        Logout
                    </Button>
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
                        <div onClick={() => setShowPasswordModal(true)} className="p-6 rounded-3xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/10 flex items-center justify-between group cursor-pointer hover:border-zinc-300 dark:hover:border-zinc-700 transition-colors">
                            <div className="flex items-center gap-4">
                                <div className="w-10 h-10 rounded-xl bg-zinc-900 dark:bg-white flex items-center justify-center text-zinc-500 dark:text-black">
                                    <Shield size={20} />
                                </div>
                                <div className="flex flex-col">
                                    <span className="font-bold text-sm uppercase text-black dark:text-white">Change Password</span>
                                    <span className="text-[10px] text-zinc-500 dark:text-zinc-600 font-medium">Update your security credentials</span>
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

            {/* Password Change Modal */}
            <AnimatePresence>
                {showPasswordModal && (
                    <div className="fixed inset-0 z-[60] flex items-center justify-center p-6">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            onClick={() => setShowPasswordModal(false)}
                            className="absolute inset-0 bg-black/50 dark:bg-black/80 backdrop-blur-sm"
                        />
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.95, y: 20 }}
                            className="relative w-full max-w-md bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-10 flex flex-col gap-8 shadow-2xl"
                        >
                            <div className="flex flex-col gap-2">
                                <h2 className="text-2xl font-bold uppercase tracking-tight text-black dark:text-white flex items-center gap-3">
                                    <Lock size={24} /> Change Password
                                </h2>
                                <p className="text-zinc-500 text-sm font-medium">Update your account security credentials.</p>
                            </div>

                            {passwordError && (
                                <div className="flex items-center gap-3 p-4 bg-red-500/10 border border-red-500 rounded-xl text-red-500 text-sm font-bold">
                                    <AlertCircle size={16} />
                                    {passwordError}
                                </div>
                            )}

                            <div className="flex flex-col gap-6">
                                <div className="flex flex-col gap-3">
                                    <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4">Current Password</label>
                                    <div className="relative">
                                        <Input
                                            type={showPasswords.current ? "text" : "password"}
                                            placeholder="••••••••"
                                            value={passwordForm.current}
                                            onChange={(e) => setPasswordForm({ ...passwordForm, current: e.target.value })}
                                        />
                                        <button
                                            onClick={() => setShowPasswords({ ...showPasswords, current: !showPasswords.current })}
                                            className="absolute right-4 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-black dark:hover:text-white"
                                        >
                                            {showPasswords.current ? <EyeOff size={16} /> : <Eye size={16} />}
                                        </button>
                                    </div>
                                </div>

                                <div className="flex flex-col gap-3">
                                    <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4">New Password</label>
                                    <div className="relative">
                                        <Input
                                            type={showPasswords.new ? "text" : "password"}
                                            placeholder="••••••••"
                                            value={passwordForm.new}
                                            onChange={(e) => setPasswordForm({ ...passwordForm, new: e.target.value })}
                                        />
                                        <button
                                            onClick={() => setShowPasswords({ ...showPasswords, new: !showPasswords.new })}
                                            className="absolute right-4 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-black dark:hover:text-white"
                                        >
                                            {showPasswords.new ? <EyeOff size={16} /> : <Eye size={16} />}
                                        </button>
                                    </div>
                                </div>

                                <div className="flex flex-col gap-3">
                                    <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4">Confirm Password</label>
                                    <div className="relative">
                                        <Input
                                            type={showPasswords.confirm ? "text" : "password"}
                                            placeholder="••••••••"
                                            value={passwordForm.confirm}
                                            onChange={(e) => setPasswordForm({ ...passwordForm, confirm: e.target.value })}
                                        />
                                        <button
                                            onClick={() => setShowPasswords({ ...showPasswords, confirm: !showPasswords.confirm })}
                                            className="absolute right-4 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-black dark:hover:text-white"
                                        >
                                            {showPasswords.confirm ? <EyeOff size={16} /> : <Eye size={16} />}
                                        </button>
                                    </div>
                                </div>
                            </div>

                            <div className="flex gap-4">
                                <Button variant="outline" className="flex-1" onClick={() => { setShowPasswordModal(false); setPasswordError(''); }}>Cancel</Button>
                                <Button className="flex-1" onClick={handlePasswordSubmit} disabled={loading}>
                                    {loading ? 'Updating...' : 'Update Password'}
                                </Button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>

            {/* Logout Confirmation Modal */}
            <AnimatePresence>
                {showLogoutConfirm && (
                    <div className="fixed inset-0 z-[60] flex items-center justify-center p-6">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            onClick={() => setShowLogoutConfirm(false)}
                            className="absolute inset-0 bg-black/50 dark:bg-black/80 backdrop-blur-sm"
                        />
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.95, y: 20 }}
                            className="relative w-full max-w-md bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-10 flex flex-col gap-8 shadow-2xl"
                        >
                            <div className="flex flex-col gap-2">
                                <h2 className="text-2xl font-bold uppercase tracking-tight text-black dark:text-white">Confirm Logout</h2>
                                <p className="text-zinc-500 text-sm font-medium">You will be logged out from all active sessions.</p>
                            </div>

                            <div className="flex gap-4">
                                <Button variant="outline" className="flex-1" onClick={() => setShowLogoutConfirm(false)} disabled={loading}>Cancel</Button>
                                <Button className="flex-1 bg-red-500 hover:bg-red-600" onClick={handleLogout} disabled={loading}>
                                    {loading ? 'Logging out...' : 'Logout'}
                                </Button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </div>
    );
}
