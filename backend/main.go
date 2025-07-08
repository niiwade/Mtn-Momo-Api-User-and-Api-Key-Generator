package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Response structure for API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MomoKeyRequest structure for incoming requests
type MomoKeyRequest struct {
	PrimaryKey   string `json:"primaryKey"`   // Subscription Key (Ocp-Apim-Subscription-Key)
	SecondaryKey string `json:"secondaryKey"` // Optional secondary key
	CallbackHost string `json:"callbackHost"` // Provider callback host
}

// CreateUserResponse structure for API user creation response
type CreateUserResponse struct {
	UserID       string `json:"userId"`
	TargetEnv    string `json:"targetEnvironment"`
	CallbackHost string `json:"providerCallbackHost"`
}

// CreateKeyResponse structure for API key creation response
type CreateKeyResponse struct {
	APIKey string `json:"apiKey"`
}

// MomoKeyResponse structure for generated keys
type MomoKeyResponse struct {
	APIKey       string `json:"apiKey"`
	APIUser      string `json:"apiUser"`
	UserID       string `json:"userId"`
	CallbackHost string `json:"callbackHost"`
	DateTime     string `json:"dateTime"`
	TargetEnv    string `json:"targetEnvironment"`
	TestCommand  string `json:"testCommand,omitempty"` // Optional curl command for testing
	Base64Auth   string `json:"base64Auth,omitempty"`  // Base64 encoded auth string (apiUser:apiKey)
}

// createAPIUser calls the MTN MoMo API to create an API user
func createAPIUser(subscriptionKey string, callbackHost string) (string, error) {
	// Generate a UUID for the API user
	apiUser := uuid.New().String()
	log.Printf("Generated new API User UUID: %s", apiUser)

	// Create the request URL
	url := "https://sandbox.momodeveloper.mtn.com/v1_0/apiuser"
	log.Printf("Preparing API request to: %s", url)

	// Create the request body
	requestBody := map[string]string{
		"providerCallbackHost": callbackHost,
	}
	log.Printf("Request body includes callback host: %s", callbackHost)

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("ERROR: Failed to marshal request body: %v", err)
		return "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request: %v", err)
		return "", err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)
	req.Header.Set("X-Reference-Id", apiUser)
	log.Println("Added required headers: Content-Type, Ocp-Apim-Subscription-Key, X-Reference-Id")

	// Send the request
	log.Println("Sending API User creation request to MTN MoMo API...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: HTTP request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	log.Printf("Received response with status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("ERROR: API returned non-success status: %d, body: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("failed to create API user: %s, status: %d", string(body), resp.StatusCode)
	}

	log.Printf("API User created successfully with ID: %s", apiUser)
	return apiUser, nil
}

