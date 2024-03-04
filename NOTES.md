## Coding Challenge Guidelines - Developer Notes

Context: Initially I understood that by Kratos and Wire, you were refering to:
- [Kratos](https://www.ory.sh/docs/welcome)
- [Wire](https://docs.wire.com/)

Therefore, I started to implement the project in a vanilla way: without any framework.
Later, I received feedback that Kratos and Wire were referring to:
- [Kratos](https://go-kratos.dev/)
- [Wire](https://pkg.go.dev/github.com/google/wire)

So, this solution has been implemented using both aproaches, vanilla and with Kratos and Wire.

The vanilla solution will be found on the vanilla folder.
The Kratos and Wire solution will be found on the kratos folder.

### Vanilla Approach
- Project skeleton copied from my own [repositories](https://github.com/guidomantilla?tab=repositories&q=go-feather-&type=&language=go&sort=)
- cmd folder contains serve command for staring the server, and migrate command for database migrations
- core folder contains the business logic
- pkg folder contains reusable code (NOTE: move to higher level)
- docker folder contains the docker-compose file for running the database and the vanilla version of this project
- resources folder will contain the database migrations, proto files and postman collection
- tools folder contains the wire code generation for the vanilla version of this project

To run the vanilla version of this project:
```bash
cd vanilla
make build
```
These commands will build the project, create a docker project with the application, migration and database apps. Check vanilla/docker/docker-compose.yml for database connection details.
* If the database have no tables created yet, please keep restarting migration-app until it creates the tables.

#### Vanilla Design
As mentioned earlier, pkg folder contains reusable code. The layout of the pkg folder is as follows:
* boot: contains the dependency injection code used for wiring the application
* config: contains the configuration coming from environment variables
* database: contains the database connection and migration code using GORM
* environment: contains the code used to represent the environment variables
* errors: common utility code for error handling
* log: contains the code for logging using log/slog
* mocks: contains utility code for mocking
* properties: contains the code for managing properties. These may come from different sources (e.g. environment variables, configuration files, etc.)
* rest: contains utility code for REST API
* security: contains the code for security (e.g. JWT, etc.)
* util: contains utility code
* validation: contains the code for validation

On core folder, the layout follows a simplistic 3 layer design, as follows:
* endpoints: contains the code for the REST API & Grpc endpoints. No business logic should be here.
* models: contains the code for the domain models
* services: contains the code for the business logic
Note: DAO (Data Access Object), Repository Pattern, or any other DAL (Data Access Layer) is not used in this project. The database logic is directly in the services using GORM


### Kratos Approach
- Project skeleton based from kratos CLI
- pkg folder reuse as much as possible 
- wire was used for wiring pkg dependencies


#### Kratos Design
As mentioned earlier, pkg folder contains reusable code.
On core folder, the layout follows a simplistic 3 layer design, as follows:
* internal/server: contains the code for the REST API & Grpc endpoints. No business logic should be here.
* internal/facade: contains code that proxies services layer
* internal/models: contains the code for the domain models
* internal/services: contains the code for the business logic
  Note: DAO (Data Access Object), Repository Pattern, or any other DAL (Data Access Layer) is not used in this project. The database logic is directly in the services using GORM



### Evaluation Criteria
- Go best practices:
  - This project follows some Uber's Best Practices, some legacy Java/SpringBoot best practices, while avoid as much as possible falling into over designing.
  - This project uses golang linting and formatting tools to ensure the code is clean and consistent following the go best practices.
  - This project uses sonarqube to ensure all technical debt is being tracked and reduced. Please check [here](http://170.187.157.212:9001/code?id=golang-engineer-udoly).
 
- Completeness: did you complete the features?
  - I think all requested features are implemented using both REST and GRPC endpoints.

- Clean code: is the code clean, easy to read, and maintainable?
  - Yes. The code passes the tests, reveal its intention, avoid duplication, and utilize the fewest elements possible; while keeping the base code robust enough to be extended and maintained.
  - The design principles used there are KISS, DRY, YAGNI, SOLID, Dependency Injection.
  - The Unit Test uses DATA-DOG/go-sqlmock, go.uber.org/mock, stretchr/testify
  - The code was designed and coded to be as easy-tested as possible in order to achieve a 100% coverage. I have the belief that a 100% coverage is consequence of having a good design and a good test strategy, not the other way around.

- Correctness: does the functionality act in sensible, thought-out ways?
  - The code supports both REST and GRPC endpoints, and the functionality is well tested.
  - The code supports explicit transaction management with the database.
  
- Maintainability: is it written in a clean, maintainable way?
  - As mentioned earlier the code was designed and coded to be as easy-tested as possible in order to achieve a 100% coverage. 
  - In order to achieve a 100% coverage, the following design principles were used: KISS, DRY, YAGNI, SOLID, Dependency Injection.
  - I have the belief that a 100% coverage is consequence of having a good design and a good test strategy, not the other way around.

- Testing: is the system adequately tested?
  - Yes and No. Yes it is adequately tested. With this I mean, the code is easy-testeable.
  - No. Because at this moment there is no 100% coverage, CI/CD pipeline, or any other automated testing strategy in place.
  
- Documentation: is the API well-documented?
  - There is no swagger documentation. But there is a postman collection in the resources folder.