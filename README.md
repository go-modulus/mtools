# Mtools
CLI application to manage a project built over [Modulus framework](https://github.com/go-modulus/modulus).

## Usage
Mtools is a CLI application. Install it using `go install github.com/go-modulus/mtools/cmd/mtools@latest`

It allows:
* init a project `mtools init`
* install modules `mtools module install` in the project directory
* create a new module `mtools module create`
* add a PostgreSQL migration `mtools db add`
* run migrations `mtools db migrate`
* update SQLs config of all modules from templates defined in the project `mtools db update-sqlc-config`
* add cli command into module `mtools module add-cli`
* add REST API endpoint into module `mtools module add-json-api`


All these mtools commands except `mtools init` are available inside the projet under makefile commands. 
Use `make help` to see them.
