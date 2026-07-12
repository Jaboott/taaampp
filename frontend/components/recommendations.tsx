import Link from 'next/link'
import Image from 'next/image'
import {Star} from 'lucide-react'
import type {Media} from '@/lib/query'
import {formatType} from '@/lib/format'

export function Recommendations({items}: { items: Media[] }) {
    if (!items.length) return null

    return (
        <section aria-labelledby="recs-heading">
            <h2
                id="recs-heading"
                className="mb-4 font-display text-sm font-semibold uppercase tracking-widest text-muted-foreground"
            >
                Recommendations
            </h2>
            <div className="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-6">
                {items.map((rec) => {
                    const title = rec.titles.english || rec.titles.romaji || 'Untitled'
                    return (
                        <Link
                            key={rec.id}
                            href={`/anime/${rec.id}`}
                            className="group flex flex-col gap-2"
                            title={title}
                        >
                            <div className="relative overflow-hidden rounded-lg border border-border bg-secondary">
                                <Image
                                    src={rec.cover_image || '/placeholder.svg'}
                                    alt={`${title} cover`}
                                    width={220}
                                    height={330}
                                    sizes="(max-width: 640px) 33vw, (max-width: 768px) 25vw, 16vw"
                                    className="aspect-[2/3] w-full object-cover transition-transform duration-300 group-hover:scale-105"
                                />
                                {rec.average_score ? (
                                    <span className="absolute left-1.5 top-1.5 flex items-center gap-1 rounded-md bg-background/85 px-1.5 py-0.5 text-[10px] font-semibold text-foreground backdrop-blur">
                                        <Star className="size-3 fill-accent text-accent" aria-hidden="true"/>
                                        {rec.average_score}
                                    </span>
                                ) : null}
                            </div>
                            <div className="min-w-0">
                                <p className="line-clamp-2 text-xs font-medium leading-tight text-foreground transition-colors group-hover:text-primary">
                                    {title}
                                </p>
                                <p className="mt-0.5 text-[11px] text-muted-foreground">
                                    {formatType(rec.format)}
                                    {rec.season_year ? ` · ${rec.season_year}` : ''}
                                </p>
                            </div>
                        </Link>
                    )
                })}
            </div>
        </section>
    )
}
