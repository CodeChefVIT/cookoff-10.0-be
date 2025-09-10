package main

import (
    "fmt"
    logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
    "github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
    "github.com/CodeChefVIT/cookoff-10.0-be/pkg/router"
    "github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
    "github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    logger.InitLogger()
    utils.LoadConfig()
    utils.InitCache()
    utils.InitDB()
    validator.InitValidator()

    
    controllers.Queries = utils.Queries

    e := echo.New()
    e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
        LogURI:        true,
        LogStatus:     true,
        LogError:      true,
        LogValuesFunc: logger.RouteLogger,
    }))

    router.RegisterRoute(e)

    for _, r := range e.Routes() {
        fmt.Println(r.Method, r.Path)
    }

    e.Start(":" + utils.Config.Port)
}
