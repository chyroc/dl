package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chyroc/dl/pkgs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

var (
	downloadDir string
	listenAddr  string = "127.0.0.1:12432"
)

//go:embed dist/umi.css
var umiCss string

//go:embed dist/umi.js
var umiJs string

func init() {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	downloadDir = filepath.Join(h, "Downloads/dl-download")
	os.MkdirAll(downloadDir, os.ModePerm)
}

type request struct {
	URL string `json:"url"`
}

func renderErr(c *gin.Context, err error) {
	if err != nil {
		c.JSON(200, gin.H{"err": err.Error()})
	} else {
		c.JSON(200, gin.H{"err": ""})
	}
}

func renderResult(c *gin.Context, data gin.H) {
	if data["err"] == nil {
		data["err"] = ""
	}
	c.JSON(200, data)
}

type taskHandler struct {
	task map[string]*taskItem
	lock sync.Mutex
}

type taskItem struct {
	URL    string `json:"url"`
	Status string `json:"status"` // running, success, error
	Err    string `json:"err"`
}

func newTaskHandler() *taskHandler {
	return &taskHandler{
		task: map[string]*taskItem{},
		lock: sync.Mutex{},
	}
}

func (r *taskHandler) addTask(url string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	id := uuid.New().String()
	r.task[id] = &taskItem{
		URL:    url,
		Status: "running",
	}
	go func() {
		err := pkgs.DownloadData(&pkgs.Argument{
			Dest: downloadDir,
			URL:  url,
		})

		r.lock.Lock()
		defer r.lock.Unlock()

		if err != nil {
			r.task[id].Status = "error"
			r.task[id].Err = err.Error()
		} else {
			r.task[id].Status = "success"
		}
	}()

	return id, nil
}

func (r *taskHandler) queryTask(taskID string) (*taskItem, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	task := r.task[taskID]
	if task == nil {
		return nil, fmt.Errorf("%q 不是一个合法的任务", taskID)
	}

	return task, nil
}

func httpServer() {
	task := newTaskHandler()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.LoadHTMLFiles("dist/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/umi.css", func(c *gin.Context) {
		c.Data(200, "text/css; charset=utf-8", []byte(umiCss))
	})
	r.GET("/umi.js", func(c *gin.Context) {
		c.Data(200, "application/javascript; charset=utf-8", []byte(umiJs))
	})
	r.POST("/api/save", func(c *gin.Context) {
		req := new(request)
		err := c.BindJSON(req)
		if err != nil {
			renderErr(c, err)
			return
		}
		taskID, err := task.addTask(req.URL)
		if err != nil {
			renderErr(c, err)
			return
		}
		renderResult(c, gin.H{"task_id": taskID})
	})
	r.GET("/api/get_task", func(c *gin.Context) {
		taskID := c.Query("task_id")
		if taskID == "" {
			renderErr(c, fmt.Errorf("任务 ID 为空"))
			return
		}

		taskItem, err := task.queryTask(taskID)
		if err != nil {
			renderErr(c, err)
			return
		}

		renderResult(c, gin.H{"err": taskItem.Err, "status": taskItem.Status})
	})
	r.Run(listenAddr)
}

func main() {
	app := &cli.App{
		Name:  "dl-app",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			go func() {
				time.Sleep(time.Second * 2)
				browser.OpenURL("http://" + listenAddr)
			}()
			httpServer()

			return nil
		},
	}
	app.Run(os.Args)
}
