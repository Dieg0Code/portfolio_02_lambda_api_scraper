
# DynamoDB table for products
resource "aws_dynamodb_table" "products_table" {
  name = "Products"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "ProductID"


  attribute {
    name = "ProductID"
    type = "S"
  }

    attribute {
        name = "Name"
        type = "S"
    }

    attribute {
      name = "Category"
      type = "S"
    }

    attribute {
      name = "OriginalPrice"
        type = "N"
    }

    attribute {
      name = "DiscountedPrice"
        type = "N"
    }
}

