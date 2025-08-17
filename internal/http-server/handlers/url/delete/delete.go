package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type Request struct {
	ID string `json:"id" validate:"required"`
}

type Response struct {
	resp.Response
}

type URLDeleter interface {
	DeleteURL(id string) error
}

func New(log *slog.Logger, deleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Error("id is empty")
			render.JSON(w, r, resp.Error("id is required"))
			return
		}

		log.Info("attempting to delete URL", slog.String("id", id))

		err := deleter.DeleteURL(id)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("id", id))
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}

		log.Info("url deleted successfully", slog.String("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
