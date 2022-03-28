package main

import (
  "fmt"
  "os"
  "time"

  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/urfave/cli/v2"
  "github.com/webview/webview"
)

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
  data["err"] = ""
  c.JSON(200, data)
}

type taskHandler struct {
}

func newTaskHandler() *taskHandler {
  return &taskHandler{}
}

func (r *taskHandler) addTask(url string) (string, error) {

}

func (r *taskHandler) queryTask(taskID string) (string, string, error) {

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

    path, status, err := task.queryTask(taskID)
    if err != nil {
      renderErr(c, err)
      return
    }

    renderResult(c, gin.H{"path": path, "status": status})
  })
  r.Run(":12432")
}

func main() {
  app := &cli.App{
    Name: "dl-app",
    Flags: []cli.Flag{
      &cli.BoolFlag{Name: "debug"},
    },
    Action: func(c *cli.Context) error {
      debug := c.Bool("debug")

      go func() {
        httpServer()
      }()

      w := webview.New(debug)
      defer w.Destroy()
      w.SetTitle("DL-Download")
      w.SetSize(800, 600, webview.HintNone)
      w.Navigate("http://192.168.0.105:8000")
      w.Run()

      return nil
    },
  }
  app.Run(os.Args)
}
