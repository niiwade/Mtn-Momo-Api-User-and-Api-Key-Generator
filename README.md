# MTN MoMo API Key Generator

A web application for generating MTN Mobile Money API Keys and API Users. This project uses Next.js for the frontend and Go for the backend.

> **Important**: This application now integrates with the actual MTN MoMo Sandbox API to register your API Users and API Keys. When you provide a valid subscription key, the application will attempt to create real credentials that can be used with the MTN MoMo API. If the API call fails, it will fall back to generating credentials locally.

## Features

- **Real API Integration**: Create and register API Users and API Keys with the actual MTN MoMo Sandbox API
- Automatic fallback to local generation if the MTN MoMo API is unavailable
- Specify custom callback hosts for your application
- Visual indication of whether credentials are registered with MTN MoMo or generated locally
- **Base64 Encoded Auth String**: Automatically generate the Base64 encoded authentication string required for API calls
- **Test Command Generation**: Provides a ready-to-use cURL command for testing your credentials
- Copy generated credentials to clipboard with one click
- Modern, responsive UI with Tailwind CSS and Next.js
- Secure backend implementation in Go

## Project Structure

- `momo-key-generator/` - Next.js frontend application
- `backend/` - Go backend service

## Getting Started

### Prerequisites

- Node.js (v14 or later)
- Go (v1.16 or later)

### Running the Backend

1. Navigate to the backend directory:
   ```
   cd backend
   ```

2. Run the Go server:
   ```
   go run main.go
   ```

   The server will start on port 8080.

### Running the Frontend

1. Navigate to the frontend directory:
   ```
   cd momo-key-generator
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm run dev
   ```

   The frontend will be available at http://localhost:3000.

## How to Use

### Starting the Application

1. **Start the backend server**:
   ```
   cd backend
   go run main.go
   ```
   The server will start on port 8080.

2. **Start the frontend server**:
   ```
   cd momo-key-generator
   npm run dev
   ```
   The frontend will be available at http://localhost:3000.

### Generating API Credentials

1. Open the application in your browser at http://localhost:3000

2. Fill in the form with the following information:
   - **Subscription Key (Primary Key)**: Your MTN MoMo subscription key (required)
   - **Secondary Key**: Optional secondary key
   - **Provider Callback Host**: Your application's callback host (e.g., example.com)

3. Click the "Generate API User & Key" button

4. The application will display:
   - **API Key**: Use this for authentication with MTN MoMo API
   - **API User (X-Reference-Id)**: Use this as your API User ID in API calls
   - **Callback Host**: Your registered callback host
   - **Target Environment**: Always "sandbox" in this simulator
   - **Base64 Encoded Auth String**: Pre-generated Base64 encoded string of `apiUser:apiKey` for use in the Authorization header
   - **Test Command**: A ready-to-use cURL command for testing your credentials against the MTN MoMo API (only shown for credentials registered with MTN MoMo)

5. Use the "Copy" buttons to copy these values to your clipboard

6. **Testing Your Credentials**:
   - If your credentials were registered with MTN MoMo, use the provided test command to verify they work
   - The test command will attempt to get an access token from the MTN MoMo API
   - A successful response will include an access token and token type

7. **Important**: Store your credentials securely. The API Key cannot be retrieved again if lost.

## API Endpoints

### Generate API User and API Key

- **URL**: `/api/generate`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "primaryKey": "your-subscription-key",
    "secondaryKey": "your-secondary-key",
    "callbackHost": "example.com"
  }
  ```
  Note: `secondaryKey` and `callbackHost` are optional. If `callbackHost` is not provided, it defaults to "example.com".

- **Response**:
  ```json
  {
    "success": true,
    "message": "API User and API Key successfully created and registered with MTN MoMo",
    "data": {
      "apiKey": "generated-api-key",
      "apiUser": "generated-api-user",
      "userId": "generated-api-user",
      "callbackHost": "example.com",
      "targetEnvironment": "sandbox",
      "dateTime": "2025-07-08T16:51:32Z",
      "base64Auth": "base64-encoded-auth-string",
      "testCommand": "curl command for testing credentials"
    }
  }
  ```
  
  Note: The `message` field will indicate whether credentials were registered with MTN MoMo or generated locally. The `testCommand` field is only included when credentials are successfully registered with MTN MoMo.

## License

This project is licensed under the MIT License.