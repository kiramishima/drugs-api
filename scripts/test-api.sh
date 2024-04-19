set -e

echo "GET http://localhost:8080/v1/credits/statistics" | vegeta attack -duration=5s -rate=200 | vegeta report --type=text

