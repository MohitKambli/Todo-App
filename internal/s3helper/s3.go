package s3helper

import (
    "bytes"
    "context"
    "fmt"
    "log"
    "mime/multipart"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/joho/godotenv"
)

var s3Client *s3.Client
var bucketName string

// InitS3Session initializes the S3 client and retrieves the bucket name from the .env file
func InitS3() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: No .env file found")
    }

    // Retrieve AWS credentials and region from environment variables
    awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
    awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
    awsRegion := os.Getenv("AWS_REGION")
    bucketName = os.Getenv("AWS_S3_BUCKET_NAME")

    if awsAccessKey == "" || awsSecretKey == "" || awsRegion == "" || bucketName == "" {
        log.Fatal("AWS credentials or region or bucket name not set in environment variables")
    }

    // Create AWS session with credentials and region
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(awsRegion),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
    )
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    // Initialize S3 client with the session
    s3Client = s3.NewFromConfig(cfg)
}

// UploadFile uploads a file to S3 and returns the URL
func UploadFile(file multipart.File, fileName string) (string, error) {
    // Initialize S3 session if not already initialized
    if s3Client == nil {
        InitS3()
    }

    // Create a buffer to read the file contents into
    buffer := new(bytes.Buffer)
    _, err := buffer.ReadFrom(file)
    if err != nil {
        return "", err
    }

    // Create an uploader and upload the file
    uploader := manager.NewUploader(s3Client)
    _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(fileName),
        Body:   bytes.NewReader(buffer.Bytes()),
        ACL:    "public-read", // Set file ACL to public-read
    })

    if err != nil {
        return "", err
    }

    // Return the URL of the uploaded file
    fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, os.Getenv("AWS_REGION"), fileName)
    return fileURL, nil
}

// DeleteFile removes a file from S3
func DeleteFile(fileName string) error {
    if s3Client == nil {
        InitS3()
    }

    // Delete the file from S3
    _, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(fileName),
    })
    return err
}
