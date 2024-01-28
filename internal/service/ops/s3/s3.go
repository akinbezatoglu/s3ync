package s3

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type BucketBasics struct {
	Clients map[string]*s3.Client
}

func NewBucketBasics() (*BucketBasics, error) {
	profiles := GetLocalAwsProfiles()
	b := BucketBasics{}
	for _, profile := range profiles {
		err := b.AddClient(profile)
		if err != nil {
			return nil, err
		}
	}
	return &b, nil
}

func (b *BucketBasics) AddClient(profile string) error {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return err
	}
	b.Clients[profile] = s3.NewFromConfig(cfg)
	return nil
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (b *BucketBasics) UploadFile(bucketName, objectKey, fileName, profile string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		_, err = b.Clients[profile].PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   bytes.NewReader(file),
		})
		if err != nil {
			fmt.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}
	return err
}

// DeleteFile deletes a file from S3.
func (b *BucketBasics) DeleteFile(bucket, key, profile string) error {
	_, err := b.Clients[profile].DeleteObject(context.Background(),
		&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	return err
}

// DeleteObjects deletes a list of objects from a bucket.
func (b *BucketBasics) DeleteObjects(bucketName, profile string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := b.Clients[profile].DeleteObjects(context.Background(),
		&s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &types.Delete{Objects: objectIds},
		})
	if err != nil {
		fmt.Printf("Couldn't delete objects from bucket %v. Here's why: %v\n", bucketName, err)
	}
	return err
}
