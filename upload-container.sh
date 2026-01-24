docker build -t change .          
docker tag go-change southamerica-east1-docker.pkg.dev/easylife-464420/microservices/change:latest

docker push southamerica-east1-docker.pkg.dev/easylife-464420/microservices/change:latest

