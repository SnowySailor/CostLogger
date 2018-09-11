# Cost Logger
### Requirements:
1. Golang
2. PostgreSQL
3. Redis
4. Python 3 (optional)

### Running
To run this application, you'll need to do some brief setup.
1. Run the PostgreSQL script to create an empty database (called `cost_logger` by default).
    1. `$ psql -d postgres -U yourpostgresuser -a -f schema.sql`
2. Create a `secrets.yaml` file in the root of the repository and base your settings off of `secrets.yaml.example`. 

After these steps, run the `runner.py` script with `$ python3.6 runner.py` to build and start the application. If there are any errors, they'll be dumped to your terminal. If you make any changes to any `.template` files or `.go` files, the python script will detect that and automatically rebuild the website and restart it.

### Testing
To run the tests (that need to have more coverage): `$ cd site; go test`
