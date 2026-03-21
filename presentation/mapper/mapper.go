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

// ToCheckResponse converts model.Check to dto.CheckResponse
func ToCheckResponse(check *model.Check) dto.CheckResponse {
	return dto.CheckResponse{
		ID:           check.ID,
		Status:       check.Status,
		ResponseTime: check.ResponseTime,
		CheckedAt:    check.CheckedAt.Format(time.RFC3339),
	}
}

// ToSiteHistoryResponse converts model.Site and []model.Check to dto.SiteHistoryResponse
func ToSiteHistoryResponse(site *model.Site, checks []model.Check) dto.SiteHistoryResponse {
	checkResponses := make([]dto.CheckResponse, 0, len(checks))
	for i := range checks {
		checkResponses = append(checkResponses, ToCheckResponse(&checks[i]))
	}

	return dto.SiteHistoryResponse{
		SiteID: site.ID,
		URL:    site.URL,
		Checks: checkResponses,
	}
}

// ToSiteResponseWithCheck converts model.Site and optional model.Check to dto.SiteResponse
func ToSiteResponseWithCheck(site *model.Site, check *model.Check) dto.SiteResponse {
	response := dto.SiteResponse{
		ID:        site.ID,
		URL:       site.URL,
		CreatedAt: site.CreatedAt.Format(time.RFC3339),
	}

	if check != nil {
		response.Status = check.Status
		response.ResponseTime = check.ResponseTime
		response.LastChecked = check.CheckedAt.Format(time.RFC3339)
	}

	return response
}

// ToSiteListResponseWithChecks converts []model.Site and []*model.Check to dto.SiteListResponse
func ToSiteListResponseWithChecks(sites []model.Site, checks []*model.Check) dto.SiteListResponse {
	responses := make([]dto.SiteResponse, 0, len(sites))
	for i, site := range sites {
		var check *model.Check
		if i < len(checks) {
			check = checks[i]
		}
		responses = append(responses, ToSiteResponseWithCheck(&site, check))
	}

	return dto.SiteListResponse{
		Sites: responses,
	}
}
