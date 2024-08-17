![infra](infra_p02.png)

# Serverless Api-Scraper

Serverless API. Scrape data from a supermarket website and store it in a dynamoDB database.

Its a mono-repo with the following structure:

- `api` - Serverless API - Products
- `scraper` - Serverless Scraper
- `shared` - Shared code

The infrastructure is defined using Terraform.

The api lambda is triggered by an API Gateway and the scraper lambda is triggered by the api.

- `[GET] /api/v1/products` - Get all products

```json
{
    "code": 200,
    "status": "OK",
    "message": "Success getting all products",
    "data": [
        {
            "product_id": "uuid",
            "name": "Producto 1",
            "category": "category 1",
            "original_price": 899,
            "discounted_price": 0
        },
        {
            "product_id": "uuid",
            "name": "Producto 2",
            "category": "category 2",
            "original_price": 999,
            "discounted_price": 0
        }
    ]
}
```

- `[GET] /api/v1/products/{ProductID}` - Get a product by ID

```json
{
    "code": 200,
    "status": "OK",
    "message": "Success getting product",
    "data": {
        "product_id": "uuid",
        "name": "Producto 1",
        "category": "category 1",
        "original_price": 899,
        "discounted_price": 0
    }
}
```

- `[POST] /api/v1/products` - Update Data needs a token

```json
{
    "update_data" : true
}
```

```json
{
    "code": 200,
    "status": "OK",
    "message": "Success updating data",
    "data": null
}
```

- `[POST] /api/v1/users` - Register a user

```json
{
    "username": "username",
    "email": "test@tes.com",
    "password": "password",
    "role": "admin"
}
```

```json
{
    "code": 200,
    "status": "OK",
    "message": "Success registering user",
    "data": "User username created successfully ID: uuid"
}
```

- `[POST] /api/v1/users/login` - Login a user

```json
{
    "email": "test@test.com",
    "password": "password"
}
```

```json
{
    "code": 200,
    "status": "OK",
    "message": "Success logging in",
    "data": "token"
}
```

- `[GET] /api/v1/users` - Get all users

```json
{
    "code": 200,
    "status": "OK",
    "message": "Success getting all users",
    "data": [
        {
            "user_id": "uuid",
            "username": "username",
            "email": "test@test.com",
        },
        {
            "user_id": "uuid",
            "username": "username",
            "email": "test1@test.com",
        }
    ]
}
```

## Class Diagram - API Products

```mermaid
classDiagram
    direction TB

    %% Interfaces en la parte superior
    class ProductRepository {
        <<interface>>
        +GetAll() []Product
        +GetByID(id: string) Product
    }

    class ProductService {
        <<interface>>
        +GetAll() []ProductResponse
        +GetByID(productID: string) ProductResponse
        +UpdateData(updateData: UpdateDataRequest) bool
    }

    class ProductController {
        <<interface>>
        +GetAll(ctx: *gin.Context)
        +GetByID(ctx: *gin.Context)
        +UpdateData(ctx: *gin.Context)
    }

    %% Implementaciones en el medio
    class ProductRepositoryImpl {
        -dynamodbiface.DynamoDBAPI db
        -string tableName
        +GetAll() []Product
        +GetByID(id: string) Product
    }

    class ProductServiceImpl {
        -ProductRepository productRepository
        +GetAll() []ProductResponse
        +GetByID(productID: string) ProductResponse
        +UpdateData(updateData: UpdateDataRequest) bool
    }

    class ProductControllerImpl {
        -ProductService productService
        +GetAll(ctx: *gin.Context)
        +GetByID(ctx: *gin.Context)
        +UpdateData(ctx: *gin.Context)
    }

    %% Clases relacionadas con productos y respuestas en la parte inferior
    class Product {
        +string ProductID
        +string Name
        +string Category
        +int OriginalPrice
        +int DiscountedPrice
    }

    class UpdateDataRequest {
        +bool UpdateData
    }

    class ProductResponse {
        +string ProductID
        +string Name
        +string Category
        +int OriginalPrice
        +int DiscountedPrice
    }

    class BaseResponse {
        +int Code
        +string Status
        +string Message
        +interface Data
    }

    %% Implementación de Interfaces
    ProductRepositoryImpl ..|> ProductRepository : implements
    ProductServiceImpl ..|> ProductService : implements
    ProductControllerImpl ..|> ProductController : implements

    %% Relaciones entre clases
    ProductResponse <|-- BaseResponse : data
    ProductRepositoryImpl o-- Product : manages
    ProductServiceImpl o-- ProductResponse : returns
    ProductServiceImpl o-- UpdateDataRequest : uses
    ProductControllerImpl o-- BaseResponse : returns


```

