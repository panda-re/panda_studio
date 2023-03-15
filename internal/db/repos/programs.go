package repos

import (
	"context"

	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const PROGRAMS_TABLE string = "programs"

type ProgramRepository interface {
	FindAll(ctx context.Context) ([]models.InteractionProgram, error)
	FindOne(ctx context.Context, id db.ObjectID) (*models.InteractionProgram, error)
	Create(ctx context.Context, obj *models.InteractionProgram) (*models.InteractionProgram, error)
	Delete(ctx context.Context, id db.ObjectID) (*models.InteractionProgram, error)
	Update(ctx context.Context, id db.ObjectID, obj *models.InteractionProgram) (*models.InteractionProgram, error)
}

func GetProgramRepository(ctx context.Context) (ProgramRepository, error) {
	mongoClient, err := db.GetMongoDatabase(ctx)
	if err != nil {
		return nil, err
	}

	return &mongoProgramRepository{
		coll: mongoClient.Collection(PROGRAMS_TABLE),
	}, nil
}

type mongoProgramRepository struct {
	coll *mongo.Collection
}

// Create implements ProgramRepository
func (r *mongoProgramRepository) Create(ctx context.Context, obj *models.InteractionProgram) (*models.InteractionProgram, error) {
	obj.ID = db.NewObjectID()

	// insert into mongo
	result, err := r.coll.InsertOne(ctx, obj)
	if err != nil {
		return nil, err
	}

	insertedId := result.InsertedID.(primitive.ObjectID)
	obj.ID = &insertedId

	return obj, nil
}

// Delete implements ProgramRepository
func (r *mongoProgramRepository) Delete(ctx context.Context, id *primitive.ObjectID) (*models.InteractionProgram, error) {
	prog, err := r.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = r.coll.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return nil, errors.Wrap(err, "db error with deleting")
	}

	return prog, nil
}

// FindAll implements ProgramRepository
func (r *mongoProgramRepository) FindAll(ctx context.Context) ([]models.InteractionProgram, error) {
	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	var programs []models.InteractionProgram
	if err = cursor.All(ctx, &programs); err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return programs, nil
}

// FindOne implements ProgramRepository
func (r *mongoProgramRepository) FindOne(ctx context.Context, id *primitive.ObjectID) (*models.InteractionProgram, error) {
	var result models.InteractionProgram

	err := r.coll.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &result, nil
}

// Update implements ProgramRepository
func (r *mongoProgramRepository) Update(ctx context.Context, id *primitive.ObjectID, obj *models.InteractionProgram) (*models.InteractionProgram, error) {
	// ensure it exists
	_, err := r.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// prevent id from being changed
	obj.ID = id

	// update in mongo
	_, err = r.coll.UpdateOne(ctx, bson.M{
		"_id": id,
	}, obj)
	if err != nil {
		return nil, errors.Wrap(err, "db error with updating")
	}

	return obj, nil
}