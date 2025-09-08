package media

import "fmt"

const mediaFields = `
	id
	title {
		romaji
		english
		native
	}
	type
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
	chapters
	volumes
	countryOfOrigin
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
	isAdult
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
`

var DiscoverMedia = fmt.Sprintf(`
	query DiscoverMedia($page: Int) {
		Page(page: $page, perPage: 50) {
			pageInfo {
				currentPage
				hasNextPage
			}
			media {
				%s
			}
		}
	}
`, mediaFields)

var UpdateFromMediaList = fmt.Sprintf(`
	query UpdateFromMediaList($ids: [Int], $page: Int) {
		Page(page: $page, perPage: 50) {
			pageInfo {
				currentPage
				hasNextPage
			}
			media(id_in: $ids) {
				%s
			}
		} 
	}
`, mediaFields)

var DiscoverNewMedia = fmt.Sprintf(`
	query DiscoverNewMedia($page: Int) {
		Page(page: $page, perPage: 50) {
			pageInfo {
				currentPage
				hasNextPage
			}
			media(sort: ID_DESC) {
			  %s
			}
		}
	}
`, mediaFields)