## Class Diagram - Scraper

```mermaid
classDiagram
    direction TB

    %% Interfaces
    class ScraperRepository {
        <<interface>>
        +Create(product models.Product) (models.Product, error)
        +DeleteAll() error
    }

    class ScraperService {
        <<interface>>
        +GetProducts() (bool, error)
    }

    class Scraper {
        <<interface>>
        +CleanPrice(price string) (int, error)
        +ScrapeData(baseURL string, maxPage int, category string) ([]models.Product, error)
    }

    %% Implementaciones
    class ScraperRepositoryImpl {
        -dynamodbiface.DynamoDBAPI db
        -string tableName
        +Create(product models.Product) (models.Product, error)
        +DeleteAll() error
    }

    class ScraperServiceImpl {
        -Scraper scraper.Scraper
        -ScraperRepository scraperRepository
        +GetProducts() (bool, error)
    }

    class ScraperImpl {
        -colly.Collector Collector
        +CleanPrice(price string) (int, error)
        +ScrapeData(baseURL string, maxPage int, category string) ([]models.Product, error)
    }

    %% Clases relacionadas
    class Product {
        +string ProductID
        +string Name
        +string Category
        +int OriginalPrice
        +int DiscountedPrice
    }

    %% Implementación de Interfaces
    ScraperRepositoryImpl ..|> ScraperRepository : implements
    ScraperServiceImpl ..|> ScraperService : implements
    ScraperImpl ..|> Scraper : implements

    %% Relaciones entre clases
    ScraperServiceImpl --> ScraperImpl : scraper
    ScraperServiceImpl --> ScraperRepositoryImpl : scraperRepository
    ScraperRepositoryImpl --> Product : manages
    ScraperImpl --> Product : returns

```

## Class Diagram - API Users

```mermaid
classDiagram
    direction TB

    %% Interfaces
    class UserRepository {
        <<interface>>
        +GetAll() ([]models.User, error)
        +GetByID(id string) (models.User, error)
        +Create(user models.User) (models.User, error)
        +GetByEmail(email string) (models.User, error)
    }

    class UserService {
        <<interface>>
        +RegisterUser(createUserReq request.CreateUserRequest) (models.User, error)
        +GetAllUsers() ([]response.UserResponse, error)
        +GetUserByID(id string) (response.UserResponse, error)
        +LogInUser(logInUserReq request.LogInUserRequest) (response.LogInUserResponse, error)
    }

    class UserController {
        <<interface>>
        +RegisterUser(c *gin.Context)
        +GetAllUsers(c *gin.Context)
        +GetUserByID(c *gin.Context)
        +LogInUser(c *gin.Context)
    }

    %% Implementaciones
    class UserRepositoryImpl {
        -dynamodbiface.DynamoDBAPI db
        -string tableName
        +GetAll() ([]models.User, error)
        +GetByID(id string) (models.User, error)
        +Create(user models.User) (models.User, error)
        +GetByEmail(email string) (models.User, error)
    }

    class UserServiceImpl {
        -UserRepository userRepository
        -*validator.Validate validator
        -utils.PasswordHasher passwordHasher
        -utils.JWTUtils jwtUtils
        +RegisterUser(createUserReq request.CreateUserRequest) (models.User, error)
        +GetAllUsers() ([]response.UserResponse, error)
        +GetUserByID(id string) (response.UserResponse, error)
        +LogInUser(logInUserReq request.LogInUserRequest) (response.LogInUserResponse, error)
    }

    class UserControllerImpl {
        -services.UserService userService
        +RegisterUser(c *gin.Context)
        +GetAllUsers(c *gin.Context)
        +GetUserByID(c *gin.Context)
        +LogInUser(c *gin.Context)
    }

    %% Clases relacionadas
    class User {
        +string UserID
        +string Username
        +string Email
        +string Password
        +string Role
    }

    class CreateUserRequest {
        +string Username
        +string Email
        +string Password
        +string Role
    }

    class LogInUserRequest {
        +string Email
        +string Password
    }

    class UserResponse {
        +string UserID
        +string Username
        +string Email
    }

    class LogInUserResponse {
        +string Token
    }

    class BaseResponse {
        +int Code
        +string Status
        +string Message
        +interface Data
    }

    %% Implementación de Interfaces
    UserRepositoryImpl ..|> UserRepository : implements
    UserServiceImpl ..|> UserService : implements
    UserControllerImpl ..|> UserController : implements

    %% Relaciones entre clases
    UserResponse <|-- BaseResponse : data
    LogInUserResponse <|-- BaseResponse : data
    UserRepositoryImpl o-- User : manages
    UserServiceImpl o-- UserResponse : returns
    UserServiceImpl o-- CreateUserRequest : uses
    UserServiceImpl o-- LogInUserRequest : uses
    UserControllerImpl o-- BaseResponse : returns
```

