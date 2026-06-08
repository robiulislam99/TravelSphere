// controllers/api/attraction.go
// AttractionAPIController — GET /api/attractions?lat=...&lon=...
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
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		utils.SendBadRequest(&c.Controller, "lon must be a valid number")
		return
	}

	radius, _ := strconv.Atoi(c.GetString("radius"))
	limit, _  := strconv.Atoi(c.GetString("limit"))

	attractions, _ := services.Attractions().GetByCoords(lat, lon, radius, limit)
	utils.SendSuccess(&c.Controller, attractions)
}