// createAPIKey calls the MTN MoMo API to create an API key for the given API user
func createAPIKey(subscriptionKey string, apiUser string) (string, error) {
	// Create the request URL
	url := fmt.Sprintf("https://sandbox.momodeveloper.mtn.com/v1_0/apiuser/%s/apikey", apiUser)
	log.Printf("Preparing API Key request for user %s", apiUser)
	log.Printf("Request URL: %s", url)

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request for API Key: %v", err)
		return "", err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)
	log.Println("Added required headers: Content-Type, Ocp-Apim-Subscription-Key")

	// Send the request
	log.Println("Sending API Key creation request to MTN MoMo API...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: HTTP request for API Key failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	log.Printf("Received API Key response with status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("ERROR: API Key creation failed with status: %d, body: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("failed to create API key: %s, status: %d", string(body), resp.StatusCode)
	}

	// Parse the response
	var result struct {
		APIKey string `json:"apiKey"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("ERROR: Failed to parse API Key response: %v", err)
		return "", err
	}

	log.Println("Successfully retrieved API Key from MTN MoMo API")
	// We don't log the actual API key for security reasons
	return result.APIKey, nil
}

// fallbackGenerateAPIKey creates an API key locally as a fallback
func fallbackGenerateAPIKey() string {
	// Generate a random API key (32 hex characters)
	randomBytes := make([]byte, 16) // 16 bytes will generate 32 hex characters
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Convert to hex string
	apiKey := fmt.Sprintf("%x", randomBytes)
	return apiKey
}

// fallbackGenerateAPIUser creates a unique API user (UUID) locally as a fallback
func fallbackGenerateAPIUser() string {
	// Generate a UUID for the API user as per MTN MoMo API documentation
	return uuid.New().String()
}

// handleGenerateKeys handles the key generation request
func handleGenerateKeys(w http.ResponseWriter, r *http.Request) {
	log.Println("=== New API Key Generation Request Received ===")

	var req MomoKeyRequest

	// Parse JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("ERROR: Invalid request format - %v", err)
		sendResponse(w, false, "Invalid request format", nil, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.PrimaryKey == "" {
		log.Println("ERROR: Missing required field - Subscription Key (Primary Key)")
		sendResponse(w, false, "Subscription Key (Primary Key) is required", nil, http.StatusBadRequest)
		return
	}

	// Default callback host if not provided
	callbackHost := req.CallbackHost
	if callbackHost == "" {
		log.Println("INFO: No callback host provided, using default: example.com")
		callbackHost = "example.com"
	} else {
		log.Printf("INFO: Using provided callback host: %s", callbackHost)
	}

	// Variables to store our API credentials
	var apiUser, apiKey string
	var useRealAPI bool = true

	log.Println("=== ATTEMPTING REAL MTN MOMO API INTEGRATION ===")
	if useRealAPI {
		// Try to use the real MTN MoMo API
		log.Println("STEP 1/2: Creating API User through MTN MoMo API...")

		// Step 1: Create API User through MTN MoMo API
		apiUserResult, err := createAPIUser(req.PrimaryKey, callbackHost)
		if err != nil {
			log.Printf("ERROR: Failed to create API User via MTN MoMo API - %v", err)
			log.Println("FALLBACK: Will use local generation instead")
			useRealAPI = false
		} else {
			apiUser = apiUserResult
			log.Printf("SUCCESS: API User created and registered with MTN MoMo: %s", apiUser)

			// Step 2: Create API Key through MTN MoMo API
			log.Println("STEP 2/2: Creating API Key through MTN MoMo API...")
			apiKeyResult, err := createAPIKey(req.PrimaryKey, apiUser)
			if err != nil {
				log.Printf("ERROR: Failed to create API Key via MTN MoMo API - %v", err)
				log.Println("FALLBACK: Will use local generation instead")
				useRealAPI = false
			} else {
				apiKey = apiKeyResult
				log.Printf("SUCCESS: API Key created and registered with MTN MoMo for user %s", apiUser)
				log.Println("=== MTN MOMO API INTEGRATION SUCCESSFUL ===")
			}
		}
	}

	// If real API failed, fall back to local generation
	if !useRealAPI {
		log.Println("=== USING LOCAL GENERATION (NOT REGISTERED WITH MTN MOMO) ===")
		log.Println("STEP 1/2: Generating API User locally...")
		apiUser = fallbackGenerateAPIUser()
		log.Printf("Generated API User locally: %s", apiUser)

		log.Println("STEP 2/2: Generating API Key locally...")
		apiKey = fallbackGenerateAPIKey()
		log.Printf("Generated API Key locally for user %s", apiUser)
		log.Println("=== LOCAL GENERATION COMPLETE ===")
		log.Println("WARNING: These credentials are NOT registered with MTN MoMo and cannot be used for API calls")
	}

	// Create response following MTN MoMo API structure
	resp := MomoKeyResponse{
		APIKey:       apiKey,
		APIUser:      apiUser,
		UserID:       apiUser, // In MTN MoMo, the API User is the same as the User ID (X-Reference-Id)
		CallbackHost: callbackHost,
		DateTime:     time.Now().Format(time.RFC3339),
		TargetEnv:    "sandbox", // Always sandbox in this simulator
	}

	// Generate Base64 auth string and test curl command for the user
	// Create the auth string (apiUser:apiKey) and encode it in base64
	authString := fmt.Sprintf("%s:%s", apiUser, apiKey)
	base64Auth := base64.StdEncoding.EncodeToString([]byte(authString))

	// Add the Base64 auth string to the response
	resp.Base64Auth = base64Auth

	// Generate the curl command if using real API
	if useRealAPI {
		// Generate the curl command
		testCommand := fmt.Sprintf("\nTest your credentials with this curl command:\n\ncurl --location --request POST 'https://sandbox.momodeveloper.mtn.com/collection/token/' \\\n--header 'Authorization: Basic %s' \\\n--header 'Ocp-Apim-Subscription-Key: %s' \\\n--header 'Content-Type: application/json'\n", base64Auth, req.PrimaryKey)

		log.Println("Generated test curl command for the user")
		log.Println(testCommand)

		// Add the test command to the response
		resp.TestCommand = testCommand
	}

	if useRealAPI {
		log.Println("Sending response with MTN MoMo registered credentials")
		sendResponse(w, true, "API User and API Key successfully created and registered with MTN MoMo", resp, http.StatusCreated)
	} else {
		log.Println("Sending response with locally generated credentials")
		sendResponse(w, true, "API User and API Key generated locally (not registered with MTN MoMo)", resp, http.StatusCreated)
	}

	log.Println("=== API Key Generation Request Completed ===")
}

// sendResponse sends a standardized JSON response
func sendResponse(w http.ResponseWriter, success bool, message string, data interface{}, statusCode int) {
	resp := Response{
		Success: success,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// setupLogger configures a more detailed logger
func setupLogger() {
	// Set log format to include timestamp
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Logger initialized with timestamp and file information")
}

func main() {
	// Setup enhanced logging
	setupLogger()
	log.Println("=== MTN MoMo API Key Generator Backend Starting ===")
	log.Println("This backend will attempt to register credentials with MTN MoMo API")
	log.Println("If MTN MoMo API is unavailable, it will fall back to local generation")

	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/api/generate", handleGenerateKeys).Methods("POST")
	log.Println("API route registered: POST /api/generate")

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	log.Println("CORS middleware configured to allow requests from http://localhost:3000")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
