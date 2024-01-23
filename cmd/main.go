package main

import (
	"context"
	"fmt"
	"os"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/akinbezatoglu/s3ync/pkg/cmd/root"
)

func main() {
	Execute()
	//--------------------------------------
	//	cfg, err := config.NewConfig()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	p, err := cfg.GetProfileNames()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(p)
	//	cfg.RemoveSync("default", "C:/Users/akin/OneDrive/Documents/s3/deneme2")
	//	cfg.AddSync("default", "C:/Users/akin/OneDrive/Documents/s3/deneme2", "s3ync-bucket-name2", "eu-central-1")
	//	ps, err := cfg.GetProfileNames()
	//	if err != nil {
	//		fmt.Println(ps)
	//	}
	//	px, err := cfg.GetSyncListFromProfile("default")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(px)
	//--------------------------------------
}

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
	}
	if err := root.NewCmdRoot(cfg).ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
