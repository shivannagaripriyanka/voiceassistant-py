package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"fmt"
	"strings"

	"voiceassistant/models"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v4"
)

type FlaskOut struct {
	Keyword []string `json:"keyword"`
}
type emotionOut struct{
	Message string `json:"text"`
}


func SpeechToText(c *gin.Context) {

	var b bytes.Buffer
	var audio bytes.Buffer
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("file doesn't exists",err)
		return
	}

	_, err3 := io.Copy(&audio, file)
	if err3 != nil {
		fmt.Println("file error",err3)
		return
	}

	audioStr := audio.String()


	// var bufferRead bytes.Buffer
	// _ = io.TeeReader(&audio, bufferRead)

	w := multipart.NewWriter(&b)

	flask_audio, err := w.CreateFormFile("file", "audio.wav")
	if err != nil {
		return
	}
	_, err = io.Copy(flask_audio, strings.NewReader(audioStr))
	if err != nil {
		return
	}

	w.Close()

	
	req, err := http.NewRequest("POST", "http://34.92.72.160:8080/gettext", &b)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("check here")
		return
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := "unable to read response data kas"
		fmt.Println("checkpoint1")
		log.Println("Msg: " + errMsg + " Error: " + err.Error())
		c.JSON(422, map[string]interface{}{
			"message": "Audio file corrupted",
		}) 
		return
	}

	myJson := string(responseData)
	jsonConfig := []byte(myJson)
	var Out FlaskOut
	err = json.Unmarshal(jsonConfig, &Out)
	if err != nil {
		errMsg := "unable to unmarshal http response mad"
		log.Println("Msg: " + errMsg + " Error: " + err.Error())
		c.JSON(422, map[string]interface{}{
			"message": "Audio file corrupted",
		})
		return
	}

	//Call ItemsFetch from DB func 

	var itemss models.Item
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	if len(Out.Keyword) > 0{
		itemss, err = models.FindItemByKeyword(Out.Keyword[0], &conn)
		if err != nil {
			fmt.Println(err)
		}
	}else{
		c.JSON(422, map[string]interface{}{
			"message": "unable to fetch products",
		})
		return
	}
	


	fmt.Println("DB Output",itemss.ProductLocation)


	// Call to emotion APi

	var b2 bytes.Buffer
	// var audio2 bytes.Buffer

	// _, err = io.Copy(&audio2, file)
	// if err != nil {
	// 	return
	// }


	w2 := multipart.NewWriter(&b2)

	flask_audio2, err2 := w2.CreateFormFile("file", "audio.wav")
	if err2 != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Audio file corrupted",
		})
		return
	}

	// fmt.Println("len audio",len(bufferRead.Bytes()),file)
	_, err = io.Copy(flask_audio2, strings.NewReader(audioStr))
	if err != nil {
		
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Audio file corrupted",
		})
		return
	}

	

	w2.WriteField("text", itemss.ProductLocation)
	w2.Close()

	req, err = http.NewRequest("POST", "http://34.92.72.160:8080/getspeech", &b2)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w2.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println()
		return
	}

	if res.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "error",
		})
	} else {


		fmt.Println("Success",res.Body)
		body,err:= ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}


		var out emotionOut
		_ = json.Unmarshal([]byte(body),&out)

		fmt.Println("enmotion output",out)

		c.JSON(http.StatusOK, map[string]interface{}{
			"message": out.Message,
		})

		
	}





	// if resp.StatusCode != http.StatusOK {
	// 	c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 		"message": "error",
	// 	})
	// } else {
	// 	c.JSON(http.StatusOK, map[string]interface{}{
	// 		"message":  "Data is inserted succesfully",
	// 		"products": Out.Keyword,
	// 	})
	// }
}




func Emotion(c *gin.Context) {
	
	var b bytes.Buffer
	var audio bytes.Buffer
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return
	}
	_, err = io.Copy(&audio, file)
	if err != nil {
		return
	}

	w := multipart.NewWriter(&b)

	flask_audio, err := w.CreateFormFile("file", "audio.wav")
	if err != nil {
		return
	}
	_, err = io.Copy(flask_audio, &audio)
	if err != nil {
		return
	}


	text := c.PostForm("productlocation")
	fmt.Println("params",text)
	w.WriteField("text", text)
	w.Close()

	req, err := http.NewRequest("POST", "http://34.92.72.160:8080/getspeech", &b)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "error",
		})
	} else {


		fmt.Println("Success",res.Body)
		body,err:= ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}


		var out emotionOut
		_ = json.Unmarshal([]byte(body),&out)

		fmt.Println("enmotion output",out)

		c.JSON(http.StatusOK, map[string]interface{}{
			"message": out.Message,
		})

		
	}
}

func FindItem(c *gin.Context) {
	keyword := c.PostForm("productname")
	fmt.Println("keyword breakpoint" ,keyword)
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	items, err := models.FindItemByKeyword(keyword, &conn)
	if err != nil {
		fmt.Println("couldn't find the item you're looking for " ,err)
		fmt.Println(keyword)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// list all the items available
func ItemsIndex(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	items, err := models.GetAllItems(&conn)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}


func ItemsForSaleByCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	items, err := models.GetItemsBeingSoldByUser(userID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// create items in dashboard
func ItemsCreate(c *gin.Context) {
	userID := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	item := models.Item{}
	// c.ShouldBindJSON(&item)
	item.ProductName = c.Request.PostFormValue("productname")
	item.Descript = c.Request.PostFormValue("descript")
	item.ProductLocation = c.Request.PostFormValue("productlocation")
	item.Product_Cat = c.Request.PostFormValue("productcat")
	err := item.Create(&conn, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid form data)": err.Error()})
		return 
	}

	c.JSON(http.StatusOK, item)
}


// update items
func ItemsUpdate(c *gin.Context) {
	userID := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	itemSent := models.Item{}
	err := c.ShouldBindJSON(&itemSent)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form sent"})
		return
	}

	itemBeingUpdated, err := models.FindItemByKeyword(itemSent.ProductName, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if itemBeingUpdated.ProductId.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this item"})
		return
	}

	itemSent.ProductId = itemBeingUpdated.ProductId
	err = itemSent.Update(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": itemSent})
}

