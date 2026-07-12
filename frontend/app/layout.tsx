import {ReloadOnHistoryNavigation} from "@/components/reload-on-history-navigation";
import {ReactNode} from "react";

import {Analytics} from '@vercel/analytics/next'
import type {Metadata, Viewport} from 'next'
import {Geist, Space_Grotesk} from 'next/font/google'
import './globals.css'

const geistSans = Geist({
    subsets: ['latin'],
    variable: '--font-geist-sans',
})

const spaceGrotesk = Space_Grotesk({
    subsets: ['latin'],
    variable: '--font-space-grotesk',
})

export const metadata: Metadata = {
    title: 'taaampp Media Summary',
}

export const viewport: Viewport = {
    colorScheme: 'dark',
    themeColor: '#0d0f16',
}

export default function RootLayout({children,}: Readonly<{ children: ReactNode }>) {
    return (
        <html lang="en" className={`${geistSans.variable} ${spaceGrotesk.variable} bg-background`}>
        <body className="font-display antialiased">
        <ReloadOnHistoryNavigation/>
        {children}
        {process.env.NODE_ENV === 'production' && <Analytics/>}
        </body>
        </html>
    )
}
