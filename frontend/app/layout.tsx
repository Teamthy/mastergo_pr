import React from 'react';
import type { Metadata } from 'next';
import '../styles/globals.css'
import { Providers } from './providers';
import { Navbar } from '../components/Navbar';

export const metadata: Metadata = {
  title: 'MY GO + NEXTJS + ETHERJS PR',
  description: 'Secure, high-performance institutional-grade Ethereum infrastructure.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="bg-white dark:bg-black text-black dark:text-white transition-colors duration-300">
        <Providers>
          <Navbar />
          <main className="pt-24 pb-20 max-w-7xl mx-auto px-6 min-h-screen">
            {children}
          </main>
        </Providers>
      </body>
    </html>
  );
}
