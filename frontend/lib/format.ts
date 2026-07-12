import type {FuzzyDate} from './query'

const FORMAT_LABELS: Record<string, string> = {
    TV: 'TV',
    TV_SHORT: 'TV Short',
    MOVIE: 'Movie',
    SPECIAL: 'Special',
    OVA: 'OVA',
    ONA: 'ONA',
    MUSIC: 'Music',
}

const STATUS_LABELS: Record<string, string> = {
    FINISHED: 'Finished',
    RELEASING: 'Releasing',
    NOT_YET_RELEASED: 'Not Yet Released',
    CANCELLED: 'Cancelled',
    HIATUS: 'Hiatus',
}

const MONTHS = [
    'Jan',
    'Feb',
    'Mar',
    'Apr',
    'May',
    'Jun',
    'Jul',
    'Aug',
    'Sep',
    'Oct',
    'Nov',
    'Dec',
]

export function formatEnum(value: string | null | undefined): string {
    if (!value) return '—'
    return value
        .toLowerCase()
        .split('_')
        .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
        .join(' ')
}

export function formatType(value: string | null | undefined): string {
    if (!value) return '—'
    return FORMAT_LABELS[value] ?? formatEnum(value)
}

export function formatStatus(value: string | null | undefined): string {
    if (!value) return '—'
    return STATUS_LABELS[value] ?? formatEnum(value)
}

export function formatDate(date: FuzzyDate | null | undefined): string {
    if (!date || !date.year) return '—'
    const month = date.month ? MONTHS[date.month - 1] : ''
    const day = date.day ?? ''
    return [month, day ? `${day},` : '', date.year].filter(Boolean).join(' ')
}

export function formatDuration(minutes: number | null | undefined): string {
    if (!minutes) return '—'
    return `${minutes} min`
}

export function formatNumber(value: number | null | undefined): string {
    if (value == null) return '—'
    return new Intl.NumberFormat('en-US').format(value)
}

export function stripHtml(html: string | null | undefined): string {
    if (!html) return ''
    return html
        .replace(/\n/g, '')
        .replace(/<br\s*\/?>/gi, '\n')
        .replace(/<[^>]+>/g, '')
        .replace(/&nbsp;/g, ' ')
        .replace(/&amp;/g, '&')
        .replace(/&quot;/g, '"')
        .replace(/&#039;/g, "'")
        .replace(/&mdash;/g, '—')
        .trim()
}