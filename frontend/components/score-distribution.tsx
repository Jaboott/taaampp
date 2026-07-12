import type {MediaDetails} from '@/lib/query'
import {formatNumber} from '@/lib/format'

export function ScoreDistribution({media}: { media: MediaDetails }) {
    const scoreAmounts = new Map(
        media.score_distribution?.map(({score, amount}) => [score, amount]) ?? [],
    )

    const dist = Array.from({length: 10}, (_, index) => {
        const score = (index + 1) * 10

        return {
            score,
            amount: scoreAmounts.get(score) ?? 0,
        }
    })

    const total = dist.reduce((sum, {amount}) => sum + amount, 0)

    if (!total) return null

    const max = Math.max(...dist.map(({amount}) => amount), 1)

    return (
        <section
            aria-labelledby="score-dist-heading"
            className="rounded-2xl border border-border bg-card p-5 md:p-6"
        >
            <div className="mb-6 flex items-baseline justify-between gap-4">
                <h2
                    id="score-dist-heading"
                    className="font-display text-sm font-semibold uppercase tracking-widest text-muted-foreground"
                >
                    Score Distribution
                </h2>

                <span className="text-xs text-muted-foreground">
                    {formatNumber(total)} ratings
                </span>
            </div>

            <div className="flex h-52 items-end justify-between gap-2 md:gap-3">
                {dist.map(({score, amount}) => {
                    const heightPct = (amount / max) * 100
                    const percentage = Math.round((amount / total) * 100)

                    const normalizedScore = (score - 10) / 90
                    const hue = normalizedScore * 120
                    const barColor = `hsl(${hue}, 75%, 50%)`

                    return (
                        <div
                            key={score}
                            className="flex h-full flex-1 flex-col items-center justify-end gap-2"
                        >
                            <div
                                className="group relative w-full rounded-md transition hover:brightness-110"
                                style={{
                                    height: `max(${heightPct.toFixed(2)}%, 4px)`,
                                    backgroundColor: barColor,
                                }}
                            >
                                <span className="pointer-events-none absolute -top-6 left-1/2 -translate-x-1/2 whitespace-nowrap rounded bg-popover px-2 py-1 text-[10px] font-medium text-popover-foreground opacity-0 shadow transition-opacity group-hover:opacity-100">
                                    {formatNumber(amount)} · {percentage}%
                                </span>
                            </div>

                            <span className="text-xs font-medium tabular-nums text-muted-foreground">
                                {score}
                            </span>
                        </div>
                    )
                })}
            </div>
        </section>
    )
}