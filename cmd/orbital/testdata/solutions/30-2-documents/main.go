package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Sample struct {
	ID   int64  `bson:"_id"`
	Tag  string `bson:"tag"`
	Mass int    `bson:"mass"`
}

var ErrSample = errors.New("invalid sample")
var ErrSampleMissing = errors.New("sample missing")
var ErrSampleDuplicate = errors.New("duplicate sample")

type SampleStore interface {
	Insert(context.Context, Sample) error
	Get(context.Context, int64) (Sample, error)
	Replace(context.Context, Sample) error
	Remove(context.Context, int64) error
	After(context.Context, int64, int) ([]Sample, error)
}
type Archive struct {
	store   SampleStore
	timeout time.Duration
}

func NewArchive(store SampleStore, timeout time.Duration) *Archive {
	return &Archive{store: store, timeout: timeout}
}
func normalizedSample(s Sample) (Sample, error) {
	s.Tag = strings.ToLower(strings.TrimSpace(s.Tag))
	if s.ID <= 0 || s.Tag == "" || s.Mass < 0 {
		return Sample{}, ErrSample
	}
	return s, nil
}
func (a *Archive) Add(ctx context.Context, s Sample) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s, err := normalizedSample(s)
	if err != nil {
		return err
	}
	bounded, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	return a.store.Insert(bounded, s)
}
func (a *Archive) Update(ctx context.Context, s Sample) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s, err := normalizedSample(s)
	if err != nil {
		return err
	}
	bounded, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	return a.store.Replace(bounded, s)
}
func (a *Archive) Get(ctx context.Context, id int64) (Sample, error) {
	if err := ctx.Err(); err != nil {
		return Sample{}, err
	}
	if id <= 0 {
		return Sample{}, ErrSample
	}
	bounded, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	return a.store.Get(bounded, id)
}
func (a *Archive) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if id <= 0 {
		return ErrSample
	}
	bounded, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	return a.store.Remove(bounded, id)
}
func (a *Archive) Page(ctx context.Context, after int64, limit int) ([]Sample, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	if after < 0 || limit < 1 || limit > 50 {
		return nil, 0, ErrSample
	}
	bounded, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	rows, err := a.store.After(bounded, after, limit+1)
	if err != nil {
		return nil, 0, err
	}
	if rows == nil {
		rows = []Sample{}
	}
	if len(rows) > limit {
		return rows[:limit], rows[limit-1].ID, nil
	}
	return rows, 0, nil
}

type MongoSamples struct{ collection *mongo.Collection }

func NewMongoSamples(collection *mongo.Collection) *MongoSamples {
	return &MongoSamples{collection: collection}
}
func (s *MongoSamples) Init(ctx context.Context) error {
	_, err := s.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "tag", Value: 1}}, Options: options.Index().SetUnique(true).SetName("unique_sample_tag"),
	})
	return err
}
func sampleMongoError(err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrSampleMissing
	}
	if mongo.IsDuplicateKeyError(err) {
		return ErrSampleDuplicate
	}
	return err
}
func (s *MongoSamples) Insert(ctx context.Context, sample Sample) error {
	_, err := s.collection.InsertOne(ctx, sample)
	return sampleMongoError(err)
}
func (s *MongoSamples) Get(ctx context.Context, id int64) (Sample, error) {
	var sample Sample
	err := s.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&sample)
	if err != nil {
		return Sample{}, sampleMongoError(err)
	}
	return sample, nil
}
func (s *MongoSamples) Replace(ctx context.Context, sample Sample) error {
	replaced, err := s.collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: sample.ID}}, sample)
	if err != nil {
		return sampleMongoError(err)
	}
	if replaced.MatchedCount == 0 {
		return ErrSampleMissing
	}
	return nil
}
func (s *MongoSamples) Remove(ctx context.Context, id int64) error {
	removed, err := s.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return sampleMongoError(err)
	}
	if removed.DeletedCount == 0 {
		return ErrSampleMissing
	}
	return nil
}
func (s *MongoSamples) After(ctx context.Context, after int64, limit int) ([]Sample, error) {
	cursor, err := s.collection.Find(ctx, bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: after}}}}, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	rows := []Sample{}
	for cursor.Next(ctx) {
		var sample Sample
		if err := cursor.Decode(&sample); err != nil {
			return nil, err
		}
		rows = append(rows, sample)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return rows, nil
}
func main() {}
