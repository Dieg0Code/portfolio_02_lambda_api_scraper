resource "aws_api_gateway_rest_api" "api" {
  name        = "api_scraper"
  description = "API Scraper"
}

# Resource for API Gateway /api/v1/products endpoint
resource "aws_api_gateway_resource" "products" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "api/v1/products"
}

# Resource for API Gateway /api/v1/products/{productId} endpoint
resource "aws_api_gateway_resource" "product" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_resource.products.id
  path_part   = "{productId}"
}

# Method for GET /api/v1/products endpoint
resource "aws_api_gateway_method" "get_products" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.products.id
  http_method   = "GET"
  authorization = "NONE"
}

# Method for GET /api/v1/products/{productId} endpoint
resource "aws_api_gateway_method" "get_product" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.product.id
  http_method   = "GET"
  authorization = "NONE"
}

# Method for POST /api/v1/products endpoint
resource "aws_api_gateway_method" "post_products" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.products.id
  http_method   = "POST"
  authorization = "NONE"
}

# Integration for GET /api/v1/products endpoint
resource "aws_api_gateway_integration" "products_lambda_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.products.id
  http_method = aws_api_gateway_method.get_products.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.api_scraper.invoke_arn
}

# Integration for GET /api/v1/products/{productId} endpoint
resource "aws_api_gateway_integration" "product_lambda_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.product.id
  http_method = aws_api_gateway_method.get_product.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.api_scraper.invoke_arn
}

# Integration for POST /api/v1/products endpoint
resource "aws_api_gateway_integration" "post_products_lambda_integration" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.products.id
  http_method = aws_api_gateway_method.post_products.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.api_scraper.invoke_arn
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api_scraper.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api.execution_arn}/*/*"
}

resource "aws_api_gateway_stage" "api_stage" {
  deployment_id = aws_api_gateway_deployment.api_deployment.id
  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name = "dev"
}

resource "aws_api_gateway_deployment" "api_deployment" {
  depends_on = [
    aws_api_gateway_method.get_products,
    aws_api_gateway_method.get_product,
    aws_api_gateway_method.post_products
  ]

  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name  = "dev"
}

output "api_gateway_invoke_url" {
  value = "https://${aws_api_gateway_rest_api.api.id}.execute-api.sa-east-1.amazonaws.com/${aws_api_gateway_stage.api_stage.stage_name}"
}
