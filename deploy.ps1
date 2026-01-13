# Build images
docker build -t cecilia08/api-gateway:latest ./api-gateway
docker build -t cecilia08/ticket-service:latest ./ticket-service
docker build -t cecilia08/waiting-room-service:latest ./waiting-room-service
docker build -t cecilia08/dashboard-service:latest ./dashboard-service
docker build -t cecilia08/website-dashboard:latest ./website-dashboard

# Push images to Docker Hub
docker push cecilia08/api-gateway:latest
docker push cecilia08/ticket-service:latest
docker push cecilia08/waiting-room-service:latest
docker push cecilia08/dashboard-service:latest
docker push cecilia08/website-dashboard:latestKs
