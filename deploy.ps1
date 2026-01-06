# Build images
docker build -t cecilia08/api-gateway:v6 ./api-gateway
docker build -t cecilia08/ticket-service:v13 ./ticket-service
docker build -t cecilia08/waiting-room-service:v6 ./waiting-room-service
docker build -t cecilia08/dashboard-service:v5 ./dashboard-service
docker build -t cecilia08/website-dashboard:v5 ./website-dashboard

# Push images to Docker Hub
docker push cecilia08/api-gateway:v6
docker push cecilia08/ticket-service:v13
docker push cecilia08/waiting-room-service:v6
docker push cecilia08/dashboard-service:v5
docker push cecilia08/website-dashboard:v5
