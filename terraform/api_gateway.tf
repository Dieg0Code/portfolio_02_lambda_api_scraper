resource "aws_api_gateway_rest_api" "api" {
  name        = "api_scraper"
  description = "API Scraper"
}

# Resource for API Gateway /products endpoint
resource "aws_api_gateway_resource" "products" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "products"
}

# Resource for API Gateway /products/{id} endpoint
resource "aws_api_gateway_resource" "product" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_resource.products.id
  path_part   = "{id}"
}

# Resource for API Gateway /Users endpoint
resource "aws_api_gateway_resource" "users" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "Users"
}

# Resource for API Gateway /Scraper endpoint
resource "aws_api_gateway_resource" "scraper" {
    rest_api_id = aws_api_gateway_rest_api.api.id
    parent_id   = aws_api_gateway_rest_api.api.root_resource_id
    path_part   = "Scraper"
}

# Method for GET /products endpoint
resource "aws_api_gateway_method" "get_products" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.products.id
  http_method   = "GET"
  authorization = "NONE"
}

# Method for GET /products/{id} endpoint
resource "aws_api_gateway_method" "get_product" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.product.id
  http_method   = "GET"
  authorization = "NONE"
}

# Method for POST /Users endpoint
resource "aws_api_gateway_method" "post_users" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.users.id
  http_method   = "POST"
  authorization = "NONE"
}

# Method for POST /Scraper endpoint (update data)
resource "aws_api_gateway_method" "post_scraper" {
    rest_api_id = aws_api_gateway_rest_api.api.id
    resource_id = aws_api_gateway_resource.scraper.id
    http_method = "POST"
    authorization = "NONE"
}


# Integration for GET /products endpoint
resource "aws_api_gateway_integration" "products_lambda_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.products.id
  http_method = aws_api_gateway_method.get_products.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.api_scraper.invoke_arn
}

# Integration for GET /products/{id} endpoint
resource "aws_api_gateway_integration" "product_lambda_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.product.id
  http_method = aws_api_gateway_method.get_product.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.api_scraper.invoke_arn
}

# Integration for POST /Users endpoint
resource "aws_api_gateway_integration" "users_lambda_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.users.id
  http_method = aws_api_gateway_method.post_users.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.api_scraper.invoke_arn
}

# Integration for POST /Scraper endpoint
resource "aws_api_gateway_integration" "scraper_lambda_integration" {
    rest_api_id = aws_api_gateway_rest_api.api.id
    resource_id = aws_api_gateway_resource.scraper.id
    http_method = aws_api_gateway_method.post_scraper.http_method

    integration_http_method = "POST"
    type                    = "AWS_PROXY"
    uri                     = aws_lambda_function.api_scraper.invoke_arn
}

