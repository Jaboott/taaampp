package media

const DiscoverAnime = `
	query DiscoverAnime($page: Int) {
		Page(page: $page, perPage: 1) {
			pageInfo {
				currentPage
				hasNextPage
			}
			media(type: ANIME) {
				id
				title {
					romaji
					english
					native
				}
				format
				status
				description
				startDate {
					year
					month
					day
				}
				endDate {
					year
					month
					day
				}
				season
				seasonYear
				episodes
				duration
				source(version: 3)
				trailer {
					id
					site
				}
				coverImage {
					large
				}
				bannerImage
				genres
				averageScore
				popularity
				trending
				favourites
				studios(isMain: true){
					nodes {
						name
					}
				}
				airingSchedule(notYetAired: true) {
					nodes {
						airingAt
						episode
					}
				}
				recommendations(sort: RATING_DESC perPage: 5) {
					nodes {
						rating
						mediaRecommendation {
							id
						}
					}
				}
				stats {
					scoreDistribution {
						score
						amount
					}   
				}
	
			}
		}
	}
`
