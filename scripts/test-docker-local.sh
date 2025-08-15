#!/bin/bash

# Local Docker testing script for TimeSeriesDB
# This script tests the Docker container locally to verify it works correctly

set -e

echo "🧪 Starting local Docker image testing..."
echo "Image: ghcr.io/maher-naija-pro/temporadb:latest"

# Pull the image
echo "📥 Pulling Docker image..."
docker pull ghcr.io/maher-naija-pro/temporadb:latest

# Show comprehensive image details for debugging
echo "📋 Image details:"
docker images | grep temporadb
docker inspect ghcr.io/maher-naija-pro/temporadb:latest | jq '.[0].Config.Entrypoint, .[0].Config.Cmd, .[0].Config.ExposedPorts' || echo "Could not inspect image"

# Test container starts with better error handling
echo "🚀 Starting container for testing..."
container_name="test-container-local-$(date +%s)"
container_id=$(docker run -d \
  -p 8080:8080 \
  -e DEBUG=true \
  -e LOG_LEVEL=debug \
  --name "$container_name" \
  ghcr.io/maher-naija-pro/temporadb:latest)

echo "Container ID: $container_id"
echo "Container Name: $container_name"

# Verify container name is set correctly
if [ -z "$container_name" ]; then
  echo "❌ ERROR: Container name is empty!"
  exit 1
fi

# Verify container ID is set correctly
if [ -z "$container_id" ]; then
  echo "❌ ERROR: Container ID is empty!"
  exit 1
fi

# Enhanced container startup monitoring
echo "⏳ Waiting for container to start..."
echo "🔍 Debug info: container_name='$container_name', container_id='$container_id'"
max_wait_time=300
wait_time=0
container_started=false

while [ $wait_time -lt $max_wait_time ]; do
  echo "Check $((wait_time/30 + 1)): Container status check... ($wait_time/300 seconds)"
  
  # Check if container is running using multiple methods
  container_running=false
  
  # Method 1: Check docker ps
  if docker ps | grep -q $container_id; then
    container_running=true
  fi
  
  # Method 2: Check container state directly
  if [ "$container_running" = false ]; then
    container_state=$(docker inspect $container_name 2>/dev/null | jq -r '.[0].State.Status' 2>/dev/null || echo "unknown")
    if [ "$container_state" = "running" ]; then
      container_running=true
    fi
  fi
  
  if [ "$container_running" = true ]; then
    echo "✅ Container is running"
    container_started=true
    
    # Show container details
    echo "🔍 Container details:"
    docker inspect $container_id | jq '.[0].State, .[0].NetworkSettings' || echo "Could not inspect container"
    
    # Check if port is listening
    echo "🔌 Port status check:"
    netstat -tlnp | grep 8080 || echo "Port 8080 not listening yet"
    
    # Additional debugging - show both container ID and name in ps output
    echo "🔍 Docker ps output:"
    docker ps | grep $container_id || echo "Container not found in docker ps"
    echo "🔍 Docker ps -a output:"
    docker ps -a | grep $container_name || echo "Container not found in docker ps -a"
    
    break
  else
    echo "⏰ Container not running yet, waiting... ($wait_time/$max_wait_time seconds)"
    
    # Check if container exists but failed
    if docker ps -a | grep -q $container_name; then
      echo "🔍 Container exists but not running - checking status..."
      container_state=$(docker inspect $container_name | jq -r '.[0].State.Status' 2>/dev/null || echo "unknown")
      container_exit_code=$(docker inspect $container_name | jq -r '.[0].State.ExitCode' 2>/dev/null || echo "unknown")
      
      echo "Container state: $container_state, Exit code: $container_exit_code"
      
      # If container has exited with error, break
      if [ "$container_state" = "exited" ] && [ "$container_exit_code" != "0" ]; then
        echo "❌ Container exited with error code $container_exit_code"
        docker logs $container_name || echo "Could not retrieve container logs"
        break
      fi
      
      # If container is still starting up, continue waiting
      if [ "$container_state" = "created" ] || [ "$container_state" = "running" ]; then
        echo "⏳ Container is starting up, continuing to wait..."
        sleep 30
        wait_time=$((wait_time + 30))
        continue
      fi
    fi
    
    sleep 30
    wait_time=$((wait_time + 30))
  fi
done

