package cmd

import (
	"github.com/Xillon/golang-todo-api/http"
	"github.com/Xillon/golang-todo-api/repository"
	"go.uber.org/fx"
)

var FxModules = fx.Options(
	fx.Provide(
		repository.ProvideDatabase,
		http.ProvideTodoHandler,
	),
)
