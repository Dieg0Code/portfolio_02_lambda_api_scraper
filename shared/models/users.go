package models

type User struct {
	UserID   string `json:"id" dynamodbav:"UserID"`
	Username string `json:"username" dynamodbav:"Username"`
	Email    string `json:"email" dynamodbav:"Email"`
	Password string `json:"password" dynamodbav:"Password"`
	Role     string `json:"role" dynamodbav:"Role"`
}
