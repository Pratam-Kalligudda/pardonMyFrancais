package controllers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler for uploading files
func UploadHandler(c echo.Context) error {
	// Source file
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error retrieving the file: "+err.Error())
	}
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error opening the file: "+err.Error())
	}
	defer src.Close()

	// Prepare a form that you will submit to the other server
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multipartWriter.CreateFormFile("file", file.Filename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating form file: "+err.Error())
	}

	// Copy the file data to the form
	if _, err = io.Copy(fileWriter, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error copying file data: "+err.Error())
	}
	multipartWriter.Close()

	// Create a request to send to the other server
	targetURL := "http://python-service:8000/upload/" // Replace with the actual URL
	request, err := http.NewRequest("POST", targetURL, &requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating request: "+err.Error())
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Send the request to the other server
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error sending file to other server: "+err.Error())
	}
	defer response.Body.Close()

	// Check the server response
	if response.StatusCode != http.StatusOK {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Server error: %s", response.Status))
	}

	// Read the server response body
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error reading response from other server: "+err.Error())
	}

	// Return the response to the client
	return c.String(http.StatusOK, string(responseData))
}
