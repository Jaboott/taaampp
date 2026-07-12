'use client'

import {Play} from 'lucide-react'
import {useState} from 'react'
import type {MediaDetails} from '@/lib/query'

type VideoData = {
    thumbnail: string
    embed: string
}

export function getVideoData(videoUrl: string | null | undefined,): VideoData | null {
    if (!videoUrl) return null

    try {
        const normalizedUrl = /^https?:\/\//i.test(videoUrl) ? videoUrl : `https://${videoUrl}`

        const url = new URL(normalizedUrl)
        const id = url.searchParams.get('v')

        if (!id) return null

        if (url.hostname.includes('youtube.com')) {
            return {
                thumbnail: `https://img.youtube.com/vi/${id}/maxresdefault.jpg`,
                embed: `https://www.youtube.com/embed/${id}?autoplay=1&rel=0`,
            }
        }

        if (url.hostname.includes('dailymotion.com')) {
            return {
                thumbnail: `https://www.dailymotion.com/thumbnail/video/${id}`,
                embed: `https://www.dailymotion.com/embed/video/${id}?autoplay=1`,
            }
        }

        return null
    } catch {
        return null
    }
}

export function Trailer({media}: { media: MediaDetails }) {
    const [playing, setPlaying] = useState(false)
    const trailer = media.trailer
    if (!trailer) return null
    const videoData = getVideoData(trailer)
    if (!videoData) return null

    return (
        <section
            aria-labelledby="trailer-heading"
            className="rounded-2xl border border-border bg-card p-5 md:p-6"
        >
            <h2
                id="trailer-heading"
                className="mb-4 font-display text-sm font-semibold uppercase tracking-widest text-muted-foreground"
            >
                Trailer
            </h2>
            <div className="relative aspect-video w-full overflow-hidden rounded-xl bg-secondary">
                {playing ? (
                    <iframe
                        src={videoData.embed}
                        title="Anime trailer"
                        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                        allowFullScreen
                        className="size-full"
                    />
                ) : (
                    <button
                        type="button"
                        onClick={() => setPlaying(true)}
                        className="group relative size-full"
                        aria-label="Play trailer"
                    >
                        <img src={videoData.thumbnail || '/placeholder.svg'} alt="" className="size-full object-cover"/>
                        <span
                            className="absolute inset-0 bg-background/30 transition-colors group-hover:bg-background/10"/>
                        <span
                            className="absolute left-1/2 top-1/2 flex size-16 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-lg transition-transform group-hover:scale-110">
                            <Play className="ml-1 size-7 fill-current" aria-hidden="true"/>
                        </span>
                    </button>
                )}
            </div>
        </section>
    )
}