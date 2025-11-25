# Build images
docker build -t cecilia08/api-gateway:v2 ./api-gateway
docker build -t cecilia08/ticket-service:v2 ./ticket-service
docker build -t cecilia08/waiting-room-service:v2 ./waiting-room-service
docker build -t cecilia08/dashboard-service:v2 ./dashboard-service
docker build -t cecilia08/website-dashboard:v2 ./website-dashboard

# Push images to Docker Hub
docker push cecilia08/api-gateway:v2
docker push cecilia08/ticket-service:v2
docker push cecilia08/waiting-room-service:v2
docker push cecilia08/dashboard-service:v2
docker push cecilia08/website-dashboard:v2
