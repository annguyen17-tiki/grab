package service

import (
	"sort"

	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/jinzhu/copier"
	"github.com/mmcloughlin/geohash"
)

func (svc *service) SaveLocation(input *dto.SaveLocation) error {
	var location model.Location
	if err := copier.Copy(&location, input); err != nil {
		return err
	}

	location.GeoHash = geohash.EncodeWithPrecision(location.Latitude, location.Longitude, svc.cfg.GeoHashPrecision)
	return svc.store.Location().Save(&location)
}

func (svc *service) NearestLocations(input *dto.NearestLocations) ([]*model.Location, error) {
	geoHash := geohash.EncodeWithPrecision(input.Latitude, input.Longitude, svc.cfg.GeoHashPrecision)
	locations, err := svc.store.Location().Search(&dto.SearchLocations{
		GeoHashes:      append(geohash.Neighbors(geoHash), geoHash),
		Vehicles:       input.Vehicles,
		DriverStatuses: input.DriverStatuses,
	}, "Driver")
	if err != nil {
		return nil, err
	}

	sort.Slice(locations, func(i, j int) bool {
		return locations[i].DistanceTo(input.Latitude, input.Longitude) < locations[j].DistanceTo(input.Latitude, input.Longitude)
	})

	var nearestLocations []*model.Location
	for _, location := range locations {
		if input.Radius != 0 && location.DistanceTo(input.Latitude, input.Longitude) > input.Radius {
			break
		}

		nearestLocations = append(nearestLocations, location)
	}

	return nearestLocations, nil
}
