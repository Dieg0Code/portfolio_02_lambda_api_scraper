resource "aws_lambda_function" "api_products" {
  filename      = "api_products_lambda.zip"
  function_name = "api_products"
  role          = aws_iam_role.lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  memory_size   = 128
  timeout       = 90

  source_code_hash = filebase64sha256("api_products_lambda.zip")

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.products_table.name
    }
  }
}

resource "aws_lambda_function" "scraper" {
  filename      = "scraper_lambda.zip"
  function_name = "scraper"
  role          = aws_iam_role.lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  memory_size   = 128
  timeout       = 150

  source_code_hash = filebase64sha256("scraper_lambda.zip")

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.products_table.name
    }
  }
}

resource "aws_lambda_function" "api_users" {
  filename      = "api_users_lambda.zip"
  function_name = "api_users"
  role          = aws_iam_role.lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  memory_size   = 128
  timeout       = 90

  source_code_hash = filebase64sha256("api_users_lambda.zip")

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.users_table.name
    }
  }
}
