package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func promptHandler(c *gin.Context) {
	payload := Ask{
		Prefix: c.PostForm("prefix"),
		Prompt: c.PostForm("prompt"),
		Secret: os.Getenv("KAHUNA_API_KEY"),
	}
	if os.Getenv("SKIP_CLAUDE") != "" {
		c.HTML(http.StatusOK, "protected.html", gin.H{})
		return
	}
	response, err := kahunaAPI(payload)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "protected.html", gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "protected.html", response.getGinH())

}

type Ask struct {
	Prefix string `json:"prefix"`
	Prompt string `json:"prompt"`
	Secret string `json:"secret"`
}

type Answer struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
}

func (a Answer) getGinH() gin.H {
	return gin.H{"response": a.Response, "prompt": a.Prompt}
}

func kahunaAPI(payload Ask) (Answer, error) {
	url := "https://kahuna.one/ask"
	//url := "http://localhost:8080/ask"
	var result Answer
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Answer{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	return result, err
}
