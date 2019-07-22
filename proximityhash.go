package proximityhash

import (
	"math"

	"github.com/mmcloughlin/geohash"
)

func inCircleCheck(latitude float64, longitude float64, centreLat float64, centreLon float64, radius float64) bool {
	xDiff := longitude - centreLon
	yDiff := latitude - centreLat

	if math.Pow(xDiff, 2)+math.Pow(yDiff, 2) <= math.Pow(radius, 2) {
		return true
	}

	return false
}

func getCentroid(latitude float64, longitude float64, height float64, width float64) (float64, float64) {
	yCen := latitude + (height / 2)
	xCen := longitude + (width / 2)

	return xCen, yCen
}

func convertToLatLon(y float64, x float64, latitude float64, longitude float64) (float64, float64) {
	rEarth := 6371000.0

	latDiff := (y / rEarth) * (180 / math.Pi)
	lonDiff := (x / rEarth) * (180 / math.Pi) / math.Cos(latitude*math.Pi/180)

	finalLat := latitude + latDiff
	finalLon := longitude + lonDiff

	if finalLon > 180 {
		finalLon = finalLon - 360
	}
	if finalLat < -180 {
		finalLon = 360 + finalLon
	}

	return finalLat, finalLon
}

// CreateGeohash creates proximity hashes around the lat/lon with size radius.
// The geohash precisions are determined by the precision parameter and the initial parameter is used as the initialization
// value for the proximity hash
func CreateGeohash(latitude float64, longitude float64, radius float64, initial float64, precision int) map[string]float64 {

	geohashes := make(map[string]float64)

	gridWidth := [12]float64{5009400.0, 1252300.0, 156500.0, 39100.0, 4900.0, 1200.0, 152.9, 38.2, 4.8, 1.2, 0.149, 0.0370}
	gridHeight := [12]float64{4992600.0, 624100.0, 156000.0, 19500.0, 4900.0, 609.4, 152.4, 19.0, 4.8, 0.595, 0.149, 0.0199}

	height := (gridHeight[precision-1]) / 2
	width := (gridWidth[precision-1]) / 2

	latMoves := int(math.Ceil(radius / height))
	lonMoves := int(math.Ceil(radius / width))

	for i := 0; i < latMoves; i++ {
		tempLat := height * float64(i)
		for j := 0; j < lonMoves; j++ {
			tempLon := width * float64(j)
			if inCircleCheck(tempLat, tempLon, 0, 0, radius) {
				var lat, lon float64
				xCen, yCen := getCentroid(tempLat, tempLon, height, width)
				lat, lon = convertToLatLon(yCen, xCen, latitude, longitude)
				geohashes[geohash.EncodeWithPrecision(lat, lon, uint(precision))] = initial
				lat, lon = convertToLatLon(-yCen, xCen, latitude, longitude)
				geohashes[geohash.EncodeWithPrecision(lat, lon, uint(precision))] = initial
				lat, lon = convertToLatLon(yCen, -xCen, latitude, longitude)
				geohashes[geohash.EncodeWithPrecision(lat, lon, uint(precision))] = initial
				lat, lon = convertToLatLon(-yCen, -xCen, latitude, longitude)
				geohashes[geohash.EncodeWithPrecision(lat, lon, uint(precision))] = initial
			}
		}
	}

	return geohashes
}
