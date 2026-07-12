import {Heart, Star, TrendingUp} from 'lucide-react'
import type {ReactNode} from 'react'
import type {MediaDetails} from '@/lib/query'
import {formatNumber, formatType} from '@/lib/format'
import {NextAiring} from './next-airing'
import Image from 'next/image'

export function AnimeHero({media}: { media: MediaDetails }) {
    const title = media.titles.english || media.titles.romaji || media.titles.native || 'Untitled'
    const cover = media.cover_image || ''

    return (
        <section aria-label="Anime overview" className="relative">
            {/* Banner */}
            <div className="relative h-[42vh] min-h-64 w-full overflow-hidden md:h-[52vh]">
                {media.banner_image ? (
                    <Image
                        src={media.banner_image}
                        alt={`${title} banner art`}
                        fill
                        loading="eager"
                        sizes="100vw"
                        className="object-cover"
                    />
                ) : (
                    <div className="size-full bg-secondary"/>
                )}
                <div className="absolute inset-0 bg-gradient-to-t from-background via-background/85 to-background/20"/>
            </div>

            {/* Overlapping content */}
            <div
                className="relative z-10 mx-auto -mt-24 flex max-w-6xl flex-col gap-6 px-4 md:-mt-28 md:flex-row md:items-end md:px-6">
                <div
                    className="relative w-40 shrink-0 self-center overflow-hidden rounded-xl border border-border shadow-2xl md:w-56 md:self-auto">
                    <Image
                        src={cover || '/placeholder.svg'}
                        alt={`${title} cover art`}
                        loading="eager"
                        width={400}
                        height={560}
                        className="aspect-[2/3] w-full object-cover"
                    />
                </div>

                <div className="flex min-w-0 flex-1 flex-col gap-4 pb-2 text-center md:text-left">
                    <div className="space-y-2">
                        <div className="flex flex-wrap items-center justify-center gap-2 md:justify-start">
                            <span
                                className="rounded-full bg-primary/15 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-primary">
                                {formatType(media.format)}
                            </span>
                            {media.season_year && (
                                <span
                                    className="rounded-full bg-secondary px-3 py-1 text-xs font-medium text-secondary-foreground">
                                    {media.season ? `${formatType(media.season)} ` : ''}
                                    {media.season_year}
                                </span>
                            )}
                        </div>
                        <h1 className="text-balance font-display text-3xl font-bold leading-tight text-foreground md:text-5xl">
                            {title}
                        </h1>
                        {media.titles.native && media.titles.native !== title && (
                            <p className="text-pretty text-sm text-muted-foreground">{media.titles.native}</p>
                        )}
                    </div>

                    {/* Genre Display*/}
                    <div className="flex flex-wrap items-center justify-center gap-2 md:justify-start">
                        {media.genres.map((g) => (
                            <span
                                key={g}
                                className="rounded-full border border-border bg-card/60 px-3 py-1 text-xs text-muted-foreground"
                            >
                                {g}
                            </span>
                        ))}
                    </div>

                    {/* Popularity Stats Display*/}
                    <div className="flex flex-wrap items-center justify-center gap-4 md:justify-start">
                        <HeroStat
                            icon={<Star className="size-4 text-accent" aria-hidden="true"/>}
                            label="Average"
                            value={media.average_score ? `${media.average_score}%` : '—'}
                        />
                        <HeroStat
                            icon={<TrendingUp className="size-4 text-primary" aria-hidden="true"/>}
                            label="Popularity"
                            value={formatNumber(media.popularity)}
                        />
                        <HeroStat
                            icon={<Heart className="size-4 text-destructive" aria-hidden="true"/>}
                            label="Favourites"
                            value={formatNumber(media.favourites)}
                        />
                    </div>
                </div>

                {media.airing_schedule ? (
                    <div className="w-full md:w-72 md:pb-2">
                        <NextAiring
                            airingAt={media.airing_schedule.airing_at}
                            episode={media.airing_schedule.episode}
                        />
                    </div>
                ) : null}
            </div>
        </section>
    )
}

function HeroStat({ icon, label, value }: { icon: ReactNode, label: string, value: string }) {
    return (
        <div className="flex items-center gap-2">
            {icon}
            <div className="text-left leading-tight">
                <p className="font-display text-base font-semibold tabular-nums text-foreground">{value}</p>
                <p className="text-[11px] uppercase tracking-wide text-muted-foreground">{label}</p>
            </div>
        </div>
    )
}
