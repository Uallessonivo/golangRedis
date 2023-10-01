package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golangRedis/model"
	"golangRedis/repository/order"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Order struct {
	Repo *order.RedisRepo
}

func (o *Order) Create(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	orderR := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	err := o.Repo.Insert(request.Context(), orderR)
	if err != nil {
		fmt.Println("error insert orderR: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(orderR)
	if err != nil {
		fmt.Println("error marshal orderR: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(res)
	writer.WriteHeader(http.StatusCreated)
}

func (o *Order) List(writer http.ResponseWriter, request *http.Request) {
	cursorStr := request.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64

	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := o.Repo.FindAll(request.Context(), order.FindAllPage{
		Size:   cursor,
		Offset: size,
	})

	if err != nil {
		fmt.Println("error find all order: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}

	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error marshal order: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(data)
	writer.WriteHeader(http.StatusOK)
}

func (o *Order) GetById(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")

	const base = 10
	const biteSize = 64

	orderID, err := strconv.ParseUint(idParam, base, biteSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	or, err := o.Repo.FindById(request.Context(), orderID)
	if errors.Is(err, order.ErrNotExists) {
		writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("error find by id: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(or); err != nil {
		fmt.Println("failed to marshal: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Order) UpdateById(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(request, "id")

	const base = 10
	const biteSize = 64

	orderID, err := strconv.ParseUint(idParam, base, biteSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	theOrder, err := o.Repo.FindById(request.Context(), orderID)
	if errors.Is(err, order.ErrNotExists) {
		writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("error find by id: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if theOrder.ShippedAt != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now
	case completedStatus:
		if theOrder.CompletedAt != nil || theOrder.ShippedAt == nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.Update(request.Context(), theOrder)
	if err != nil {
		fmt.Println("error update order: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(theOrder); err != nil {
		fmt.Println("failed to marshal: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Order) DeleteById(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")

	const base = 10
	const biteSize = 64

	orderID, err := strconv.ParseUint(idParam, base, biteSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.DeleteByID(request.Context(), orderID)
	if errors.Is(err, order.ErrNotExists) {
		writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("error find by id: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
