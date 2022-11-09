package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RideRepository struct {
	datastore *mongo.Collection
}

func NewRideRepository(ds *mongo.Collection) *RideRepository {
	return &RideRepository{
		datastore: ds,
	}
}

func (r *RideRepository) FindById(ctx context.Context, uuid string) (RideLocationSchema, error) {
	var rideLocation RideLocationSchema

	if err := r.datastore.FindOne(ctx, bson.M{"_id": uuid}).Decode(&rideLocation); err != nil {
		return RideLocationSchema{}, nil
	}

	return rideLocation, nil
}

func (r *RideRepository) UpdateLocation(ctx context.Context, uuid string, lat, lon float64) (RideLocationSchema, error) {
	cord := Coordinates{
		Latitude:  lat,
		Longitude: lon,
	}

	schema := RideLocationSchema{
		UUID:        uuid,
		Coordinates: cord,
	}

	updateWith := bson.D{{Key: "$set", Value: schema}}

	filter := bson.M{"_id": uuid}
	opts := options.Update().SetUpsert(true)

	// @todo need to update list name
	_, err := r.datastore.UpdateOne(ctx, filter, updateWith, opts)

	if err != nil {
		return RideLocationSchema{}, err
	}

	return schema, err
}
