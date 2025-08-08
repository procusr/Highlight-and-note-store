#!/bin/bash

# Test script for authentication
echo "Testing MH-Go Authentication..."

# Test 1: Register a new user
echo "1. Testing user registration..."
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=test@example.com&password=testpassword&confirm_password=testpassword" \
  -v

echo -e "\n\n"

# Test 2: Login with the registered user
echo "2. Testing user login..."
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=test@example.com&password=testpassword" \
  -v

echo -e "\n\n"

# Test 3: Login with wrong credentials (should fail)
echo "3. Testing login with wrong password..."
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=test@example.com&password=wrongpassword" \
  -v

echo "Authentication tests completed."
