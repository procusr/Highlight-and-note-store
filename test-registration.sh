#!/bin/bash

echo "🧪 Testing Registration Functionality"
echo "===================================="

# Test successful registration
echo "Testing successful registration..."
response=$(curl -s -w "%{http_code}" -X POST http://localhost:8080/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=test$(date +%s)@example.com&password=password123&confirm_password=password123")

http_code="${response: -3}"
if [ "$http_code" = "200" ]; then
    echo "✅ Registration successful"
else
    echo "❌ Registration failed with code: $http_code"
fi

# Test password mismatch
echo "Testing password mismatch..."
response=$(curl -s -w "%{http_code}" -X POST http://localhost:8080/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=test2@example.com&password=password123&confirm_password=different123")

http_code="${response: -3}"
if [ "$http_code" = "400" ]; then
    echo "✅ Password mismatch properly handled"
else
    echo "❌ Password mismatch not handled correctly: $http_code"
fi

# Test empty fields
echo "Testing empty fields..."
response=$(curl -s -w "%{http_code}" -X POST http://localhost:8080/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=&password=&confirm_password=")

http_code="${response: -3}"
if [ "$http_code" = "400" ]; then
    echo "✅ Empty fields properly handled"
else
    echo "❌ Empty fields not handled correctly: $http_code"
fi

echo ""
echo "🎉 Registration testing complete!"
echo "You can now test the web interface at: http://localhost:8080"
