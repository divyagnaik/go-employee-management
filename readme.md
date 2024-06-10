go-employee-management

# Framework - Gofr (Links - https://github.com/gofr-dev/gofr & https://gofr.dev/)

# Migrations
    $ migrate create -ext=sql -dir=migrations -seq init
# Create a postgres docker container for employee
    $ docker run --name employee -e POSTGRES_DB=employee -e POSTGRES_PASSWORD=root123 -p 2007:5432 -d postgres:latest

# To get inside the docker container
    $ docker exec -it employee psql -U postgres

# Connect to the database
    postgres=#  \c employee

# To check the relations present in database, after running the service with successful migration
    postgres=# \dt