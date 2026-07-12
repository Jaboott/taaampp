'use client'

import {useState} from 'react'
import {stripHtml} from '@/lib/format'

export function Synopsis({description}: { description: string | null }) {
    const [expanded, setExpanded] = useState(false)
    const text = stripHtml(description)

    if (!text) return null

    const isLong = text.length > 420

    return (
        <section
            aria-labelledby="synopsis-heading"
            className="rounded-2xl border border-border bg-card p-5 md:p-6"
        >
            <h2
                id="synopsis-heading"
                className="mb-4 font-display text-sm font-semibold uppercase tracking-widest text-muted-foreground"
            >
                Synopsis
            </h2>
            <p
                className={`whitespace-pre-line text-pretty leading-relaxed text-foreground/90 ${
                    isLong && !expanded ? 'line-clamp-6' : ''
                }`}
            >
                {text}
            </p>
            {isLong ? (
                <button
                    type="button"
                    onClick={() => setExpanded((v) => !v)}
                    className="mt-3 text-sm font-medium text-primary transition-colors hover:text-primary/80"
                >
                    {expanded ? 'Show less' : 'Read more'}
                </button>
            ) : null}
        </section>
    )
}
