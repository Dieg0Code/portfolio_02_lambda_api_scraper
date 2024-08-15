
# DynamoDB table for products
resource "aws_dynamodb_table" "products_table" {
  name         = "Products"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ProductID"


  attribute {
    name = "ProductID"
    type = "S"
  }

  lifecycle {
    ignore_changes = [name]
  }
}

resource "aws_dynamodb_table" "users_table" {
  name        = "Users"
  billing_mode = "PAY_PER_REQUEST"
  hash_key    = "UserID"

  attribute {
    name = "UserID"
    type = "S"
  }

  lifecycle {
    ignore_changes = [name]
  }
}

