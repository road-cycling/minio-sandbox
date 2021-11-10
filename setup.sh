docker build -t minio -f Dockerfile-minio .
docker build -t entry -f Dockerfile-entrypoint .

curl --output mc https://dl.min.io/client/mc/linux-amd64/mc && chmod +x mc
# ./mc alias set ENV http://localhost:9001 minioaccesskey miniosecretkey
