trap ctrl_c INT

function ctrl_c() {
  docker-compose stop
}

docker-compose build
docker-compose run start_dependencies
docker-compose up
