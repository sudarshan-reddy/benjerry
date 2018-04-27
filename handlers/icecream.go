package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sudarshan-reddy/benjerry/httputils"
	"github.com/sudarshan-reddy/benjerry/models"
)

//IceCreamHandler holds handler related data
type IceCreamHandler struct {
	iceCreamStore models.IceCreamStore
}

//NewIceCreamHandler returns a new instance of IceCreamHandler
func NewIceCreamHandler(iceCreamStore models.IceCreamStore) *IceCreamHandler {
	return &IceCreamHandler{
		iceCreamStore: iceCreamStore,
	}
}

//PostIceCreamData creates ice cream data
func (i *IceCreamHandler) PostIceCreamData(w http.ResponseWriter, r *http.Request) {
	var iceCreamTask models.IceCream
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&iceCreamTask); err != nil {
		msg := fmt.Sprintf("invalid input format. error: %s", err)
		httputils.WriteHandlerError(httputils.NewFormatError(msg), r, w)
		return
	}

	ctx := r.Context()
	if err := i.iceCreamStore.StoreContext(ctx, iceCreamTask); err != nil {
		httputils.WriteHandlerError(httputils.NewUnexpectedError(err), r, w)
		return
	}

	if err := httputils.WriteJSON(http.StatusCreated, "", w); err != nil {
		httputils.WriteHandlerError(httputils.NewUnexpectedError(err), r, w)
		return
	}
}

//GetIceCreamData gets ice cream data for a particular name
func (i *IceCreamHandler) GetIceCreamData(w http.ResponseWriter, r *http.Request) {
	iceCreamName := chi.URLParam(r, "ice-cream-name")

	iceCreamData, err := i.iceCreamStore.Get(iceCreamName)

	if err != nil {
		if err == models.ErrNoRows {
			httputils.WriteHandlerError(httputils.
				NewNotFoundError(fmt.Sprintf("Icecream: %s Not Found", iceCreamName)), r, w)
			return
		}
		httputils.WriteHandlerError(httputils.NewUnexpectedError(err), r, w)
		return
	}

	if err := httputils.WriteJSON(http.StatusOK, iceCreamData, w); err != nil {
		httputils.WriteHandlerError(httputils.NewUnexpectedError(err), r, w)
		return
	}
}

//UpdateIceCreamData updates records based on change
//This would not update primary key
func (i *IceCreamHandler) UpdateIceCreamData(w http.ResponseWriter, r *http.Request) {
	var iceCreamTask models.IceCream
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&iceCreamTask); err != nil {
		msg := fmt.Sprintf("invalid input format. error: %s", err)
		httputils.WriteHandlerError(httputils.NewFormatError(msg), r, w)
		return
	}

	err := i.iceCreamStore.Update(iceCreamTask)
	if err != nil {
		httputils.WriteHandlerError(httputils.NewUnexpectedError(err), r, w)
		return
	}

	if err := httputils.WriteJSON(http.StatusOK, "", w); err != nil {
		httputils.WriteHandlerError(httputils.NewUnexpectedError(err), r, w)
		return
	}
}

//DeleteIceCreamData delete records based on primary key
func (i *IceCreamHandler) DeleteIceCreamData(w http.ResponseWriter, r *http.Request) {
	iceCreamName := chi.URLParam(r, "ice-cream-name")

	fmt.Println(iceCreamName)
	err := i.iceCreamStore.Delete(iceCreamName)

	if err != nil {
		httputils.WriteHandlerError(httputils.NewUnexpectedError(err), r, w)
		return
	}

	if err := httputils.WriteJSON(http.StatusOK, "", w); err != nil {
		httputils.WriteHandlerError(httputils.NewUnexpectedError(err), r, w)
		return
	}
}
