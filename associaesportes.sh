#!/bin/bash

# Step 1: Login and store the token in a variable.
# The `jq` utility is used here to parse the JSON response and extract the token.
# You might need to install it (`sudo apt-get install jq` or `brew install jq`).
echo "Attempting to log in..."
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email": "jogador2@email.com","password": "Senha123!"}' | jq -r .token)

# Check if the token was successfully retrieved.
if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo "Login failed. Please check your credentials and ensure the server is running."
  exit 1
fi

echo "Login successful. Token received."

# Step 2: Make the authenticated request to associate sports.
# - Replace '1' with the actual ID of the user you want to modify.
# - The Authorization header now includes the token.
# - The JSON payload uses the correct "esporte_ids" key.
echo "Associating sports with user ID 1..."
curl -i -X POST http://localhost:8080/api/usuarios/2/esportes \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"esporte_ids": [8, 6]}'

