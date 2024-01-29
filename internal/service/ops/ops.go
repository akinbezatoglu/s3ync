package ops

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"gopkg.in/ini.v1"
)

type BucketBasics struct {
	Clients map[string]*s3.Client
}

func NewBucketBasics() (*BucketBasics, error) {
	profiles := GetLocalAwsProfiles()
	b := BucketBasics{
		Clients: make(map[string]*s3.Client),
	}
	for _, profile := range profiles {
		err := b.AddClient(profile)
		if err != nil {
			return nil, err
		}
	}
	return &b, nil
}

func GetLocalAwsProfiles() []string {
	fname := config.DefaultSharedConfigFilename()
	f, err := ini.Load(fname)
	if err != nil {
		return nil
	}
	profiles := f.SectionStrings()

	// There is no profile in config file
	if len(profiles) == 1 && profiles[0] == "DEFAULT" {
		return nil
	}

	// ["DEFAULT", "default", "user1", ...]
	return profiles[1:] // remove DEFAULT
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

func (b *BucketBasics) DeleteDirectory(bucket, key, profile string) error {
	ctx := context.Background()
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	}
	paginator := s3.NewListObjectsV2Paginator(b.Clients[profile], listObjectsInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		deleteObjectsInput := &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &types.Delete{
				Objects: make([]types.ObjectIdentifier, 0, len(page.Contents)),
			},
		}
		for _, object := range page.Contents {
			deleteObjectsInput.Delete.Objects = append(deleteObjectsInput.Delete.Objects, types.ObjectIdentifier{
				Key: object.Key,
			})
		}
		_, err = b.Clients[profile].DeleteObjects(ctx, deleteObjectsInput)
		if err != nil {
			return err
		}
	}
	// Delete the directory itself (empty prefix)
	_, err := b.Clients[profile].DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
