
# DynamoDB table for products
resource "aws_dynamodb_table" "products_table" {
  name         = "Products"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ProductID"


  attribute {
    name = "ProductID"
    type = "S"
  }
}

