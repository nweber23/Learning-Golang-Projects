package handlers

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"image-process-service/middleware"
	"image-process-service/models"
	"image-process-service/processor"
	"image-process-service/storage"
)

const maxFileSize = 50 << 20 // 50MB

type ImageHandler struct {
	store     *models.Store
	storage   storage.Storage
	processor *processor.Processor
}

func NewImageHandler(store *models.Store, storage storage.Storage, processor *processor.Processor) *ImageHandler {
	return &ImageHandler{
		store:     store,
		storage:   storage,
		processor: processor,
	}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		middleware.JSONError(w, http.StatusBadRequest, "File too large, max 50MB")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		middleware.JSONError(w, http.StatusBadRequest, "Missing file field")
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			middleware.JSONError(w, http.StatusInternalServerError, "Error closing file")
			return
		}
	}()

	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" && contentType != "image/webp" {
		middleware.JSONError(w, http.StatusBadRequest, "Invalid file type, only images are accepted")
		return
	}

	img, format, err := image.Decode(file)
	if err != nil {
		middleware.JSONError(w, http.StatusBadRequest, "Could not decode image")
		return
	}

	bounds := img.Bounds()
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), header.Filename)

	if _, err := file.Seek(0, 0); err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error processing file")
		return
	}

	url, err := h.storage.Upload(filename, file)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error uploading file")
		return
	}

	imageRecord := &models.Image{
		ID:          uuid.New().String(),
		UserID:      userID,
		Filename:    header.Filename,
		OriginalURL: url,
		Width:       bounds.Max.X - bounds.Min.X,
		Height:      bounds.Max.Y - bounds.Min.Y,
		Size:        header.Size,
		Format:      format,
		CreatedAt:   time.Now(),
	}

	if err := h.store.SaveImage(imageRecord); err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error saving image metadata")
		return
	}

	WriteJSON(w, http.StatusCreated, models.ImageResponse{
		ID:        imageRecord.ID,
		Filename:  imageRecord.Filename,
		URL:       imageRecord.OriginalURL,
		Width:     imageRecord.Width,
		Height:    imageRecord.Height,
		Size:      imageRecord.Size,
		Format:    imageRecord.Format,
		CreatedAt: imageRecord.CreatedAt,
	})
}

func (h *ImageHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	images, total, err := h.store.ListImagesByUser(userID, page, limit)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error retrieving images")
		return
	}

	var responses []*models.ImageResponse
	for _, img := range images {
		responses = append(responses, &models.ImageResponse{
			ID:        img.ID,
			Filename:  img.Filename,
			URL:       img.OriginalURL,
			Width:     img.Width,
			Height:    img.Height,
			Size:      img.Size,
			Format:    img.Format,
			CreatedAt: img.CreatedAt,
		})
	}

	if responses == nil {
		responses = []*models.ImageResponse{}
	}

	WriteJSON(w, http.StatusOK, models.ImageListResponse{
		Images: responses,
		Page:   page,
		Limit:  limit,
		Total:  total,
	})
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	imageID := mux.Vars(r)["id"]

	img, err := h.store.FindImageByID(imageID)
	if err != nil || img == nil {
		middleware.JSONError(w, http.StatusNotFound, "Image not found")
		return
	}

	if img.UserID != userID {
		middleware.JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	WriteJSON(w, http.StatusOK, models.ImageResponse{
		ID:        img.ID,
		Filename:  img.Filename,
		URL:       img.OriginalURL,
		Width:     img.Width,
		Height:    img.Height,
		Size:      img.Size,
		Format:    img.Format,
		CreatedAt: img.CreatedAt,
	})
}

func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	imageID := mux.Vars(r)["id"]

	img, err := h.store.FindImageByID(imageID)
	if err != nil || img == nil {
		middleware.JSONError(w, http.StatusNotFound, "Image not found")
		return
	}

	if img.UserID != userID {
		middleware.JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.storage.Delete(img.OriginalURL); err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error deleting file")
		return
	}

	if err := h.store.DeleteImage(imageID); err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error deleting image metadata")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ImageHandler) TransformImage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	imageID := mux.Vars(r)["id"]

	img, err := h.store.FindImageByID(imageID)
	if err != nil || img == nil {
		middleware.JSONError(w, http.StatusNotFound, "Image not found")
		return
	}

	if img.UserID != userID {
		middleware.JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.TransformationRequest
	if err := DecodeJSON(r, &req); err != nil {
		middleware.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	reader, err := h.storage.Download(img.OriginalURL)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error reading original image")
		return
	}

	srcImg, _, err := image.Decode(reader)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error decoding original image")
		return
	}

	resultImg, err := h.processor.Process(srcImg, &req.Transformations)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error processing image")
		return
	}

	outputFormat := req.Transformations.Format
	if outputFormat == "" {
		outputFormat = img.Format
	}

	var buf bytes.Buffer
	if err := processor.Encode(resultImg, outputFormat, &buf); err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error encoding result image")
		return
	}

	resultFilename := fmt.Sprintf("transformed-%s.%s", uuid.New().String(), outputFormat)
	resultURL, err := h.storage.Upload(resultFilename, &buf)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error saving transformed image")
		return
	}

	transformation := &models.Transformation{
		ID:              uuid.New().String(),
		UserID:          userID,
		OriginalImageID: imageID,
		ResultURL:       resultURL,
		CreatedAt:       time.Now(),
	}

	if err := h.store.SaveTransformation(transformation); err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error saving transformation record")
		return
	}

	WriteJSON(w, http.StatusOK, transformation)
}
