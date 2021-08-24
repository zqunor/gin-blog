package v1

import (
	"gin-blog/pkg/logging"
	"log"
	"net/http"

	"gin-blog/models"
	"gin-blog/pkg/e"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetArticle(c *gin.Context)  {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	
	code := e.SUCCESS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Fatal("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

func GetArticles(c *gin.Context)  {
	title := c.Query("title")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	valid := validation.Validation{}

	if title != "" {
		maps["title"] = title
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo("state").MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}
	
	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
		
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

func AddArticle(c *gin.Context)  {
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(tagId, "tag_id").Message("标签不能为空")
	valid.Min(tagId, 1, "tag_id").Message("标签必须大于0")
	valid.Required(title, "title").Message("文章标题不能为空")
	valid.Required(desc, "desc").Message("文章简述不能为空")
	valid.Required(content, "content").Message("文章内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Fatal("err.key=%s, err.message=%s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H {
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

func EditArticle(c *gin.Context)  {
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
    if arg := c.Query("state"); arg != "" {
        state = com.StrTo(arg).MustInt()
        valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
    }

	// valid.Required(id, "id").Message("ID不能为空")
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.Required(tagId, "tag_id").Message("标签不能为空")
	valid.Min(tagId, 1, "tag_id").Message("标签必须大于0")
	valid.Required(title, "title").Message("文章标题不能为空")
	valid.MaxSize(title, 100, "title").Message("文章标题最长100个字符")
	valid.Required(desc, "desc").Message("文章简述不能为空")
	valid.MaxSize(desc, 244, "desc").Message("文章简述最长255字符")
	valid.Required(content, "content").Message("文章内容不能为空")
	valid.MaxSize(content, 65535, "content").Message("文章内容最长65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(tagId) {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy

				code = e.SUCCESS  
				models.EditArticle(id, data)
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}

	} else {

		for _, err := range valid.Errors {
			log.Printf("err.key=%s, err.message=%s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

func DeleteArticle(c *gin.Context)  {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}

	valid.Required(id, "id").Message("文章ID不能为空")
	valid.Min(id, 1, "id").Message("文章ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
			models.DeleteArticle(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key= %s,err.message=%s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}

