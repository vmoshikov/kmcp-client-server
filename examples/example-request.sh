#!/bin/bash

# Пример запроса для анализа подов
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "clusterId": "example-cluster",
    "toolNames": ["kubectl_get_pods"],
    "parameters": {
      "namespace": "default"
    }
  }' | jq .
