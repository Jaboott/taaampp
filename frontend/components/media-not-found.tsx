import {Frown} from 'lucide-react'

export function MediaNotFound({ mediaType }: { mediaType: string }) {
    return (
        <div className="mx-auto flex max-w-md flex-col items-center gap-3 px-4 py-32 text-center">
      <span className="flex size-14 items-center justify-center rounded-2xl bg-secondary text-muted-foreground">
        <Frown className="size-14" aria-hidden="true"/>
      </span>
            <h1 className="font-display text-xl font-semibold text-foreground">No {mediaType} found</h1>
            <p className="text-pretty text-sm text-muted-foreground">
                Something went wrong fetching {mediaType} data. Please try again.
            </p>
        </div>
    )
}