package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type Item struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"_"`
	UpdatedAt    time.Time `json:"_"`
	ProductName  string    `json:"productname"`
	Descript     string    `json:"descript"`
	ProductId    uuid.UUID `json:"productid"`
	// Product_Cat int64     `json:"productcat"`
	ProductLocation string `json:"productlocation"`
	Product_Cat  string    `json:"productcat"`
}

func (i *Item) Create(conn *pgx.Conn, userID string) error {
	i.ProductName = strings.Trim(i.ProductName, " ")
	if len(i.ProductName) < 1 {
		return fmt.Errorf("ProductName must not be empty.")
	}
	if len(i.ProductLocation) < 1 {
		return fmt.Errorf("productlocation must not be empty.")
	}
	if i.Product_Cat < "" {
		i.Product_Cat = "Generic"
	}
	now := time.Now()

	row := conn.QueryRow(context.Background(), "INSERT INTO item (productname, descript, productid, productcat, productlocation,created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6,$7) RETURNING id, productid", i.ProductName, i.Descript, userID, i.Product_Cat, i.ProductLocation, now, now)

	err := row.Scan(&i.ID, &i.ProductId)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the item")
	}

	return nil
}

func GetAllItems(conn *pgx.Conn) ([]Item, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, productname, descript, productid, productlocation, productcat FROM item")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error getting items")
	}

	var items []Item
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ID, &item.ProductName, &item.ProductLocation, &item.Descript, &item.ProductId, &item.Product_Cat)
		if err != nil {
			fmt.Printf("items liost checkpoint")
			fmt.Println(err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func GetItemsBeingSoldByUser(userID string, conn *pgx.Conn) ([]Item, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, productname, productlocation, productcat, descript, productid  FROM item WHERE productid = $1", userID)
	if err != nil {
		fmt.Printf("Error getting items %v", err)
		return nil, fmt.Errorf("There was an error getting the items")
	}

	var items []Item
	for rows.Next() {
		i := Item{}
		err = rows.Scan(&i.ID, &i.ProductName, &i.ProductLocation, &i.Product_Cat, &i.Descript, &i.ProductId)
		if err != nil {
			fmt.Printf("Error scaning item: %v", err)
			continue
		}
		items = append(items, i)
	}

	return items, nil
}

func (i *Item) Update(conn *pgx.Conn) error {
	i.ProductName = strings.Trim(i.ProductName, " ")
	if len(i.ProductName) < 1 {
		return fmt.Errorf("ProductName must not be empty")
	}
	if len(i.ProductLocation) < 1 {
		return fmt.Errorf("ProductLocation must not be empty")
	}

	// if i.ProductLocation < 0 {
	// 	i.ProductLocation = 0
	// }
	now := time.Now()
	_, err := conn.Exec(context.Background(), "UPDATE item SET productname=$1, productlocation=$6, descript=$2, productcat=$3, updated_at=$4 WHERE id=$5", i.ProductName, i.ProductLocation, i.Descript, i.Product_Cat, now, i.ID)

	if err != nil {
		fmt.Printf("Error updating item: (%v)", err)
		return fmt.Errorf("Error updating item")
	}

	return nil
}

func FindItemById(id uuid.UUID, conn *pgx.Conn) (Item, error) {
	row := conn.QueryRow(context.Background(), "SELECT productname, descript, productid, productcat FROM item WHERE id=$1", id)
	fmt.Printf("find item checkpoint")
	item := Item{
		ID: id,
	}
	err := row.Scan(&item.ProductName, &item.Descript, &item.ProductId, &item.Product_Cat)
	if err != nil {
		fmt.Printf("returning nil")
		return item, fmt.Errorf("The item doesn't exist")
	}

	return item, nil
}
