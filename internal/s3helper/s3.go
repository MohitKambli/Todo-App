package s3helper

import (
    "bytes"
    "context"
    "fmt"
    "mime/multipart"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var bucketName string

func InitS3() {
    bucketName = os.Getenv("AWS_S3_BUCKET")

    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
    if err != nil {
        panic(fmt.Sprintf("Unable to load AWS SDK config: %v", err))
    }

    s3Client = s3.NewFromConfig(cfg)
}

// UploadFile uploads a file to S3 and returns the URL
func UploadFile(file multipart.File, fileName string) (string, error) {
    buffer := new(bytes.Buffer)
    _, err := buffer.ReadFrom(file)
    if err != nil {
        return "", err
    }

    uploader := manager.NewUploader(s3Client)
    _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(fileName),
        Body:   bytes.NewReader(buffer.Bytes()),
        ACL:    "public-read",
    })

    if err != nil {
        return "", err
    }

    fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, os.Getenv("AWS_REGION"), fileName)
    return fileURL, nil
}

// DeleteFile removes a file from S3
func DeleteFile(fileName string) error {
    _, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(fileName),
    })
    return err
}
