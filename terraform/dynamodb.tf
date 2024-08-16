resource "aws_dynamodb_table" "products_table" {
  name         = "Products"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ProductID"

  attribute {
    name = "ProductID"
    type = "S"
  }
}

resource "aws_dynamodb_table" "users_table" {
  name        = "Users"
  billing_mode = "PROVISIONED"
  hash_key    = "UserID"

  attribute {
    name = "UserID"
    type = "S"
  }

  attribute {
    name = "Email"
    type = "S"
  }

  global_secondary_index {
    name               = "EmailIndex"
    hash_key           = "Email"
    projection_type    = "ALL"
    write_capacity     = 10
    read_capacity      = 10
  }

  write_capacity = 10
  read_capacity  = 10
}
