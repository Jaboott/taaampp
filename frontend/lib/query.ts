export interface FuzzyDate {
    year: number | null
    month: number | null
    day: number | null
}

interface Titles {
    english: string | null
    native: string | null
    romaji: string | null
}

export interface Media {
    average_score: number | null
    chapters: number | null
    cover_image: string | null
    episodes: number | null
    format: string | null
    genres: string[]
    id: number
    is_adult: boolean
    last_updated: string
    season: string | null
    season_year: number | null
    status: string
    studios: string[]
    titles: Titles
    type: string
    volumes: number | null
}

export interface MediaDetails {
    airing_schedule: {
        airing_at: number
        episode: number
    } | null
    average_score: number | null
    banner_image: string | null
    chapters: number | null
    country: string | null
    cover_image: string | null
    description: string | null
    duration: number | null
    end_date: FuzzyDate | null
    episodes: number | null
    favourites: number
    format: string | null
    genres: string[]
    id: number
    is_adult: boolean
    last_updated: string
    popularity: number
    recommendations: {
        id: number
        likes: number
    }[] | null
    score_distribution: {
        score: number
        amount: number
    }[] | null
    season: string | null
    season_year: number | null
    source: string | null
    start_date: FuzzyDate | null
    status: string | null
    studios: string[] | null
    titles: Titles
    trailer: string | null
    trending: number
    type: string
    volumes: number | null
}

interface MediaResult {
    media: MediaDetails
    recommendations: Media[]
}

export async function getMedia(id: number): Promise<MediaResult | null> {
    try {
        const mediaResponse = await fetch(`http://127.0.0.1:5000/api/media_detail/${id}`)

        if (!mediaResponse.ok) {
            console.error(`Failed to fetch media ${id}: ${mediaResponse.status}`)
            return null
        }

        const mediaJson: { data?: MediaDetails | null } = await mediaResponse.json()
        const media = mediaJson.data

        if (!media) return null

        const recommendationIds = media.recommendations?.map(({ id }) => id) ?? []

        if (recommendationIds.length === 0) {
            return {
                media,
                recommendations: [],
            }
        }

        const recommendationsResponse = await fetch(
            'http://127.0.0.1:5000/api/medias/batch',
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    ids: recommendationIds,
                }),
            }
        )

        if (!recommendationsResponse.ok) {
            console.error(
                `Failed to fetch recommendations for media ${id}: ` +
                `${recommendationsResponse.status}`
            )

            return {
                media,
                recommendations: [],
            }
        }

        const recommendationsJson: { data?: Media[] | null } = await recommendationsResponse.json()

        return {
            media,
            recommendations: recommendationsJson.data ?? [],
        }
    } catch (error) {
        console.error(`Failed to fetch media ${id}:`, error)
        return null
    }
}
