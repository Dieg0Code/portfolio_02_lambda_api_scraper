![infra](infra_p02.png)

# Serverless Api-Scraper

Serverless API. Scrape data from a supermarket website and store it in a dynamoDB database.

Its a mono-repo with the following structure:

- `api` - Serverless API
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

- `[POST] /api/v1/products` - Update Data (this will take a while, 1 min aprox)

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

    %% Conexiones
    ProductControllerImpl --> ProductServiceImpl : productService
    ProductServiceImpl --> ProductRepositoryImpl : productRepository
    ProductRepositoryImpl --> Product : Product
    ProductServiceImpl --> ProductResponse : ProductResponse
    ProductServiceImpl --> UpdateDataRequest : UpdateDataRequest
    ProductControllerImpl --> BaseResponse : BaseResponse


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

## Interaction Diagram between Lambdas and DynamoDB

```mermaid
sequenceDiagram
    participant User
    participant APIGateway
    participant APILambda as api_products Lambda
    participant ScraperLambda as Scraper Lambda
    participant DynamoDB

    User->>APIGateway: Request (GetAll/GetByID)
    APIGateway->>APILambda: Invoke Lambda
    APILambda->>DynamoDB: Query Products
    DynamoDB-->>APILambda: Return Products
    APILambda-->>APIGateway: Respond with Products
    APIGateway-->>User: Return Response

    User->>APIGateway: Request (POST /api/v1/products)
    APIGateway->>APILambda: Invoke Lambda
    APILambda->>ScraperLambda: Trigger Scraper Lambda
    ScraperLambda->>SupermarketWebpage: Scrape Data
    ScraperLambda->>DynamoDB: Update Products
    ScraperLambda-->>APILambda: Return Success
    APILambda-->>APIGateway: Respond with Success
    APIGateway-->>User: Return Response
```

## Terraform Infrastructure

```mermaid
graph TD
    %% Subgraph for S3 and DynamoDB, which are part of Terraform backend
    subgraph "Terraform State Management"
        s3_bucket[S3 Bucket - terraform-state-api-scraper]
        dynamodb_table_locks[DynamoDB Table - terraform_locks]
    end

    %% Subgraph for IAM roles and policies
    subgraph "IAM Roles and Policies"
        iam_role[Lambda IAM Role - lambda_role]
        lambda_policy[Lambda Policy - lambda_policy]
        invoke_policy[Invoke Policy - lambda_invoke_policy]
        basic_exec_role[AWSLambdaBasicExecutionRole]

        %% Attachments
        iam_role --> lambda_policy
        iam_role --> invoke_policy
        iam_role --> basic_exec_role
    end

    %% Subgraph for Lambda functions
    subgraph "Lambda Functions"
        lambda_api_products[Lambda - api_products]
        lambda_scraper[Lambda - scraper]
    end

    %% Subgraph for API Gateway
    subgraph "API Gateway"
        api_gateway[API Gateway - api_scraper]
        api_resource[Resource - /api/v1/products]
        get_method[GET /api/v1/products]
        post_method[POST /api/v1/products]

        %% Methods and resources association
        api_resource --> get_method
        api_resource --> post_method
    end

    %% DynamoDB Table for Products
    dynamodb_table[(DynamoDB Table - Products)]

    %% Relationships between resources
    s3_bucket --> terraform_backend
    dynamodb_table_locks --> terraform_backend

    dynamodb_table --> lambda_api_products
    dynamodb_table --> lambda_scraper

    iam_role --> lambda_api_products
    iam_role --> lambda_scraper

    lambda_api_products --> api_gateway
    lambda_scraper --> lambda_api_products

    api_gateway --> api_resource
    api_resource --> get_method
    api_resource --> post_method

    get_method --> lambda_api_products
    post_method --> lambda_api_products

```

## CI/CD Pipeline

```mermaid
graph TD
    A[Push to main branch] -->|Trigger| B[CI/CD Pipeline]
    B --> C[Test and Build API]
    B --> D[Test and Build Scraper]
    
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
    
    C11 --> E[Deploy]
    D10 --> E
    
    subgraph "Deploy"
    E --> E1[Checkout code]
    E1 --> E2[Download API artifact]
    E2 --> E3[Download Scraper artifact]
    E3 --> E4[Set up Terraform]
    E4 --> E5[Initialize Terraform]
    E5 --> E6[Plan Terraform]
    E6 --> E7[Apply Terraform]
    end
```