## Class Diagram - Authorizer

```mermaid
classDiagram
    direction TB

    %% Interfaces
    class JWTValidator {
        <<interface>>
        +ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error)
    }

    class Policy {
        <<interface>>
        +GeneratePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse
    }

    class AuthorizerHandler {
        <<interface>>
        +HandleAuthorizer(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)
    }

    %% Implementaciones
    class JWTValidatorImpl {
        +ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error)
    }

    class PolicyImpl {
        +GeneratePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse
    }

    class AuthorizerHandlerImpl {
        -jwtValidator auth.JWTValidator
        -policy aws.Policy
        +HandleAuthorizer(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error)
    }

    %% Implementación de Interfaces
    JWTValidatorImpl ..|> JWTValidator : implements
    PolicyImpl ..|> Policy : implements
    AuthorizerHandlerImpl ..|> AuthorizerHandler : implements

    %% Relaciones entre clases
    AuthorizerHandlerImpl --> JWTValidatorImpl : uses
    AuthorizerHandlerImpl --> PolicyImpl : uses
```

## Interaction Diagram between Lambdas and DynamoDB

```mermaid
sequenceDiagram
    participant User
    participant APIGateway
    participant AuthorizerLambda as authorizer Lambda
    participant UserLambda as api_users Lambda
    participant APILambda as api_products Lambda
    participant ScraperLambda as Scraper Lambda
    participant DynamoDB


    %% Register User request - response
    User->>APIGateway: Request (Register) [POST] /api/v1/users
    APIGateway->>UserLambda: Invoke Lambda
    UserLambda->>DynamoDB: Register User
    DynamoDB-->>UserLambda: Return Success
    UserLambda-->>APIGateway: Respond with Success
    APIGateway-->>User: User Created successfully

    %% Login User request - response
    User->>APIGateway: Request (Login) [POST] /api/v1/users/login
    APIGateway->>UserLambda: Invoke Lambda
    UserLambda->>DynamoDB: GetUserByEmail(email)
    DynamoDB-->>UserLambda: Return User
    UserLambda-->>APIGateway: Respond with Token
    APIGateway-->>User: Return Token

    %% Update Product List request - response
    User->>APIGateway: Request (Update Product List) [POST] /api/v1/products : Token
    APIGateway->>AuthorizerLambda: Validate Token
    AuthorizerLambda->>APIGateway: Return Success
    APIGateway->>APILambda: Invoke Lambda
    APILambda->>ScraperLambda: Trigger Scraper Lambda
    ScraperLambda->>SupermarketWebpage: Scrape Data
    ScraperLambda->>DynamoDB: Update Products
    AuthorizerLambda-->>APIGateway: Return Success
    ScraperLambda-->>APILambda: Return Success
    APILambda-->>APIGateway: Respond with Success
    APIGateway-->>User: Return Scraper Started Successfully
    
    %% Get All Products request - response
    User->>APIGateway: Request (Get all products) [GET] /api/v1/products
    APIGateway->>APILambda: Invoke Lambda
    APILambda->>DynamoDB: Query Products
    DynamoDB-->>APILambda: Return Products
    APILambda-->>APIGateway: Respond with Products
    APIGateway-->>User: List of Products

    %% Get Product by ID request - response
    User->>APIGateway: Request (Get product by ID) [GET] /api/v1/products/{ProductID}
    APIGateway->>APILambda: Invoke Lambda
    APILambda->>DynamoDB: Get Product by ID
    DynamoDB-->>APILambda: Return Product
    APILambda-->>APIGateway: Respond with Product
    APIGateway-->>User: Return Product

    %% Get All Users request - response
    User->>APIGateway: Request (Get all users) [GET] /api/v1/users
    APIGateway->>UserLambda: Invoke Lambda
    UserLambda->>DynamoDB: Get All Users
    DynamoDB-->>UserLambda: Return Users
    UserLambda-->>APIGateway: Respond with Users
    APIGateway-->>User: Return Users

    %% Get User by ID request - response
    User->>APIGateway: Request (Get user by ID) [GET] /api/v1/users/{UserID}
    APIGateway->>UserLambda: Invoke Lambda
    UserLambda->>DynamoDB: Get User by ID
    DynamoDB-->>UserLambda: Return User
    UserLambda-->>APIGateway: Respond with User
    APIGateway-->>User: Return User
```

