package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/minio/minio-go"
)

var servingEndpoints = []string{
	"localhost:9001",
	"localhost:9002",
	"localhost:9003",
	"localhost:9004",
	"localhost:9005",
	"localhost:9006",
}

var seededFiles = []string{
	"seed_1.txt",
	"seed_2.txt",
	"seed_3.txt",
	"seed_4.txt",
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func NewMinioConfig() *MinioConfig {

	return &MinioConfig{
		Endpoint:  servingEndpoints[rand.Intn(len(servingEndpoints))],
		AccessKey: "minioaccesskey",
		SecretKey: "miniosecretkey",
	}

}

func main() {

	rand.Seed(time.Now().Unix())

	minioConfig := NewMinioConfig()

	minioClient, err := minio.New(minioConfig.Endpoint, minioConfig.AccessKey, minioConfig.SecretKey, false)
	if err != nil {
		fmt.Println("[ERROR] Unable to create minio connections: ", err)
		os.Exit(1)
	}

	minioBucketInfo, err := minioClient.ListBuckets()
	if err != nil {
		fmt.Println("[ERROR] Unable to list buckets: ", err)
		os.Exit(1)
	}

	// Likely fresh cluster start
	if len(minioBucketInfo) == 0 {
		if err := minioClient.MakeBucket("test-bucket", ""); err != nil {
			fmt.Println("[ERROR] Unable to make bucket `test-bucket`: ", err)
			os.Exit(1)
		}
	}

	for _, file := range seededFiles {

		filePath := fmt.Sprintf("to_seed/%s", file)
		fileToSeed, err := os.Open(filePath)
		if err != nil {
			fmt.Println("[ERROR] Unable to open file... ", filePath, err)
			os.Exit(1)
		}

		fileInfo, err := fileToSeed.Stat()
		if err != nil {
			fmt.Println("[ERROR] Unable to stat file...", filePath, err)
			os.Exit(1)
		}

		bytesWritten, err := minioClient.PutObject("test-bucket", file, fileToSeed, fileInfo.Size(), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println("[ERROR] Unable to upload file to minio.", err)
			os.Exit(1)
		}

		fmt.Printf("[INFO] Created file %s - wrote %d bytes\n", file, bytesWritten)
	}

}
