import Link from "next/link";
import Image from "next/image"

export function SiteHeader() {
    return (
        <header className="fixed inset-x-0 top-0 z-30 border-b border-border bg-secondary/30 backdrop-blur-md">
            <div className="mx-auto flex max-w-6xl items-center gap-4 px-4 py-2 md:px-6">
                <Link href="/" className="flex items-center">
            <span className="relative flex size-16 items-center justify-center rounded-lg">
                {/*placeholder for logo*/}
                <Image src="/taaampp_logo.png" alt="taaampp logo" fill loading="eager"/>
            </span>
                    <span className="font-display text-lg font-bold tracking-tight text-foreground">
              taaampp
            </span>
                </Link>
            </div>
        </header>
    )
}
