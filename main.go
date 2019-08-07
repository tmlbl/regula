package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Registry API spec: https://docs.docker.com/registry/spec/api/

func main() {
	api := gin.New()
	api.Use(RequestLogger())

	api.HEAD("/v2/:name/blobs/:digest", handleBlobExists)
	api.POST("/v2/:name/blobs/uploads", handleStartBlobUpload)
	api.PATCH("/v2/:name/blobs/uploads/:uuid", handleChunkUpload)
	api.PUT("/v2/:name/blobs/uploads/:uuid", handleUploadCompleted)

	api.Run(":5000")
}

func handleBlobExists(c *gin.Context) {
	c.Status(404)
}

func handleStartBlobUpload(c *gin.Context) {
	id := uuid.New()
	location := fmt.Sprintf("/v2/%s/blobs/uploads/%s",
		c.Param("name"), id.String())
	c.Header("Location", location)
	c.Header("Docker-Upload-UUID", id.String())
	c.Header("Content-Length", "0")
	c.Status(202)
}

func handleChunkUpload(c *gin.Context) {
	name := c.Param("name")
	id := c.Param("uuid")

	// When do these get sent? Only for chunk n > 0 ?
	length := c.GetHeader("Content-Length")
	rangeStr := c.GetHeader("Content-Range")
	if length != "" || rangeStr != "" {
		log.Println("Length:", length, "Range:", rangeStr)
	}

	defer c.Request.Body.Close()
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading chunk", err)
	}

	location := fmt.Sprintf("/v2/%s/blobs/uploads/%s",
		name, id)

	c.Status(http.StatusAccepted)
	c.Header("Location", location)
	rang := fmt.Sprintf("0-%d", len(data))
	fmt.Println("Range:", rang)
	c.Header("Range", rang)
	c.Header("Content-Length", "0")
	c.Header("Docker-Upload-UUID", id)
}

func handleUploadCompleted(c *gin.Context) {
	name := c.Param("name")
	digest := c.Query("digest")
	c.Status(http.StatusCreated)
	c.Header("Location", fmt.Sprintf("/v2/%s/blobs/%s", name, digest))
	c.Header("Docker-Content-Digest", digest)
	c.Header("Content-Length", "0")
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Method, c.Request.URL.Path)
	}
}
