package user

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"profile/internal/cfg"
	"profile/platform/dynamo"
	"strconv"
)

type Repository interface {
	CreateUser(user *User) (*User, error)
	FindUserById(id int) (*User, error)
	UpdateUser(user *User) (*User, error)
	ListUsers(userIDs []int64) ([]*User, error)
	DeleteUser(id int) error
}

type repository struct {
	db  dynamo.Client
	cfg *cfg.Config
}

func (r repository) CreateUser(user *User) (*User, error) {
	value, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	item, err := r.db.DB().PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(r.cfg.DynamodbConfig.UserTable),
		Item:      value,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(item.Attributes, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r repository) ExistWithCpf(cpf string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindUserById(id int) (*User, error) {
	value, err := r.db.DB().GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(r.cfg.DynamodbConfig.UserTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},
		},
	})
	if err != nil {
		return nil, err
	}

	var user User
	err = attributevalue.UnmarshalMap(value.Item, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r repository) UpdateUser(user *User) (*User, error) {
	upd := expression.
		Set(expression.Name("Name"), expression.Value(user.Name)).
		Set(expression.Name("Email"), expression.Value(user.Email)).
		Set(expression.Name("Adress"), expression.Value(user.Address)).
		Set(expression.Name("Phone"), expression.Value(user.Phone)).
		Set(expression.Name("Birthday"), expression.Value(user.Birthday))

	exp, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return nil, errors.New("failed to build expression")
	}
	item, err := r.db.DB().UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: aws.String(r.cfg.DynamodbConfig.UserTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberN{Value: strconv.FormatInt(user.Id, 10)},
		},
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		UpdateExpression:          exp.Update(),
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(item.Attributes, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r repository) ListUsers(userIds []int64) ([]*User, error) {
	keys := make([]map[string]types.AttributeValue, len(userIds))
	for i, v := range userIds {
		keys[i] = map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: strconv.FormatInt(v, 10)},
		}
	}

	value, err := r.db.DB().BatchGetItem(context.Background(), &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			r.cfg.DynamodbConfig.UserTable: {
				Keys: keys,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	listUser := make([]*User, len(value.Responses[r.cfg.DynamodbConfig.UserTable]))
	for i := range value.Responses[r.cfg.DynamodbConfig.UserTable] {
		err = attributevalue.UnmarshalMap(value.Responses[r.cfg.DynamodbConfig.UserTable][i], &listUser[i])
		if err != nil {
			return nil, err
		}
	}
	return listUser, nil
}

func (r repository) DeleteUser(id int) error {
	_, err := r.db.DB().DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(r.cfg.DynamodbConfig.UserTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},
		},
	})
	return err
}

func NewRepository(db dynamo.Client) Repository {
	return &repository{
		db: db,
	}
}
