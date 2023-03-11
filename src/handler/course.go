package handler

import (
	"CareerBoost/src/entity"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *handler) addNewCourse(ctx *gin.Context) {
	var courseBody entity.CourseAdd
	if err := h.BindBody(ctx, &courseBody); err != nil {
		fmt.Println(err)
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request", nil)
		return
	}

	var courseDB entity.Course
	courseDB.Judul = courseBody.Judul
	courseDB.Deskripsi = courseBody.Deskripsi
	courseDB.Intro = courseBody.Intro
	courseDB.Rate = courseBody.Rate
	courseDB.Price = courseBody.Price

	var playlists []entity.Playlist
	for _, playlist := range courseBody.Playlist {
		var videos []entity.Video
		count := 0
		for _, video := range playlist.Video {

			durasi, err := time.ParseDuration(video.Durasi)
			if err != nil {
				h.ErrorResponse(ctx, http.StatusBadRequest, "invalid video duration", nil)
				return
			}

			count += int(durasi)
			videos = append(videos, entity.Video{
				Link:       video.Link,
				Judul:      video.Judul,
				Durasi:     video.Durasi,
				PlaylistID: playlist.ID,
			})
		}
		var playl entity.Playlist

		playl.Nama = playlist.Nama
		playl.Video = videos
		playl.Course = courseDB
		playl.Durasi = time.Duration(count)

		playlists = append(playlists, playl)

		if err := h.db.Create(&playl).Error; err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, "failed to add playlist", nil)
			return
		}
	}

	courseDB.Playlist = playlists

	if err := h.db.Create(&courseDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create course", nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Course berhasil ditambahkan", nil, nil)
}
