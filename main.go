package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"
	"github.com/tashima42/restaurant-manager/controllers"
	"github.com/tashima42/restaurant-manager/database"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	version = "dev"
	date    = "unknown"
)

type Context struct {
	Port      int
	DB        *sqlx.DB
	JWTSecret []byte
	Logger    *zap.SugaredLogger
	Validate  *validator.Validate
}

func main() {
	app := &cli.App{
		Name:                   "restaurant-manager",
		UseShortOptionHandling: true,
		Version:                version,
		Suggest:                true,
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Serve the api",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "port",
						Usage:    "port to run the service",
						Aliases:  []string{"p"},
						EnvVars:  []string{"PORT"},
						Value:    8080,
						Required: false,
					},
					&cli.StringFlag{
						Name:     "path",
						Aliases:  []string{"d"},
						EnvVars:  []string{"DB_PATH"},
						Value:    "restaurant.db",
						Usage:    "database path",
						Required: false,
					},
				},
				Action: run,
			},
			{
				Name:  "database",
				Usage: "database management tasks",
				Subcommands: []*cli.Command{
					{
						Name:  "migrate",
						Usage: "perform database migrations",
						Subcommands: []*cli.Command{
							{
								Name:  "down",
								Usage: "migrate down",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "path",
										Aliases:  []string{"d"},
										EnvVars:  []string{"DB_PATH"},
										Value:    "restaurant.db",
										Usage:    "database path",
										Required: false,
									},
								},
								Action: func(ctx *cli.Context) error {
									path := ctx.String("path")
									return database.MigrateDown(ctx.Context, path)
								},
							},
						},
					},
				},
			},
		},
	}
	compiled, err := time.Parse("2006-01-02T15:04:05Z", date)
	if err == nil {
		app.Compiled = compiled
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	db, err := database.Open(c.String("path"), c.Bool("migrate-down"))
	if err != nil {
		return err
	}
	defer database.Close(db)

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	ec := &Context{
		Port:      c.Int("port"),
		DB:        db,
		JWTSecret: []byte(c.String("jwt-secret")),
		Logger:    logger.Sugar(),
		Validate:  validator.New(validator.WithRequiredStructEnabled()),
	}
	return runServer(ec)
}

func runServer(ec *Context) error {
	cr := controllers.Controller{
		DB:        ec.DB,
		JWTSecret: ec.JWTSecret,
		Logger:    ec.Logger,
		Validate:  ec.Validate,
	}
	app := fiber.New(fiber.Config{ErrorHandler: cr.ErrorHandler})
	app.Use(requestid.New())
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})
	app.Post("/user", cr.CreateUser)
	app.Post("/signin", cr.SignIn)
	app.Use(cr.ValidateToken)
	app.Get("/hello", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*database.User)
		return c.SendString("Hello, " + user.Name)
	})
	app.Post("/table", cr.CreateTable)
	app.Post("/item", cr.CreateItem)
	app.Post("/items", cr.CreateItems)
	app.Get("/items", cr.GetItems)

	return app.Listen(":" + strconv.Itoa(ec.Port))
}
