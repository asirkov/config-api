package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type ConfigDao interface {
	List(limit, offset int) ([]*ConfigModel, error)
	Count() (*int, error)
	Get(string) (*ConfigModel, error)
	Create(*ConfigModel) (*ConfigModel, error)
	Update(string, *ConfigModel) (*ConfigModel, error)
	Delete(string) error
	Exists(string) bool
}

type MongoDbConfigDao struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewMongoDbConfigDao(ctx context.Context, db *mongo.Collection) ConfigDao {
	return &MongoDbConfigDao{
		ctx:        ctx,
		collection: db,
	}
}

func (s *MongoDbConfigDao) List(limit, offset int) ([]*ConfigModel, error) {
	var results []*ConfigModel

	options := options.Find()
	options.SetLimit(int64(limit))
	options.SetSkip(int64(offset))

	filter := bson.M{}

	cursor, err := s.collection.Find(s.ctx, filter, options)
	if err != nil {
		return nil, err
	}

	for cursor.Next(s.ctx) {
		var result ConfigModel
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}
	defer cursor.Close(s.ctx)

	return results, nil
}

func (d *MongoDbConfigDao) Count() (*int, error) {
	var count int64

	filter := bson.M{}

	count, err := d.collection.CountDocuments(d.ctx, filter)
	if err != nil {
		return nil, err
	}
	countValue := int(count)
	return &countValue, nil
}

func (d *MongoDbConfigDao) Get(id string) (*ConfigModel, error) {
	var result ConfigModel

	if err := d.collection.FindOne(d.ctx, bson.M{"_id": id}).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *MongoDbConfigDao) Create(config *ConfigModel) (*ConfigModel, error) {
	if _, err := d.collection.InsertOne(d.ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (d *MongoDbConfigDao) Update(id string, config *ConfigModel) (*ConfigModel, error) {
	if _, err := d.collection.DeleteOne(d.ctx, bson.M{"_id": id}); err != nil {
		return nil, err
	}
	if _, err := d.collection.InsertOne(d.ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (d *MongoDbConfigDao) Delete(id string) error {
	if _, err := d.collection.DeleteOne(d.ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}

func (d *MongoDbConfigDao) Exists(id string) bool {
	return d.collection.FindOne(d.ctx, bson.M{"_id": id}).Err() != mongo.ErrNoDocuments
}