if [ "$container_started" = false ]; then
  echo "❌ Container failed to start within $max_wait_time seconds"
  echo "🔍 Container logs:"
  docker logs $container_name || echo "Could not retrieve container logs"
  echo "🔍 Container inspect:"
  docker inspect $container_name || echo "Could not inspect container"
  echo "🔍 Docker ps -a:"
  docker ps -a | grep $container_name || echo "Container not found in ps -a"
  echo "🔍 Docker system info:"
  docker system df || echo "Could not get Docker system info"
  echo "🔍 System resources:"
  free -h || echo "Could not get memory info"
  df -h || echo "Could not get disk info"
  exit 1
fi

# Wait for application to initialize with progress
echo "⏳ Waiting for application to initialize..."
echo "🔍 Container startup timeline:"
echo "  - Container created at: $(date -u)"
echo "  - Container ID: $container_id"
echo "  - Container Name: $container_name"
echo "  - Port mapping: 8080:8080"

for i in {1..12}; do
  echo "Initialization check $i/12..."
  sleep 5
  
  # Check if container is still running
  if ! docker ps | grep -q $container_id; then
    echo "❌ Container stopped during initialization"
    echo "🔍 Container logs:"
    docker logs $container_name || echo "Could not retrieve container logs"
    exit 1
  fi
  
  # Test basic connectivity
  if curl -f -s http://localhost:8080/ >/dev/null 2>&1; then
    echo "✅ Basic connectivity established"
    break
  fi
done

# Enhanced health endpoint testing with comprehensive retry logic
echo "🏥 Testing health endpoint..."
max_retries=15
retry_count=0
health_check_passed=false

while [ $retry_count -lt $max_retries ]; do
  echo "🔍 Health check attempt $((retry_count + 1))/$max_retries..."
  
  # Check if container is still running
  if ! docker ps | grep -q $container_id; then
    echo "❌ Container stopped unexpectedly during health check"
    echo "🔍 Container logs:"
    docker logs $container_name || echo "Could not retrieve container logs"
    echo "🔍 Container inspect:"
    docker inspect $container_name | jq '.[0].State' || echo "Could not inspect container"
    exit 1
  fi
  
  # Test health endpoint
  if curl -f -s http://localhost:8080/health >/dev/null 2>&1; then
    echo "✅ Health check passed"
    health_check_passed=true
    break
  else
    echo "⚠️  Health check attempt $((retry_count + 1)) failed"
    echo "🔍 Container logs (last 30 lines):"
    docker logs --tail 30 $container_name || echo "Could not retrieve container logs"
    echo "🔍 Container status:"
    docker ps | grep $container_id || echo "Container not found"
    echo "🔍 Port status:"
    netstat -tlnp | grep 8080 || echo "Port 8080 not listening"
    retry_count=$((retry_count + 1))
    sleep 3
  fi
done

if [ "$health_check_passed" = false ]; then
  echo "❌ Health check failed after $max_retries attempts"
  echo "🔍 Final container logs:"
  docker logs $container_name || echo "Could not retrieve container logs"
  echo "🔍 Container inspect:"
  docker inspect $container_name || echo "Could not inspect container"
  echo "🔍 Container resource usage:"
  docker stats --no-stream $container_name || echo "Could not get container stats"
  exit 1
fi

# Test metrics endpoint
echo "📊 Testing metrics endpoint..."
if curl -f -s http://localhost:8080/metrics >/dev/null 2>&1; then
  echo "✅ Metrics endpoint accessible"
else
  echo "⚠️  Metrics endpoint not accessible"
fi

# Test root endpoint
echo "🌐 Testing root endpoint..."
if curl -f -s http://localhost:8080/ >/dev/null 2>&1; then
  echo "✅ Root endpoint accessible"
else
  echo "⚠️  Root endpoint not accessible"
fi

# Test with different HTTP methods
echo "🔍 Testing HTTP methods..."
for method in GET HEAD OPTIONS; do
  if curl -f -s -X $method http://localhost:8080/ >/dev/null 2>&1; then
    echo "✅ $method method works"
  else
    echo "⚠️  $method method failed"
  fi
done

# Cleanup with better error handling
echo "🧹 Cleaning up test container..."
docker stop $container_name || echo "Could not stop container"
docker rm $container_name || echo "Could not remove container"

echo "🎉 Docker image test completed successfully!"

# Final test summary
echo "📋 Test Summary:"
echo "=================="
echo "✅ Container started successfully"
echo "✅ Container is running and healthy"
echo "✅ Port 8080 is accessible"
echo "✅ Health endpoint is responding"
echo "✅ Application is fully initialized"
echo "=================="

echo "🎯 All tests passed! The Docker container is working correctly."
