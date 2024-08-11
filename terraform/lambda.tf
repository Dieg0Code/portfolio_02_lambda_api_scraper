resource "aws_lambda_function" "api_scraper" {
  filename      = "api_scraper.zip"
  function_name = "api_scraper"
  role          = aws_iam_role.lambda_role.arn
  handler       = "api_scraper_lambda"
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 150

  source_code_hash = filebase64sha256("../api_scraper_lambda.zip")

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.products_table.name
    }
  }
}
