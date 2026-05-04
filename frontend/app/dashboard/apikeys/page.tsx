'use client';

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "motion/react";
import { Plus, Copy, Check, Trash2, AlertCircle, RefreshCw, Eye, EyeOff } from "lucide-react";
import { Button, Input } from "@/lib/components/ui";
import { ApiKey } from "@/lib/types";
import { fetchJSON } from "@/lib/fetcher";
import { apiFetch, apiKeyAPI } from "@/lib/api";
import { useAuth } from "@/lib/AuthContext";

export default function ApiKeysPage() {
    const { token } = useAuth();
    const [keys, setKeys] = useState<ApiKey[]>([]);
    const [loading, setLoading] = useState(true);
    const [newKeyData, setNewKeyData] = useState<{ publicKey: string; secret: string; name: string } | null>(null);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [name, setName] = useState("");
    const [copied, setCopied] = useState<string | null>(null);
    const [hiddenKeys, setHiddenKeys] = useState<Set<string>>(new Set());
    const [regeneratingId, setRegeneratingId] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [creatingKey, setCreatingKey] = useState(false);

    // Load API keys from backend
    useEffect(() => {
        if (token) {
            loadKeys();
        }
    }, [token]);

    const loadKeys = async () => {
        setLoading(true);
        setError(null);
        try {
            const res = await apiKeyAPI.list(token!);
            console.log("API Keys Response:", res);
            if (Array.isArray(res)) {
                const formattedKeys = res.map((key: any) => ({
                    id: key.id,
                    name: key.name,
                    publicKey: key.public_key,
                    isHidden: true,
                    createdAt: key.created_at,
                    lastUsed: null
                }));
                setKeys(formattedKeys);
                // Hide all secret keys initially
                setHiddenKeys(new Set(formattedKeys.map((k: ApiKey) => k.id)));
            }
        } catch (err: any) {
            console.error("Failed to load API keys:", err);
            setError(`Failed to load API keys: ${err?.response?.data?.message || err?.message || 'Unknown error'}`);
        }
        setLoading(false);
    };

    // Generate a random key
    const generateKey = (): string => {
        const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        let result = '';
        for (let i = 0; i < 32; i++) {
            result += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        return result;
    };

    const handleCreate = async () => {
        if (!name.trim() || !token) {
            setError("Please enter a key name");
            return;
        }

        setCreatingKey(true);
        setError(null);
        try {
            console.log("Creating API key with name:", name);
            const res = await apiKeyAPI.create(token, { name });
            console.log("API Key Creation Response:", res);

            const newKey: ApiKey = {
                id: res.id,
                name: res.name,
                publicKey: res.public_key,
                isHidden: true,
                createdAt: res.created_at,
                lastUsed: null
            };

            // Show the new key modal with secret
            setNewKeyData({
                publicKey: res.public_key,
                secret: res.secret_key,
                name: res.name
            });

            // Add to keys list
            setKeys(prev => [newKey, ...prev]);
            setHiddenKeys(prev => new Set([...prev, newKey.id]));

            setName("");
            setShowCreateModal(false);
        } catch (err: any) {
            console.error("Failed to create API key:", err);
            setError(`Failed to create API key: ${err?.response?.data?.message || err?.message || 'Unknown error'}`);
        } finally {
            setCreatingKey(false);
        }
    };

    const handleDelete = (id: string) => {
        if (!token) return;

        apiKeyAPI.delete(token, id).then(() => {
            setKeys(prev => prev.filter(k => k.id !== id));
            setHiddenKeys(prev => {
                const newSet = new Set(prev);
                newSet.delete(id);
                return newSet;
            });
        }).catch(err => {
            console.error("Failed to delete API key:", err);
        });
    };

    const handleRegenerate = (id: string) => {
        if (!token) return;

        setRegeneratingId(id);
        apiKeyAPI.regenerate(token, id).then((res) => {
            // Show regenerate modal with new secret
            const key = keys.find(k => k.id === id);
            if (key) {
                setNewKeyData({
                    publicKey: key.publicKey,
                    secret: res.secret_key,
                    name: key.name
                });
            }
            setRegeneratingId(null);
        }).catch(err => {
            console.error("Failed to regenerate API key:", err);
            setRegeneratingId(null);
        });
    };

    const toggleVisibility = (id: string) => {
        setHiddenKeys(prev => {
            const newSet = new Set(prev);
            if (newSet.has(id)) {
                newSet.delete(id);
            } else {
                newSet.add(id);
            }
            return newSet;
        });
    };

    const copyToClipboard = (text: string, type: string) => {
        navigator.clipboard.writeText(text);
        setCopied(type);
        setTimeout(() => setCopied(null), 2000);
    };

    return (
        <div className="flex flex-col gap-10">
            <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
                <div className="flex flex-col gap-1">
                    <h1 className="text-4xl font-bold tracking-tighter uppercase">Developer Platform</h1>
                    <p className="text-zinc-500 font-medium tracking-tight">Manage your infrastructure access keys.</p>
                </div>
                <Button onClick={() => setShowCreateModal(true)} size="lg" className="h-12 px-6">
                    <Plus size={18} className="mr-2" /> Generate Key
                </Button>
            </header>

            {/* Error Alert */}
            <AnimatePresence>
                {error && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="bg-red-500/10 border border-red-500/20 rounded-2xl p-6 flex items-center gap-3 text-red-500"
                    >
                        <AlertCircle size={20} />
                        <span className="text-sm font-medium">{error}</span>
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Secrets Warning (Conditional) */}
            <AnimatePresence>
                {newKeyData && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="bg-amber-500/10 border border-amber-500/20 rounded-2xl p-6 flex flex-col gap-6 overflow-hidden"
                    >
                        <div className="flex items-center gap-3 text-amber-500">
                            <AlertCircle size={20} />
                            <span className="text-sm font-bold uppercase tracking-wider">Secret Key Generated</span>
                        </div>
                        <p className="text-zinc-500 dark:text-zinc-400 text-sm leading-relaxed max-w-2xl">
                            This secret key will be shown <span className="text-black dark:text-white font-bold underline">only once</span>.
                            If you lose it, you must regenerate the key. Vanguard does not store plain-text secrets.
                        </p>
                        <div className="flex flex-col md:flex-row gap-4">
                            <div className="flex-1 bg-zinc-50 dark:bg-black rounded-xl border border-zinc-200 dark:border-zinc-800 p-4 flex items-center justify-between group">
                                <div className="flex flex-col gap-0.5">
                                    <span className="text-[10px] uppercase font-bold text-zinc-500 dark:text-zinc-600">Public Key</span>
                                    <code className="text-zinc-600 dark:text-zinc-300 font-mono text-xs">{newKeyData.publicKey}</code>
                                </div>
                                <button onClick={() => copyToClipboard(newKeyData.publicKey, 'public')} className="text-zinc-400 hover:text-black dark:hover:text-white transition-colors">
                                    {copied === 'public' ? <Check size={16} /> : <Copy size={16} />}
                                </button>
                            </div>
                            <div className="flex-1 bg-zinc-50 dark:bg-black rounded-xl border border-zinc-200 dark:border-zinc-800 p-4 flex items-center justify-between group">
                                <div className="flex flex-col gap-0.5">
                                    <span className="text-[10px] uppercase font-bold text-zinc-500 dark:text-zinc-600">Secret Key</span>
                                    <code className="text-amber-600 dark:text-amber-500 font-mono text-xs">{newKeyData.secret}</code>
                                </div>
                                <button onClick={() => copyToClipboard(newKeyData.secret, 'secret')} className="text-zinc-400 hover:text-black dark:hover:text-white transition-colors">
                                    {copied === 'secret' ? <Check size={16} /> : <Copy size={16} />}
                                </button>
                            </div>
                        </div>
                        <Button variant="secondary" className="self-end" onClick={() => setNewKeyData(null)}>
                            I've stored it safely
                        </Button>
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Keys List */}
            <div className="rounded-3xl border border-zinc-200 dark:border-zinc-900 overflow-hidden bg-zinc-50 dark:bg-zinc-900/10">
                <div className="grid grid-cols-4 p-6 border-b border-zinc-200 dark:border-zinc-900 text-[10px] uppercase font-bold tracking-widest text-zinc-500">
                    <span>Key Name</span>
                    <span>Preview</span>
                    <span>Status</span>
                    <span className="text-right">Actions</span>
                </div>
                <div className="flex flex-col">
                    {keys.length === 0 ? (
                        <div className="p-12 text-center text-zinc-500 font-medium tracking-tight">
                            No active API keys found.
                        </div>
                    ) : (
                        keys.map((key) => (
                            <div key={key.id} className="grid grid-cols-4 p-6 border-b border-zinc-100 dark:border-zinc-900 last:border-0 items-center hover:bg-zinc-100/50 dark:hover:bg-zinc-900/30 transition-colors group">
                                <div className="flex flex-col gap-1">
                                    <span className="font-bold tracking-tight text-black dark:text-white">{key.name}</span>
                                    <span className="text-[10px] uppercase font-bold text-zinc-400 dark:text-zinc-600 font-mono">{new Date(key.createdAt).toLocaleDateString()}</span>
                                </div>
                                <code className="text-xs font-mono text-zinc-400 opacity-50 group-hover:opacity-100 transition-opacity">
                                    {key.publicKey.slice(0, 12)}...
                                </code>
                                <div>
                                    <span className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-500 text-[10px] uppercase font-bold tracking-widest">
                                        <span className="w-1 h-1 rounded-full bg-emerald-500" />
                                        Active
                                    </span>
                                </div>
                                <div className="flex justify-end gap-2">
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        onClick={() => toggleVisibility(key.id)}
                                        className="text-zinc-400 hover:text-black dark:hover:text-white"
                                    >
                                        {hiddenKeys.has(key.id) ? <EyeOff size={14} /> : <Eye size={14} />}
                                    </Button>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        onClick={() => handleRegenerate(key.id)}
                                        disabled={regeneratingId === key.id}
                                        className="text-zinc-400 hover:text-amber-500"
                                    >
                                        <RefreshCw size={14} />
                                    </Button>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="text-zinc-400 hover:text-red-500"
                                        onClick={() => handleDelete(key.id)}
                                    >
                                        <Trash2 size={14} />
                                    </Button>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>

            {/* Create Modal */}
            <AnimatePresence>
                {showCreateModal && (
                    <div className="fixed inset-0 z-[60] flex items-center justify-center p-6">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            onClick={() => setShowCreateModal(false)}
                            className="absolute inset-0 bg-black/50 dark:bg-black/80 backdrop-blur-sm"
                        />
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.95, y: 20 }}
                            className="relative w-full max-w-md bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-10 flex flex-col gap-8 shadow-2xl"
                        >
                            <div className="flex flex-col gap-1">
                                <h2 className="text-2xl font-bold uppercase tracking-tight text-black dark:text-white">Generate Service Key</h2>
                                <p className="text-zinc-500 text-sm font-medium">Assign a label to this key for internal tracking.</p>
                            </div>
                            <div className="flex flex-col gap-3">
                                <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 ml-4">Key Identity Label</label>
                                <Input
                                    autoFocus
                                    placeholder="e.g. Production Mobile App"
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                />
                            </div>
                            <div className="flex gap-4">
                                <Button variant="outline" className="flex-1" onClick={() => setShowCreateModal(false)} disabled={creatingKey}>Cancel</Button>
                                <Button className="flex-1" onClick={handleCreate} disabled={creatingKey || !name.trim()}>
                                    {creatingKey ? "Creating..." : "Generate"}
                                </Button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </div>
    );
}
