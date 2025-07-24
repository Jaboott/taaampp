package media

import (
	"fmt"
	"strings"
	"time"
)

func FormatAnimeQueryResponse(resp AnimeQueryResponse) string {
	var sb strings.Builder
	pageInfo := resp.Page.PageInfo
	sb.WriteString(fmt.Sprintf("📄 Page %d | More pages: %v\n\n", pageInfo.CurrentPage, pageInfo.HasNextPage))

	for _, media := range resp.Page.Media {
		sb.WriteString("==================================\n")
		sb.WriteString(fmt.Sprintf("🎬 Title: %s (%s / %s)\n", media.Titles.Romaji, media.Titles.English, media.Titles.Native))
		sb.WriteString(fmt.Sprintf("🆔 ID: %d | 📅 Season: %s %d\n", media.ID, media.Season, media.SeasonYear))
		sb.WriteString(fmt.Sprintf("📊 Status: %s | Format: %s | Episodes: %d | Duration: %d\n",
			media.Status, media.Format, media.Episodes, media.Duration))
		sb.WriteString(fmt.Sprintf("🏷️ Genres: %s\n", strings.Join(media.Genres, ", ")))
		sb.WriteString(fmt.Sprintf("🎯 Score: %d | ❤️ Favorites: %d | 🔥 Trending: %d\n",
			media.AverageScore, media.Favourites, media.Trending))
		if len(media.Studios.Nodes) > 0 {
			sb.WriteString(fmt.Sprintf("🏢 Studio: %s\n", media.Studios.Nodes[0].Name))
		}
		sb.WriteString(fmt.Sprintf("📈 Popularity: %d\n", media.Popularity))
		sb.WriteString(fmt.Sprintf("🖼️ Cover: %s\n", media.CoverImage.Large))
		if media.BannerImage != "" {
			sb.WriteString(fmt.Sprintf("🖼️ Banner: %s\n", media.BannerImage))
		}

		start := formatDate(media.StartDate)
		end := formatDate(media.EndDate)
		sb.WriteString(fmt.Sprintf("🗓️ Aired: %s → %s\n", start, end))

		if len(media.AiringSchedule.Nodes) > 0 {
			air := media.AiringSchedule.Nodes[0]
			sb.WriteString(fmt.Sprintf("📺 Next Episode: %d at %s\n", air.Episode, time.Unix(int64(air.AiringAt), 0).Format(time.RFC1123)))
		}

		if len(media.Stats.ScoreDistribution) > 0 {
			sb.WriteString("📉 Score Distribution:\n")
			for _, s := range media.Stats.ScoreDistribution {
				sb.WriteString(fmt.Sprintf("  %2d: %d\n", s.Score, s.Amount))
			}
		}

		sb.WriteString("==================================\n\n")
	}

	return sb.String()
}

func formatDate(d struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}) string {
	if d.Year == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}
