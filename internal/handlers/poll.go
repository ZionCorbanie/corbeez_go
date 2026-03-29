package handlers

import (
	"goth/internal/middleware"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)


type PollHandler struct {
	pollStore store.PollStore
}

type PollHandlerParams struct {
	PollStore store.PollStore
}

func NewPollHandler(params PollHandlerParams) *PollHandler {
	return &PollHandler{
		pollStore: params.PollStore,
	}
}

func (h *PollHandler) Show(w http.ResponseWriter, r *http.Request) {
	//function shows current active poll
	user := middleware.GetUser(r.Context())
	active, err := h.pollStore.GetActivePoll()
    poll, voted := h.pollStore.GetPollVotes(active.ID, user.ID)
	if err != nil{
		return
	}

	c := templates.Poll(poll, voted)
	err = templates.Layout(c, "Corbeez").Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *PollHandler) Vote(w http.ResponseWriter, r *http.Request) {

    pollId, err := strconv.ParseInt(chi.URLParam(r, "pollId"), 10, 64)

    r.ParseForm()
    option, err := strconv.ParseInt(r.Form.Get("option"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid option id", http.StatusBadRequest)
        return
    }

	user := middleware.GetUser(r.Context())

    if user == nil {
        http.Error(w, "Not authorized", http.StatusUnauthorized)
        return
    }

    poll, voted := h.pollStore.GetPollVotes(uint(pollId), user.ID)
	if poll == nil || voted{
        http.Error(w, "Error voting", http.StatusBadRequest)
		return
	}

	if time.Now().Before(poll.StartDate){
        http.Error(w, "Te vroeg", http.StatusBadRequest)
		return
	}else if time.Now().After(poll.EndDate){
        http.Error(w, "Te laat", http.StatusBadRequest)
		return
	}

    err = h.pollStore.VotePoll(uint(pollId), uint(option), user.ID)
    if err != nil {
        http.Error(w, "Error voting", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/poll", http.StatusSeeOther)
}

func (h *PollHandler) RemoveVote(w http.ResponseWriter, r *http.Request) {
    pollId, err := strconv.ParseInt(chi.URLParam(r, "pollId"), 10, 64)

    user := middleware.GetUser(r.Context())

    if user == nil {
        http.Error(w, "Not authorized", http.StatusUnauthorized)
        return
    }

    err = h.pollStore.DeletePollVote(uint(pollId), user.ID)
    if err != nil {
        http.Error(w, "Error deleting vote", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/poll", http.StatusSeeOther)
}
