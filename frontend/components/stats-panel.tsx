import type {MediaDetails} from '@/lib/query'
import {
    formatDate,
    formatDuration,
    formatEnum,
    formatNumber,
    formatStatus,
    formatType,
} from '@/lib/format'

export function StatsPanel({media}: { media: MediaDetails }) {
    const studios = media.studios? media.studios
        .slice(0, 3)
        .join(', ') : null

    const rows: { label: string; value: string }[] = [
        {label: 'Format', value: formatType(media.format)},
        {label: 'Episodes', value: media.episodes ? String(media.episodes) : '—'},
        {label: 'Episode Duration', value: formatDuration(media.duration)},
        {label: 'Status', value: formatStatus(media.status)},
        {label: 'Start Date', value: formatDate(media.start_date)},
        {
            label: 'Season',
            value: media.season
                ? `${formatEnum(media.season)}${media.season_year ? ` ${media.season_year}` : ''}`
                : media.season_year
                    ? String(media.season_year)
                    : '—',
        },
        {label: 'Average Score', value: media.average_score ? `${media.average_score}%` : '—'},
        {label: 'Popularity', value: formatNumber(media.popularity)},
        {label: 'Favourites', value: formatNumber(media.favourites)},
        {label: 'Studios', value: studios || '—'},
        {label: 'Source', value: formatEnum(media.source)},
        {label: 'Genres', value: media.genres.length ? media.genres.join(', ') : '—'},
    ]

    return (
        <section
            aria-labelledby="stats-heading"
            className="rounded-2xl border border-border bg-card p-5 md:p-6"
        >
            <h2
                id="stats-heading"
                className="mb-4 font-display text-sm font-semibold uppercase tracking-widest text-muted-foreground"
            >
                Details
            </h2>
            <dl className="divide-y divide-border">
                {rows.map((row) => (
                    row.value != "—" &&
                    <div key={row.label} className="flex items-start justify-between gap-4 py-3">
                        <dt className="text-sm text-muted-foreground">{row.label}</dt>
                        <dd className="text-right text-sm font-medium text-foreground">{row.value}</dd>
                    </div>
                ))}
            </dl>
        </section>
    )
}
