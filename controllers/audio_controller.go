package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func FetchAudioHandler(c echo.Context) error {
	fileName := c.Param("fileName");
	var err = godotenv.Load()
	if err != nil {
		c.Logger().Fatal(err)
		return err
	}
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	s3Bucket := os.Getenv("s3Bucket")
	awsRegion := os.Getenv("awsRegion")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretKey, ""),
	})
	if err != nil {
		c.Logger().Fatal(err)
		return err
	}
	s3Client := s3.New(sess)
	// Fetch the audio file from S3
	params := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String("/pronunciation/"+fileName+".mp3"),
	}
	resp, err := s3Client.GetObject(params)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	defer resp.Body.Close()

	// Set response headers
	c.Response().Header().Set("Content-Type", "audio/mpeg")
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", *resp.ContentLength))
	c.Response().Header().Set("Content-Disposition", "inline; filename=audio-file.mp3")

	// Send the audio file as the response
	_, err = io.Copy(c.Response(), resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return nil
}
