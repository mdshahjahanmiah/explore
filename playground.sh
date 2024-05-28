#!/bin/bash

# Function to install jq if not already installed
install_jq() {
    echo "jq could not be found, installing jq..."
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get &> /dev/null; then
            sudo apt-get install -y jq
        elif command -v yum &> /dev/null; then
            sudo yum install -y jq
        else
            echo "Package manager not supported. Please install jq manually."
            exit 1
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        brew install jq
    else
        echo "OS not supported. Please install jq manually."
        exit 1
    fi
}

# Install jq if not already installed
if ! command -v jq &> /dev/null; then
    install_jq
fi

# Step 1: Call the ciphertext endpoint and save the response to response.json
curl -X GET http://localhost:9000/ds/ciphertext -H "Content-Type: application/json" -o response.json

# Step 2: Extract the ciphertext value from the response.json file
ciphertext=$(jq -r '.ciphertext' response.json)

# Print the ciphertext value
echo "Retrieved ciphertext: $ciphertext"

# Step 3: Prepare the JSON payload
payload=$(jq -n --arg ct "$ciphertext" '{"ciphertext": $ct}')

# Print the JSON payload for debugging
echo "JSON payload to send to decrypt endpoint: $payload"

# Step 4: Call the decrypt endpoint with the extracted ciphertext
response=$(curl -X POST http://localhost:9000/ds/decrypt -H "Content-Type: application/json" -d "$payload")

# Print the response from the decrypt endpoint
echo "Response from decrypt endpoint: $response"
