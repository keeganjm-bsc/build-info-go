package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jfrog/build-info-go/build"
	"github.com/jfrog/build-info-go/utils"
	clitool "github.com/urfave/cli/v2"
)

func GetCommands(logger utils.Log) []*clitool.Command {
	return []*clitool.Command{
		{
			Name:      "go",
			Usage:     "Collect build-info for a Go project",
			UsageText: "bi go",
			Action: func(context *clitool.Context) (err error) {
				service := build.NewBuildInfoService()
				service.SetLogger(logger)
				bld, err := service.GetOrCreateBuild("", "")
				if err != nil {
					return
				}
				defer func() {
					e := bld.Clean()
					if err == nil {
						err = e
					}
				}()
				goModule, err := bld.AddGoModule("")
				if err != nil {
					return
				}
				err = goModule.CalcDependencies()
				if err != nil {
					return
				}
				buildInfo, err := bld.ToBuildInfo()
				if err != nil {
					return
				}
				b, err := json.Marshal(buildInfo)
				if err != nil {
					return
				}
				var content bytes.Buffer
				err = json.Indent(&content, b, "", "  ")
				if err != nil {
					return
				}
				fmt.Println(content.String())
				return
			},
		},
	}
}
