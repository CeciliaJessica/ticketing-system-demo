# Build images
docker build -t cecilia08/api-gateway:v7 ./api-gateway
docker build -t cecilia08/ticket-service:v14 ./ticket-service
docker build -t cecilia08/waiting-room-service:v7 ./waiting-room-service
docker build -t cecilia08/dashboard-service:v7 ./dashboard-service
docker build -t cecilia08/website-dashboard:v7 ./website-dashboard

# Push images to Docker Hub
docker push cecilia08/api-gateway:v7
docker push cecilia08/ticket-service:v14
docker push cecilia08/waiting-room-service:v7
docker push cecilia08/dashboard-service:v7
docker push cecilia08/website-dashboard:v7
