package mapper

import (
	"time"
	"uptime_monitor/application/dto"
	"uptime_monitor/model"
)

// ToSiteResponse converts model.Site to dto.SiteResponse
func ToSiteResponse(site *model.Site) dto.SiteResponse {
	return dto.SiteResponse{
		ID:        site.ID,
		URL:       site.URL,
		CreatedAt: site.CreatedAt.Format(time.RFC3339),
	}
}

// ToSiteListResponse converts []model.Site to dto.SiteListResponse
func ToSiteListResponse(sites []model.Site) dto.SiteListResponse {
	responses := make([]dto.SiteResponse, 0, len(sites))
	for i := range sites {
		responses = append(responses, ToSiteResponse(&sites[i]))
	}
	return dto.SiteListResponse{
		Sites: responses,
	}
}
