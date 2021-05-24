# NOTE: This is the simplest way of achieving a replicaset in mongodb with Docker.
# However if you would like a more automated approach, please see the setup.sh file and the docker-compose file which includes this startup script.

# run this after setting up the docker-compose This will instantiate the replica set.
# The id and hostname's can be tailored to your liking, however they MUST match the docker-compose file above.
docker-compose up -d
docker exec -it localmongo1 mongo

rs.initiate(
  {
    _id : 'rs0',
    members: [
      { _id : 0, host : "mongo1:27017" },
      { _id : 1, host : "mongo2:27017" },
      { _id : 2, host : "mongo3:27017", arbiterOnly: true }
    ]
  }
)

exit