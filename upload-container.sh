docker build -t go-itudy .          
docker tag go-itudy southamerica-west1-docker.pkg.dev/itudy-485221/itudy/itudy:latest
docker push southamerica-west1-docker.pkg.dev/itudy-485221/itudy/itudy:latest

