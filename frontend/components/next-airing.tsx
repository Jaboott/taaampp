'use client'

import {CalendarClock} from 'lucide-react'
import {useEffect, useState, useMemo} from 'react'

function formatCountDown(seconds: number): string {
    if (seconds <= 0) return 'Airing now'
    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = Math.floor(seconds % 60)
    const parts: string[] = []
    if (days) parts.push(`${days}d`)
    if (hours || days) parts.push(`${hours}h`)
    parts.push(`${minutes}m`)
    if (!days) parts.push(`${secs}s`)
    return parts.join(' ')
}

export function NextAiring({airingAt, episode}: { airingAt: number, episode: number }) {
    const [mounted, setMounted] = useState(false)
    const [remaining, setRemaining] = useState(0)

    useEffect(() => {
        setMounted(true)
        const tick = () =>
            setRemaining(Math.max(0, airingAt - Math.floor(Date.now() / 1000)))
        tick()
        const id = setInterval(tick, 1000)
        return () => clearInterval(id)
    }, [airingAt])

    const airDate = useMemo(() => {
        return new Date(airingAt * 1000).toLocaleString('en-US', {
            weekday: 'short',
            month: 'short',
            day: 'numeric',
            hour: 'numeric',
            minute: '2-digit',
        })
    }, [airingAt])

    return (
        <div className="flex items-center gap-3 rounded-xl border border-primary/25 bg-primary/10 px-4 py-3">
        <span className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary/20 text-primary">
        <CalendarClock className="size-5" aria-hidden="true"/>
        </span>
            <div className="min-w-0">
                <p className="text-xs font-medium uppercase tracking-wide text-primary">
                    {remaining > 0 ? `Episode ${episode} airing in` : `Episode ${episode}`}
                </p>
                <p className="truncate font-display text-lg font-semibold tabular-nums text-foreground">
                    {mounted ? formatCountDown(remaining) : 'Loading...'}
                </p>
                <p className="truncate text-xs text-muted-foreground">
                    {mounted ? airDate : '\u00A0'}
                </p>
            </div>
        </div>
    )
}
