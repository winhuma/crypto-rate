#/bin/bash

cp .dockerignore service1
cp -r libs service1
cd service1
docker build -t crypto-get-rate .
rm -rf libs .dockerignore
cd ..

cp .dockerignore service2
cp -r libs service2
cd service2
docker build -t crypto-set-db .
rm -rf libs .dockerignore
cd ..

cp .dockerignore service3
cp -r libs service3
cd service3
docker build -t crypto-api .
rm -rf libs .dockerignore