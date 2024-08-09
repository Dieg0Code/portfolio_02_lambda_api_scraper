resource "aws_lambda_function" "api_scraper" {
  filename = "api_scraper.zip"
  function_name = "api_scraper"
  role = aws_iam_role.lambda_role.arn
  handler = "api_scraper_lambda"
  runtime = "go1.x"

  source_code_hash = filebase64sha256("api_scraper.zip")
}