package dynamodbdriver

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/pkg/errors"
)

// DynamoDB setting
type DynamoDB struct {
	TableName string
	dynamodbiface.DynamoDBAPI
}

// New DynamoDB
func New(p client.ConfigProvider, tableName string) *DynamoDB {
	return &DynamoDB{
		DynamoDBAPI: dynamodb.New(p),
		TableName:   tableName,
	}
}

// ScanAll get all table data
func (d *DynamoDB) ScanAll(out interface{}) error {
	var data []map[string]*dynamodb.AttributeValue
	err := d.ScanPages(
		&dynamodb.ScanInput{TableName: &d.TableName},
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			data = append(data, page.Items...)
			return true
		},
	)
	if err != nil {
		return err
	}
	return dynamodbattribute.UnmarshalListOfMaps(data, &out)
}

// Get is get from table
func (d *DynamoDB) Get(key string, out interface{}) error {
	input := &dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"ID": {S: &key}},
		TableName: &d.TableName,
	}
	result, err := d.GetItem(input)
	if err != nil {
		return errors.Wrapf(err, "Dynamodb GetItem table=%s,key=%s", d.TableName, key)
	}
	return dynamodbattribute.UnmarshalMap(result.Item, out)
}

// Put put to table
func (d *DynamoDB) Put(in interface{}) error {
	attr, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return errors.Wrap(err, "MarshalMap err")
	}
	_, err = d.PutItem(&dynamodb.PutItemInput{Item: attr})
	return err
}
