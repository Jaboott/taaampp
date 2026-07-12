import {getMedia} from '@/lib/query'
import {SiteHeader} from '@/components/site-header'
import {AnimeHero} from '@/components/anime-hero'
import {StatsPanel} from '@/components/stats-panel'
import {ScoreDistribution} from '@/components/score-distribution'
import {Trailer} from '@/components/trailer'
import {Synopsis} from '@/components/synopsis'
import {Recommendations} from '@/components/recommendations'
import {MediaNotFound} from "@/components/media-not-found";

export default async function Page({params,}: { params: Promise<{ id: string }> }) {
    const {id} = await params
    const animeId = Number(id)
    const result = await getMedia(animeId)

    return (
        <main className="min-h-dvh pb-16">
            <SiteHeader/>

            {!result ? (
                <MediaNotFound mediaType = {"anime"}/>
            ) : (
                <>
                    <AnimeHero media={result.media}/>

                    <div
                        className="mx-auto mt-10 grid max-w-6xl grid-cols-1 gap-6 px-4 md:px-6 lg:grid-cols-[minmax(0,1fr)_20rem]">
                        {/* Main column */}
                        <div className="order-2 flex flex-col gap-6 lg:order-1">
                            <Synopsis description={result.media.description}/>
                            <Trailer media={result.media}/>
                            <ScoreDistribution media={result.media}/>
                            <Recommendations items={result.recommendations}/>
                        </div>

                         {/*Grouped stats sidebar*/}
                        <aside className="order-1 lg:order-2">
                            <div className="lg:sticky lg:top-20">
                                <StatsPanel media={result.media}/>
                            </div>
                        </aside>
                    </div>
                </>
            )}
        </main>
    )
}
