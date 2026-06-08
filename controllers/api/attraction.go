// controllers/api/attraction.go
package api

import (
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/robiulislam99/TravelSphere/services"
	"github.com/robiulislam99/TravelSphere/utils"
)

type AttractionAPIController struct {
	web.Controller
}

// GetByCoords handles GET /api/attractions?lat=...&lon=...&radius=...&limit=...
func (c *AttractionAPIController) GetByCoords() {
	latStr := c.GetString("lat")
	lonStr := c.GetString("lon")

	if latStr == "" || lonStr == "" {
		utils.SendBadRequest(&c.Controller, "lat and lon query params are required")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		utils.SendBadRequest(&c.Controller, "lat must be a valid number")
		return
	}
	if lat < -90 || lat > 90 {
		utils.SendBadRequest(&c.Controller, "lat must be between -90 and 90")
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		utils.SendBadRequest(&c.Controller, "lon must be a valid number")
		return
	}
	if lon < -180 || lon > 180 {
		utils.SendBadRequest(&c.Controller, "lon must be between -180 and 180")
		return
	}

	radius, _ := strconv.Atoi(c.GetString("radius"))
	limit, _ := strconv.Atoi(c.GetString("limit"))

	if radius <= 0 {
		radius = 10000
	} // default 10 km
	if radius > 50000 {
		radius = 50000
	} // OpenTripMap free tier cap
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	} // sanity cap

	attractions, err := services.Attractions().GetByCoords(lat, lon, radius, limit)
	if err != nil {
		utils.SendInternalError(&c.Controller)
		return
	}

	utils.SendSuccess(&c.Controller, attractions)
}
