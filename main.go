package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

// https://aws.github.io/aws-sdk-go-v2/docs/
func main() {
	log.Println("Starting ...")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load environment variables")
	}

	log.Println("Configuring AWS ...")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           os.Getenv("AWS_REGION_ENDPOINT"),
			SigningRegion: "default",
		}, nil
	})))
	if err != nil {
		log.Fatal("Error loading default AWS config: ", err.Error())
	}

	log.Println("Initializing S3 Client ...")
	s3client := s3.NewFromConfig(cfg)

	log.Println("Fetching list objects ...")
	output, err := s3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("thor-sandbox"),
	})
	if err != nil {
		log.Fatal("Error listing S3 objects: ", err.Error())
	}
	log.Println("First page results: ")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
	log.Println("Done!")
}
