package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func init() {
	setUpEnv("")
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("BUCKET")
	key := os.Getenv("LOG_FILE")

	params := s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		ACL:    "private",
		Body:   bytes.NewReader(readFile(key)),
	}

	s3Client.PutObject(context.TODO(), &params)
}

func readFile(name string) []byte {
	content, err := ioutil.ReadFile(name)

	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	return content
}

func setUpEnv(envFile string) {
	if envFile == "" {
		envFile = ".env"
	}

	file, err := os.Open(envFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		line := scanner.Text()
		values := strings.Split(line, "=")
		os.Setenv(values[0], values[1])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
