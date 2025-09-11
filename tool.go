package main

import (
	"context"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type Game struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type InputParams struct {
	Name string `json:"name" jsonschema:"description=the name of game"`
}

func GetGame(ctx context.Context, params *InputParams) (string, error) {
	GameSet := []Game{
		{Name: "原神", Url: "https://ys.mihoyo.com/tool"},
		{Name: "鸣潮", Url: "https://mc.kurogames.com/tool"},
		{Name: "明日方舟", Url: "https://ak.hypergryph.com/tool"},
	}
	for _, game := range GameSet {
		if game.Name == params.Name {
			return game.Url, nil
		}
	}
	return "", nil
}

func CreateTool() tool.InvokableTool {
	getGameTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_game",
		Desc: "get a game url by name",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": &schema.ParameterInfo{
					Type:     schema.String,
					Desc:     "game's name",
					Required: true,
				},
			},
		),
	}, GetGame)
	return getGameTool
}
