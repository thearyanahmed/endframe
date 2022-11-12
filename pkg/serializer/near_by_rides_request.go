package serializer

import (
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"strconv"
)

type NearByRidesRequest struct {
	X1 string  `json:"x1" schema:"x1"`
	X3 float64 `json:"x3" schema:"x3"`
	Y1 float64 `json:"y1" schema:"y1"`
	Y3 float64 `json:"y3" schema:"y3"`
}

func (r *NearByRidesRequest) Rules() govalidator.MapData {
	return govalidator.MapData{
		"x1": []string{"required", "lat"},
		"x3": []string{"required", "lat"},
		"y1": []string{"required", "lon"},
		"y3": []string{"required", "lon"},
	}
}

func (_ *NearByRidesRequest) ToArea(r *http.Request) entity.Area {
	// ignoring error, because request should be validated
	x1, _ := strconv.ParseFloat(r.URL.Query().Get("x1"), 64)
	y1, _ := strconv.ParseFloat(r.URL.Query().Get("y1"), 64)
	x3, _ := strconv.ParseFloat(r.URL.Query().Get("x3"), 64)
	y3, _ := strconv.ParseFloat(r.URL.Query().Get("y3"), 64)

	return entity.Area{
		X1Y1: entity.Coordinate{Lat: x1, Lon: y1},
		//X2Y2: entity.Coordinate{Lat: r.X2, Lon: r.Y2},
		X3Y3: entity.Coordinate{Lat: x3, Lon: y3},
		//X4Y4: entity.Coordinate{Lat: r.X4, Lon: r.Y4},
	}
}
