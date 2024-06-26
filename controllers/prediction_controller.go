package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/labstack/echo/v4"
)
type ServerResponse struct {
	Result struct {
			PronunciationScore float64 `json:"pronunciation_score"`
	} `json:"result"`
}


// // Handler for uploading files
// func UploadHandler(c echo.Context) error {
// 	// Source file
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Error retrieving the file: "+err.Error())
// 	}
// 	src, err := file.Open()
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Error opening the file: "+err.Error())
// 	}
// 	defer src.Close()

// 	// Prepare a form that you will submit to the other server
// 	var requestBody bytes.Buffer
// 	multipartWriter := multipart.NewWriter(&requestBody)
// 	fileWriter, err := multipartWriter.CreateFormFile("file", file.Filename)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating form file: "+err.Error())
// 	}

// 	// Copy the file data to the form
// 	if _, err = io.Copy(fileWriter, src); err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Error copying file data: "+err.Error())
// 	}
// 	multipartWriter.Close()

// 	// Create a request to send to the other server
// 	targetURL := "http://fastApi:8000/upload/" // Replace with the actual URL
// 	request, err := http.NewRequest("POST", targetURL, &requestBody)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating request: "+err.Error())
// 	}
// 	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

// 	// Send the request to the other server
// 	client := &http.Client{}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Error sending file to other server: "+err.Error())
// 	}
// 	defer response.Body.Close()

// 	// Check the server response
// 	if response.StatusCode != http.StatusOK {
// 		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Server error: %s", response.Status))
// 	}

// 	// Read the server response body
// 	responseData, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Error reading response from other server: "+err.Error())
// 	}

// 	// Return the response to the client
// 	return c.String(http.StatusOK, string(responseData))
// }

func UploadHandlers(c echo.Context) error {
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

	// Use the system's temporary directory
	tempDir := os.TempDir()

	// Save the video file to a temporary location
	videoPath := filepath.Join(tempDir, file.Filename)
	dst, err := os.Create(videoPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating temporary video file: "+err.Error())
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error saving video file: "+err.Error())
	}

	// Extract audio from the video
	audioOutputPath := filepath.Join(tempDir, "audio.wav")
	if err := extractAudio(videoPath, audioOutputPath); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error extracting audio: "+err.Error())
	}
	defer os.Remove(audioOutputPath) // Clean up the audio file after sending it

	// Prepare the audio file to be sent to the other server
	audioFile, err := os.Open(audioOutputPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error opening audio file: "+err.Error())
	}
	defer audioFile.Close()

	// Prepare a form that you will submit to the other server
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multipartWriter.CreateFormFile("file", "audio.wav")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating form file: "+err.Error())
	}

	// Copy the audio file data to the form
	if _, err = io.Copy(fileWriter, audioFile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error copying audio file data: "+err.Error())
	}
	multipartWriter.Close()

	// Create a request to send to the other server
	targetURL := "http://127.0.0.1:8000/upload/" // Replace with the actual URL
	request, err := http.NewRequest("POST", targetURL, &requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating request: "+err.Error())
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Send the request to the other server
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error sending audio file to other server: "+err.Error())
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
	return c.JSONBlob(http.StatusOK, responseData)
}

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

	// Use the system's temporary directory
	tempDir := os.TempDir()

	// Save the video file to a temporary location
	videoPath := filepath.Join(tempDir, file.Filename)
	dst, err := os.Create(videoPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating temporary video file: "+err.Error())
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error saving video file: "+err.Error())
	}

	// Extract audio from the video
	audioOutputPath := filepath.Join(tempDir, "audio.wav")
	if err := extractAudio(videoPath, audioOutputPath); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error extracting audio: "+err.Error())
	}
	defer os.Remove(audioOutputPath) // Clean up the audio file after sending it

	// Prepare the audio file to be sent to the other server
	audioFile, err := os.Open(audioOutputPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error opening audio file: "+err.Error())
	}
	defer audioFile.Close()

	// Prepare a form that you will submit to the other server
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multipartWriter.CreateFormFile("file", "audio.wav")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating form file: "+err.Error())
	}

	// Copy the audio file data to the form
	if _, err = io.Copy(fileWriter, audioFile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error copying audio file data: "+err.Error())
	}
	multipartWriter.Close()

	// Create a request to send to the other server
	targetURL := "http://fastApi:8000/upload/" // Replace with the actual URL
	request, err := http.NewRequest("POST", targetURL, &requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating request: "+err.Error())
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Send the request to the other server
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error sending audio file to other server: "+err.Error())
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

	// Unmarshal the response data to extract the pronunciation score
	var serverResponse ServerResponse
	if err := json.Unmarshal(responseData, &serverResponse); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing server response: "+err.Error())
	}

	// Return only the pronunciation score
	return c.JSON(http.StatusOK, map[string]float64{
		"pronunciation_score": serverResponse.Result.PronunciationScore,
	})
}

func extractAudio(videoPath string, audioOutputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-q:a", "0", "-map", "a", audioOutputPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error extracting audio: %w", err)
	}
	return nil
}