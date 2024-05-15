package controllers

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "net/http"

  "github.com/labstack/echo/v4"
)

func UploadAudio(c echo.Context) error {
  // Access request data
  req := c.Request()

  // Check for POST method (optional, since route is already defined for POST)
  if req.Method != http.MethodPost {
    return echo.ErrMethodNotAllowed
  }

  // Read audio data from request body
  audioBytes, err := io.ReadAll(req.Body)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "Error reading audio data")
  }

  // **Optional:** Validate audio data format (if applicable)
  // You can add checks here to ensure the received data is in the expected format (e.g., WAV)

  // **Optional:** Preprocess audio data (if needed)
  // You can add preprocessing steps here (e.g., normalization, silence removal)

  // Send audio data to Django (assuming WAV format)
  err = sendAudioToDjango(audioBytes)
  if err != nil {
    return echo.NewHTTPError(http.StatusInternalServerError, "Error sending audio to Django")
  }

  // Send success response
  return c.String(http.StatusOK, "Audio received successfully")
}

func sendAudioToDjango(audioData []byte) error {
  // Prepare request data (send audio data directly)
  requestData := map[string]interface{}{
    "audio_data": audioData, // Send bytes directly (assuming WAV format)
    // Add other data as needed (e.g., session_id, metadata)
  }
  requestBody, err := json.Marshal(requestData)
  if err != nil {
    return err
  }

  // Create HTTP POST request
  req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8000/predict/", bytes.NewReader(requestBody))
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/json") // Set request content type

  // Send request and handle response (optional)
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  // Check response status code (optional)
  if resp.StatusCode != http.StatusOK {
    // Handle non-200 status code
    return fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
  }
  return nil
  // You can optionally handle the response
}
