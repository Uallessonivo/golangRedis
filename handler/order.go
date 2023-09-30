package handler

import (
	"fmt"
	"net/http"
)

type Order struct {
}

func (o *Order) Create(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("create order")
}

func (o *Order) List(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("list order")
}

func (o *Order) GetById(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("get order")
}

func (o *Order) UpdateById(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("update order")
}

func (o *Order) DeleteById(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("delete order")
}
