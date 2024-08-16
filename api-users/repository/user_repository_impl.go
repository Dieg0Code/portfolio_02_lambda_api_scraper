package repository

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dieg0code/shared/models"
	"github.com/sirupsen/logrus"
)

type UserRepositoryImpl struct {
	db        dynamodbiface.DynamoDBAPI
	tableName string
}

// GetByEmail implements UserRepository.
func (u *UserRepositoryImpl) GetByEmail(email string) (models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              &u.tableName,
		IndexName:              aws.String("EmailIndex"),
		KeyConditionExpression: aws.String("Email = :email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(email),
			},
		},
	}

	result, err := u.db.Query(input)
	if err != nil {
		logrus.WithError(err).Error("[UserRepositoryImpl.GetByEmail] error getting user")
		return models.User{}, errors.New("error getting user")
	}

	if len(result.Items) == 0 {
		return models.User{}, errors.New("user not found")
	}

	var user models.User
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		logrus.WithError(err).Error("[UserRepositoryImpl.GetByEmail] error unmarshalling user")
		return models.User{}, errors.New("error getting user")
	}

	return user, nil
}

// Create implements UserRepository.
func (u *UserRepositoryImpl) Create(user models.User) (models.User, error) {
	input := &dynamodb.PutItemInput{
		TableName: &u.tableName,
		Item: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(user.UserID),
			},
			"Username": {
				S: aws.String(user.Username),
			},
			"Email": {
				S: aws.String(user.Email),
			},
			"Password": {
				S: aws.String(user.Password),
			},
			"Role": {
				S: aws.String(user.Role),
			},
		},
	}

	_, err := u.db.PutItem(input)
	if err != nil {
		logrus.WithError(err).Error("[UserRepositoryImpl.Create] error creating user")
		return models.User{}, errors.New("error creating user")
	}

	return user, nil
}

// GetAll implements UserRepository.
func (u *UserRepositoryImpl) GetAll() ([]models.User, error) {
	input := &dynamodb.ScanInput{
		TableName: &u.tableName,
	}

	result, err := u.db.Scan(input)
	if err != nil {
		logrus.WithError(err).Error("[UserRepositoryImpl.GetAll] error getting users")
		return nil, errors.New("error getting users")
	}

	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		logrus.WithError(err).Error("[UserRepositoryImpl.GetAll] error unmarshalling users")
		return nil, errors.New("error getting users")
	}

	return users, nil
}

// GetByID implements UserRepository.
func (u *UserRepositoryImpl) GetByID(id string) (models.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: &u.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(id),
			},
		},
	}

	result, err := u.db.GetItem(input)
	if err != nil {
		logrus.WithError(err).Error("[UserRepositoryImpl.GetByID] error getting user")
		return models.User{}, errors.New("error getting user")
	}

	if result.Item == nil {
		return models.User{}, errors.New("user not found")
	}

	var user models.User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		logrus.WithError(err).Error("[UserRepositoryImpl.GetByID] error unmarshalling user")
		return models.User{}, errors.New("error getting user")
	}

	return user, nil
}

func NewUserRepositoryImpl(db dynamodbiface.DynamoDBAPI, tableName string) UserRepository {
	return &UserRepositoryImpl{
		db:        db,
		tableName: tableName,
	}
}
