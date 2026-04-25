'use client';

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "motion/react";
import {
    Wallet as WalletIcon,
    ArrowUpRight,
    ArrowDownLeft,
    Shield,
    Copy,
    Check,
    Search,
    TrendingUp,
    Clock,
    History,
    Download,
    Building
} from "lucide-react";
import { Button, Input } from "@/lib/components/ui";
import { Wallet as WalletType, Transaction } from "@/lib/types";
import { cn } from "@/lib/utils";
import { isAddress, formatEther, parseEther } from "ethers";

type FilterType = 'all' | 'deposit' | 'withdrawal';

export default function WalletPage() {
    const [wallet, setWallet] = useState<WalletType | null>(null);
    const [loading, setLoading] = useState(true);
    const [copied, setCopied] = useState(false);
    const [filter, setFilter] = useState<FilterType>('all');
    const [search, setSearch] = useState("");
    const [showModal, setShowModal] = useState<'deposit' | 'withdrawal' | null>(null);
    const [targetAddress, setTargetAddress] = useState("");
    const [addressError, setAddressError] = useState("");

    const [transactions] = useState<Transaction[]>([
        { id: 'TX-9021', walletId: 'w1', amount: '1.250', type: 'deposit', status: 'completed', createdAt: '2023-10-24T14:20:00Z', txHash: '0x34a...2e1' },
        { id: 'TX-9022', walletId: 'w1', amount: '0.400', type: 'withdrawal', status: 'pending', createdAt: '2023-10-23T09:12:00Z', txHash: '0x91b...f42' },
        { id: 'TX-9023', walletId: 'w1', amount: '2500.00', type: 'deposit', status: 'completed', createdAt: '2023-10-22T18:45:00Z', txHash: '0x45d...a3e' },
        { id: 'TX-9024', walletId: 'w1', amount: '0.150', type: 'withdrawal', status: 'failed', createdAt: '2023-10-21T11:30:00Z', txHash: '0xa2c...11d' },
    ]);

    useEffect(() => {
        initWallet();
    }, []);

    const initWallet = async () => {
        setLoading(true);
        try {
            const res = await fetch('/api/wallet/balance');
            const data = await res.json();

            const fullWallet = {
                id: 'w1',
                userId: 'u1',
                publicKey: '0x71C7656EC7ab88b098defB751B7401B5f6d8976F',
                balance: data.balance || "12.45",
                createdAt: new Date().toISOString()
            };
            setWallet(fullWallet);
        } catch (e) {
            console.error(e);
            setWallet({
                id: 'w1',
                userId: 'u1',
                publicKey: '0x71C7656EC7ab88b098defB751B7401B5f6d8976F',
                balance: "12.45",
                createdAt: new Date().toISOString()
            });
        }
        setLoading(false);
    };

    const copyAddress = () => {
        if (wallet) {
            navigator.clipboard.writeText(wallet.publicKey);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        }
    };

    const handleAddressChange = (val: string) => {
        setTargetAddress(val);
        if (val && !isAddress(val)) {
            setAddressError("Invalid Ethereum address format via ethers.js validation");
        } else {
            setAddressError("");
        }
    };

    const filteredTransactions = transactions.filter(t => {
        if (filter !== 'all' && t.type !== filter) return false;
        const searchLower = search.toLowerCase();
        if (search && !t.id.toLowerCase().includes(searchLower) && !t.txHash?.toLowerCase().includes(searchLower)) return false;
        return true;
    });

    if (loading) return <div className="flex items-center justify-center h-[60vh] text-zinc-500 font-bold uppercase tracking-widest animate-pulse">Establishing Secure Connection...</div>;

    return (
        <div className="flex flex-col gap-10">
            <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
                <div className="flex flex-col gap-1">
                    <h1 className="text-4xl font-bold tracking-tighter uppercase">Vault Ledger</h1>
                    <p className="text-zinc-500 font-medium tracking-tight">Institutional assets managed via ethers.js infrastructure.</p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" size="lg" onClick={() => setShowModal('deposit')} className="h-12">
                        <ArrowDownLeft size={18} className="mr-2" /> Deposit
                    </Button>
                    <Button size="lg" onClick={() => setShowModal('withdrawal')} className="h-12">
                        <ArrowUpRight size={18} className="mr-2" /> Withdraw
                    </Button>
                </div>
            </header>

            {/* Overview Cards */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div className="md:col-span-2 p-8 rounded-[2.5rem] bg-gradient-to-br from-zinc-900 to-black border border-zinc-800 flex flex-col justify-between min-h-[220px] relative overflow-hidden group">
                    <div className="relative z-10 flex justify-between items-start">
                        <div className="flex flex-col gap-2">
                            <span className="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-500">Total Valuation</span>
                            <h2 className="text-5xl font-bold tracking-tighter text-white">{wallet?.balance} <span className="text-zinc-600">ETH</span></h2>
                        </div>
                        <div className="p-4 rounded-2xl bg-white/5 border border-white/10 backdrop-blur-md">
                            <WalletIcon size={24} className="text-white" />
                        </div>
                    </div>

                    <div className="relative z-10 flex flex-col gap-4">
                        <div className="flex items-center justify-between bg-white/5 rounded-xl border border-white/10 px-4 py-2">
                            <code className="text-zinc-400 font-mono text-xs break-all">{wallet?.publicKey}</code>
                            <button onClick={copyAddress} className="text-zinc-500 hover:text-white transition-colors ml-4">
                                {copied ? <Check size={16} /> : <Copy size={16} />}
                            </button>
                        </div>
                        <div className="flex gap-4">
                            <div className="flex items-center gap-1.5 text-[10px] font-bold text-emerald-500 uppercase tracking-widest">
                                <div className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                                Synchronized
                            </div>
                            <div className="flex items-center gap-1.5 text-[10px] font-bold text-zinc-600 uppercase tracking-widest">
                                <Shield size={12} />
                                AES-256
                            </div>
                        </div>
                    </div>
                    <div className="absolute top-0 right-0 w-48 h-48 bg-white/5 blur-[60px] rounded-full -translate-y-1/2 translate-x-1/2 group-hover:bg-white/10 transition-all duration-700" />
                </div>

                {[
                    { label: "30D Volume", value: "$842k", icon: TrendingUp, color: "text-white" },
                    { label: "Active Queues", value: "02", icon: Clock, color: "text-amber-500" },
                ].map((stat, i) => (
                    <div key={i} className="p-8 rounded-[2.5rem] bg-zinc-50 dark:bg-zinc-900/40 border border-zinc-200 dark:border-zinc-900 flex flex-col justify-between">
                        <div className="w-12 h-12 rounded-2xl bg-white dark:bg-black border border-zinc-200 dark:border-zinc-800 flex items-center justify-center">
                            <stat.icon size={20} className={stat.color} />
                        </div>
                        <div className="flex flex-col gap-1">
                            <span className="text-[10px] uppercase font-bold tracking-widest text-zinc-500">{stat.label}</span>
                            <span className={cn("text-3xl font-bold tracking-tight text-black dark:text-white")}>{stat.value}</span>
                        </div>
                    </div>
                ))}
            </div>

            {/* Transactions Section */}
            <div className="flex flex-col gap-6">
                <div className="flex flex-col md:flex-row justify-between items-center gap-4">
                    <div className="flex items-center gap-1 bg-zinc-100 dark:bg-zinc-900/50 p-1 rounded-full border border-zinc-200 dark:border-zinc-800 self-start">
                        {(['all', 'deposit', 'withdrawal'] as const).map((t) => (
                            <button
                                key={t}
                                onClick={() => setFilter(t)}
                                className={cn(
                                    "px-6 py-2 rounded-full text-xs font-bold uppercase tracking-widest transition-all",
                                    filter === t
                                        ? "bg-black text-white dark:bg-white dark:text-black shadow-lg"
                                        : "text-zinc-500 hover:text-black dark:hover:text-white"
                                )}
                            >
                                {t}
                            </button>
                        ))}
                    </div>

                    <div className="flex items-center gap-3 w-full md:w-auto">
                        <div className="relative flex-1 md:w-64">
                            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-400 dark:text-zinc-600" size={16} />
                            <Input
                                placeholder="Search TXIDs..."
                                className="pl-12 py-2.5 h-11"
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                            />
                        </div>
                        <Button variant="outline" size="md" className="h-11">
                            <Download size={16} />
                        </Button>
                    </div>
                </div>

                <div className="rounded-[2.5rem] border border-zinc-200 dark:border-zinc-900 overflow-hidden bg-zinc-50 dark:bg-zinc-900/10 shadow-xl">
                    <div className="grid grid-cols-5 p-6 bg-zinc-100/50 dark:bg-zinc-900/30 border-b border-zinc-200 dark:border-zinc-900 text-[10px] uppercase font-bold tracking-widest text-zinc-500">
                        <span className="col-span-2">Transaction Details</span>
                        <span>Amount</span>
                        <span>Status</span>
                        <span className="text-right">Timestamp</span>
                    </div>

                    <div className="flex flex-col">
                        {filteredTransactions.length === 0 ? (
                            <div className="p-20 text-center flex flex-col items-center gap-4">
                                <div className="w-16 h-16 rounded-full bg-zinc-100 dark:bg-zinc-900 flex items-center justify-center text-zinc-300 dark:text-zinc-700">
                                    <History size={32} />
                                </div>
                                <p className="text-zinc-500 font-medium tracking-tight uppercase text-xs font-bold tracking-[0.2em]">No records found</p>
                            </div>
                        ) : (
                            filteredTransactions.map((tx) => (
                                <div
                                    key={tx.id}
                                    className="grid grid-cols-5 p-6 border-b border-zinc-100 dark:border-zinc-900 last:border-0 items-center hover:bg-zinc-100/50 dark:hover:bg-zinc-900/30 transition-colors group"
                                >
                                    <div className="col-span-2 flex items-center gap-4">
                                        <div className={cn(
                                            "w-10 h-10 rounded-xl flex items-center justify-center",
                                            tx.type === 'deposit' ? "bg-emerald-500/10 text-emerald-500" : "bg-black/5 dark:bg-white/5 text-black dark:text-white"
                                        )}>
                                            {tx.type === 'deposit' ? <ArrowDownLeft size={20} /> : <ArrowUpRight size={20} />}
                                        </div>
                                        <div className="flex flex-col gap-0.5">
                                            <span className="font-mono text-sm text-zinc-600 dark:text-zinc-300 group-hover:text-black dark:group-hover:text-white transition-colors">{tx.id}</span>
                                            <div className="flex items-center gap-2">
                                                <span className="text-[10px] font-bold uppercase text-zinc-400 dark:text-zinc-600">{tx.type} ETH</span>
                                                <div className="w-1 h-1 rounded-full bg-zinc-300 dark:bg-zinc-800" />
                                                <span className="text-[10px] font-mono text-zinc-400">{tx.txHash}</span>
                                            </div>
                                        </div>
                                    </div>
                                    <div>
                                        <span className={cn(
                                            "text-sm font-bold tracking-tight",
                                            tx.type === 'deposit' ? "text-emerald-600 dark:text-emerald-500" : "text-black dark:text-white"
                                        )}>
                                            {tx.type === 'deposit' ? '+' : '-'}{tx.amount} ETH
                                        </span>
                                    </div>
                                    <div>
                                        <span className={cn(
                                            "inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-widest",
                                            tx.status === 'completed' && "bg-emerald-500/10 text-emerald-600 dark:text-emerald-500",
                                            tx.status === 'pending' && "bg-amber-500/10 text-amber-600 dark:text-amber-500",
                                            tx.status === 'failed' && "bg-red-500/10 text-red-600 dark:text-red-500",
                                        )}>
                                            {tx.status}
                                        </span>
                                    </div>
                                    <div className="text-right">
                                        <span className="text-[10px] font-bold uppercase text-zinc-400 dark:text-zinc-500">{new Date(tx.createdAt).toLocaleDateString()}</span>
                                    </div>
                                </div>
                            )
                            ))}
                    </div>
                </div>
            </div>

            {/* Transaction Modals */}
            <AnimatePresence>
                {showModal && (
                    <div className="fixed inset-0 z-[100] flex items-center justify-center p-6">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            onClick={() => setShowModal(null)}
                            className="absolute inset-0 bg-black/50 dark:bg-black/90 backdrop-blur-md"
                        />
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.95, y: 20 }}
                            className="relative w-full max-w-lg bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-[2.5rem] p-12 shadow-2xl flex flex-col gap-10"
                        >
                            <div className="flex items-center gap-4">
                                <div className={cn(
                                    "w-16 h-16 rounded-3xl flex items-center justify-center",
                                    showModal === 'deposit' ? "bg-emerald-500 text-white" : "bg-black dark:bg-white text-white dark:text-black"
                                )}>
                                    {showModal === 'deposit' ? <ArrowDownLeft size={32} /> : <ArrowUpRight size={32} />}
                                </div>
                                <div className="flex flex-col">
                                    <h2 className="text-3xl font-bold uppercase tracking-tighter text-black dark:text-white">{showModal} Assets</h2>
                                    <p className="text-zinc-500 text-sm font-medium tracking-tight text-black dark:text-zinc-500">Access Vanguard's global liquidity network.</p>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div className="p-6 rounded-3xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/50 transition-all flex flex-col gap-4 text-left group">
                                    <WalletIcon size={24} className="text-zinc-400 dark:text-zinc-500" />
                                    <div className="flex flex-col">
                                        <span className="font-bold text-sm uppercase tracking-tight text-black dark:text-white">Web3 Node</span>
                                        <span className="text-[10px] font-bold text-zinc-400 dark:text-zinc-600 uppercase tracking-widest">Connect MetaMask</span>
                                    </div>
                                </div>
                                <div className="p-6 rounded-3xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/50 transition-all flex flex-col gap-4 text-left group">
                                    <Building className="text-zinc-400 dark:text-zinc-500" size={24} />
                                    <div className="flex flex-col">
                                        <span className="font-bold text-sm uppercase tracking-tight text-black dark:text-white">FIAT Node</span>
                                        <span className="text-[10px] font-bold text-zinc-400 dark:text-zinc-600 uppercase tracking-widest">Institutional Swift</span>
                                    </div>
                                </div>
                            </div>

                            <div className="flex flex-col gap-6">
                                <div className="flex flex-col gap-3">
                                    <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 dark:text-zinc-600 ml-4">
                                        {showModal === 'deposit' ? 'Source' : 'Destination'} Wallet Address
                                    </label>
                                    <Input
                                        placeholder="0x..."
                                        className={cn(
                                            "h-14 px-6 font-mono text-sm",
                                            addressError ? "border-red-500" : ""
                                        )}
                                        value={targetAddress}
                                        onChange={(e) => handleAddressChange(e.target.value)}
                                    />
                                    {addressError && (
                                        <span className="text-[10px] text-red-500 font-bold ml-4 uppercase">{addressError}</span>
                                    )}
                                </div>

                                <div className="flex flex-col gap-3">
                                    <label className="text-[10px] uppercase tracking-widest font-bold text-zinc-500 dark:text-zinc-600 ml-4">Input ETH Amount</label>
                                    <Input placeholder="0.00" className="text-2xl h-16 px-6 font-mono" />
                                </div>
                            </div>

                            <Button size="lg" className="h-14 font-bold text-base" onClick={() => setShowModal(null)} disabled={!!addressError}>
                                Confirm Intent
                            </Button>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </div>
    );
}
