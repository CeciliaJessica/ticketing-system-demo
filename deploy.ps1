# Build images
docker build -t cecilia08/api-gateway:v1 ./api-gateway
docker build -t cecilia08/ticket-service:v1 ./ticket-service
docker build -t cecilia08/waiting-room-service:v1 ./waiting-room-service
docker build -t cecilia08/dashboard-service:v1 ./dashboard-service
docker build -t cecilia08/website-dashboard:v1 ./website-dashboard

# Push images to Docker Hub
docker push cecilia08/api-gateway:v1
docker push cecilia08/ticket-service:v1
docker push cecilia08/waiting-room-service:v1
docker push cecilia08/dashboard-service:v1
docker push cecilia08/website-dashboard:v1
