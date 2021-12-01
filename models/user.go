package models

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	// "strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	tokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"_"`
	UpdatedAt       time.Time `json:"_"`
	FullName		string	  `json:"fullname"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Password        string    `json:"password"`
	PasswordConfirm string    `json:"password_confirm"`
	PhoneNo 		uint64	  `json:"phoneno"`
    StoreName		string	  `json:"storename"`
    StoreAddress	string	  `json:"storeaddress"`
    PinCode			int64    `json:"pincode"`
}


// type User struct {
// 	ID              uuid.UUID `json:"id"`
// 	CreatedAt       time.Time `json:"_"`
// 	UpdatedAt       time.Time `json:"_"`
// 	Product_id		string	  `json:"productid"`
// 	Email           string    `json:"email"`
// 	PasswordHash    string    `json:"-"`
// 	Password        string    `json:"password"`
// 	PasswordConfirm string    `json:"password_confirm"`
// 	PhoneNo 		uint64	  `json:"phoneno"`
//     StoreName		string	  `json:"storename"`
//     StoreAddress	string	  `json:"storeaddress"`
//     PinCode			int64     `json:"pincode"`
// }

func (u *User) Register(conn *pgx.Conn) error {
	if len(u.Password) < 4 || len(u.PasswordConfirm) < 4 {
		return fmt.Errorf("Password must be at least 4 characters long.")
	}

	if u.Password != u.PasswordConfirm {
		return fmt.Errorf("Passwords do not match.")
	}

	if len(u.Email) < 4 {
		return fmt.Errorf("Email must be at least 4 characters long.")
	}

	// ph := strconv.Itoa(u.PhoneNo)
	// if len(ph) < 10 && len(ph) >10 {
	// 	return fmt.Errorf("phone no must be at least 10 numbers long.")
	// }	

	// if len(u.PinCode) < 6  {
	// 	return fmt.Errorf("pin no must be at least 6 numbers long.")
	// }	

	u.Email = strings.ToLower(u.Email)
	row := conn.QueryRow(context.Background(), "SELECT id from user_account WHERE email = $1", u.Email)
	userLookup := User{}
	err := row.Scan(&userLookup)
	if err != pgx.ErrNoRows {
		fmt.Println("found user")
		fmt.Println(userLookup.Email)
		return fmt.Errorf("A user with that email already exists")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("There was an error creating your account.")
	}
	u.PasswordHash = string(pwdHash)

	now := time.Now()
	_, err = conn.Exec(context.Background(), "INSERT INTO user_account (created_at, updated_at, fullname, email, password_hash, phoneno, storename, storeaddress, pincode) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)", now, now, u.FullName, u.Email, u.PasswordHash, u.PhoneNo, u.StoreName, u.StoreAddress, u.PinCode)

	return err
}

// GetAuthToken returns the auth token to be used
func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err
}

// IsAuthenticated checks to make sure password is correct and user is active
func (u *User) IsAuthenticated(conn *pgx.Conn) error {
	row := conn.QueryRow(context.Background(), "SELECT id, password_hash from user_account WHERE email = $1", u.Email)
	err := row.Scan(&u.ID, &u.PasswordHash)
	if err == pgx.ErrNoRows {
		fmt.Println("User with email not found")
		return fmt.Errorf("Invalid login credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return fmt.Errorf("wrong password credentials")
	}

	return nil
}

func IsTokenValid(tokenString string) (bool, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// fmt.Printf("Parsing: %v \n", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			fmt.Printf("step1")
			return nil, fmt.Errorf("Token signing method is not valid: %v", token.Header["alg"])
		}
		fmt.Printf("step2")
		return tokenSecret, nil
	})

	if err != nil {
		fmt.Printf("Err %v \n", err)
		return false, ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Println(claims)
		userID := claims["user_id"]
		return true, userID.(string)
	} else {
		fmt.Printf("The alg header %v \n", claims["alg"])
		fmt.Println(err)
		return false, "uuid.UUID{}"
	}
}
