import Script from 'next/script'

export function ReloadOnHistoryNavigation() {
    return (
        <Script
            id="reload-on-history-navigation"
            strategy="beforeInteractive"
            dangerouslySetInnerHTML={{
                __html: `
                    window.addEventListener('pageshow', function (event) {
                        var navEntry = performance.getEntriesByType('navigation')[0];

                        var isBackForward =
                            event.persisted ||
                            (navEntry && navEntry.type === 'back_forward');

                        if (isBackForward) {
                            window.location.reload();
                        }
                    });
                `,
            }}
        />
    )
}