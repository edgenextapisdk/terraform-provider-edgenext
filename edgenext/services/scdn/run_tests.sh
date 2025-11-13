#!/bin/bash

# SCDN Service Test Runner
# This script demonstrates how to run the SCDN service tests

echo "=== SCDN Service Test Runner ==="
echo

# Check if we're in the right directory
if [ ! -f "service_en_scdn_test.go" ]; then
    echo "Error: Please run this script from the scdn service directory"
    exit 1
fi

echo "1. Running unit tests (no API credentials required)..."
go test -v -run "TestNewScdnService|TestDomainInfo|TestOrigin|TestOriginRecord|TestCnameInfo|TestStatus|TestDomainListData" .

echo
echo "2. Running integration tests (requires config file)..."
echo "   Create test_config.json file with your API credentials:"
echo "   {"
echo "     \"access_key\": \"your_access_key\","
echo "     \"secret_key\": \"your_secret_key\","
echo "     \"endpoint\": \"https://api.edgenextscdn.com\","
echo "     \"region\": \"us-east-1\","
echo "     \"timeout_seconds\": 30,"
echo "     \"enable_integration_tests\": true"
echo "   }"
echo "   Copy test_config.json.example to test_config.json and update with your credentials"
echo

if [ -f "test_config.json" ]; then
    echo "   Config file found, running integration tests..."
    go test -v -run "Integration" .
else
    echo "   No config file found, skipping integration tests."
    echo "   To run integration tests, create test_config.json file with your credentials:"
    echo "   cp test_config.json.example test_config.json"
    echo "   # Edit test_config.json with your API credentials"
    echo "   go test -v -run 'Integration' ."
fi

echo
echo "3. Running all tests (unit + integration if credentials available)..."
go test -v .

echo
echo "=== Test run complete ==="