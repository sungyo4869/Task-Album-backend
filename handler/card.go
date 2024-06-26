package handler

import (
	"encoding/json"
	"net/http"
	"context"

	"github.com/sungyo4869/portfolio/model"
	"github.com/sungyo4869/portfolio/service"
)

type CardHandler struct {
	svc *service.CardService
}

func NewCardHandler(svc *service.CardService) *CardHandler {
	return &CardHandler{
		svc: svc,
	}
}

func (c *CardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		res, err := c.Read(r.Context())
		if err != nil {
			http.Error(w, "Error reading card", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		// リクエストボディをmodel.CreateCardRequestにデコード
		var req model.CreateCardRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.Title == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// reqをDBに格納する関数の呼び出し
		res, err := c.Create(r.Context(), &req)
		if err != nil {
			http.Error(w, "Error creating card", http.StatusInternalServerError)
			return
		}

		// resをjsonにエンコード
		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		// リクエストボディをmodel.CreateCardRequestにデコード
		var req model.UpdateCardRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		// reqをDBに格納する関数の呼び出し
		res, err := c.Update(r.Context(), &req)
		if err != nil {
			http.Error(w, "Error Updating card", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}

	case http.MethodDelete:
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CardHandler) Create(ctx context.Context, req *model.CreateCardRequest) (*model.CreateCardResponse, error) {
	var res model.CreateCardResponse

	result, err := h.svc.CreateCard(ctx, req.Title, req.Summary, req.TimeLimit)
	if err != nil {
		return nil, err
	}
	
	res.Card = *result
	
	return &res, nil
}

func (h *CardHandler) Read(ctx context.Context) (*model.ReadCardsResponse, error) {
	var res model.ReadCardsResponse

	result, err := h.svc.ReadCard(ctx)
	if err != nil {
		return nil, err
	}

	res.Cards = result

	return &res, nil
}

func (h *CardHandler) Update(ctx context.Context, req *model.UpdateCardRequest) (*model.UpdateCardResponse, error) {
	var res model.UpdateCardResponse

	result, err := h.svc.UpdateCard(ctx, req.Title, req.Summary, req.Description, req.TimeLimit, req.CardID)
	if err != nil {
		return nil, err
	}

	res.Card = *result

	return &res, nil
}

func (h *CardHandler) Delete(ctx context.Context, req *model.DeleteCardRequest) (*model.DeleteCardResponse, error) {
	var res model.DeleteCardResponse

	err := h.svc.DeleteCard(ctx, req.CardID)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
