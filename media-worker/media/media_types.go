package media

type AnimeQueryResponse struct {
	Page struct {
		PageInfo struct {
			CurrentPage int  `json:"currentPage"`
			HasNextPage bool `json:"hasNextPage"`
		} `json:"pageInfo"`
		Media []struct {
			ID     int `json:"id"`
			Titles struct {
				Romaji  string `json:"romaji"`
				English string `json:"english"`
				Native  string `json:"native"`
			} `json:"title"`
			Format      string `json:"format"`
			Status      string `json:"status"`
			Description string `json:"description"`
			StartDate   struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"startDate"`
			EndDate struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"endDate"`
			Season     string `json:"season"`
			SeasonYear int    `json:"seasonYear"`
			Episodes   int    `json:"episodes"`
			Duration   int    `json:"duration"`
			Source     string `json:"source"`
			Trailer    struct {
				ID   string `json:"id"`
				Site string `json:"site"`
			} `json:"trailer"`
			CoverImage struct {
				Large string `json:"large"`
			} `json:"coverImage"`
			BannerImage  string   `json:"bannerImage"`
			Genres       []string `json:"genres"`
			AverageScore int      `json:"averageScore"`
			Popularity   int      `json:"popularity"`
			Trending     int      `json:"trending"`
			Favourites   int      `json:"favourites"`
			Studios      struct {
				Nodes []struct {
					Name string `json:"name"`
				} `json:"nodes"`
			} `json:"studios"`
			AiringSchedule struct {
				Nodes []struct {
					AiringAt int `json:"airingAt"`
					Episode  int `json:"episode"`
				} `json:"nodes"`
			} `json:"airingSchedule"`
			Recommendations struct {
				Nodes []struct {
					Rating              int `json:"rating"`
					MediaRecommendation struct {
						ID int `json:"id"`
					} `json:"mediaRecommendation"`
				} `json:"nodes"`
			} `json:"recommendations"`
			Stats struct {
				ScoreDistribution []struct {
					Score  int `json:"score"`
					Amount int `json:"amount"`
				} `json:"scoreDistribution"`
			} `json:"stats"`
		} `json:"media"`
	} `json:"Page"`
}