## CI/CD Pipeline

```mermaid
graph TD
    A[Push to main branch] -->|Trigger| B[CI/CD Pipeline]
    B --> C[Test and Build API]
    B --> F[Test and Build API Users]
    
    subgraph "Test and Build API"
    C --> C1[Checkout code]
    C1 --> C2[Set up Go]
    C2 --> C3[Sync modules]
    C3 --> C4[Cache dependencies]
    C4 --> C5[Install dependencies]
    C5 --> C6[Run golangci-lint]
    C6 --> C7[Run tests]
    C7 --> C8[Upload coverage to codecov]
    C8 --> C9[Build API binary]
    C9 --> C10[Zip API binary]
    C10 --> C11[Upload API artifact]
    end

    C11 -->|Trigger| D[Test and Build Scraper]
    
    subgraph "Test and Build Scraper"
    D --> D1[Checkout code]
    D1 --> D2[Set up Go]
    D2 --> D3[Cache dependencies]
    D3 --> D4[Install dependencies]
    D4 --> D5[Run golangci-lint]
    D5 --> D6[Run tests]
    D6 --> D7[Upload coverage to codecov]
    D7 --> D8[Build Scraper binary]
    D8 --> D9[Zip Scraper binary]
    D9 --> D10[Upload Scraper artifact]
    end

    subgraph "Test and Build API Users"
    F --> F1[Checkout code]
    F1 --> F2[Set up Go]
    F2 --> F3[Cache dependencies]
    F3 --> F4[Install dependencies]
    F4 --> F5[Run golangci-lint]
    F5 --> F6[Run tests]
    F6 --> F7[Upload coverage to codecov]
    F7 --> F8[Build API Users binary]
    F8 --> F9[Zip API Users binary]
    F9 --> F10[Upload API Users artifact]
    end

    F10 -->|Trigger| G[Test and Build Authorizer]

    subgraph "Test and Build Authorizer"
    G --> G1[Checkout code]
    G1 --> G2[Set up Go]
    G2 --> G3[Cache dependencies]
    G3 --> G4[Install dependencies]
    G4 --> G5[Run golangci-lint]
    G5 --> G6[Run tests]
    G6 --> G7[Upload coverage to codecov]
    G7 --> G8[Build Authorizer binary]
    G8 --> G9[Zip Authorizer binary]
    G9 --> G10[Upload Authorizer artifact]
    end
    
    D10 -->|Trigger| E[Deploy]
    G10 -->|Trigger| E[Deploy]
    
    subgraph "Deploy"
    E --> E1[Checkout code]
    E1 --> E2[Download API artifact]
    E2 --> E3[Download Scraper artifact]
    E3 --> E4[Download API Users artifact]
    E4 --> E5[Download Authorizer artifact]
    E5 --> E6[Set up Terraform]
    E6 --> E7[Initialize Terraform]
    E7 --> E8[Plan Terraform]
    E8 --> E9[Apply Terraform]
    end
```