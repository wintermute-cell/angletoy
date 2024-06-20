package main

import (
	"fmt"
	"os"
	"scripts/collision"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "gorl"}

    // CATEGORIES
	var cmdCollision = &cobra.Command{
		Use:   "collision",
		Short: "Collision related commands",
	}

	var cmdEntities = &cobra.Command{
		Use:   "entities",
		Short: "Entity related commands",
	}


    // COMMANDS

    // collision
	var cmdGenerateFromImage = &cobra.Command{
		Use:   "from-image",
		Short: "Generate a collider string (list of polygons) from an image. See documentation for details.",
		Run: func(cmd *cobra.Command, args []string) {
            collision.GenerateCollidersFromImage(args[0])
		},
	}
	var cmdPreviewFromString = &cobra.Command{
		Use:   "preview-string",
		Short: "Preview something from string",
		Run: func(cmd *cobra.Command, args []string) {
            collision.PreviewPolygonsFromString(args)
		},
	}

    // entities
	var cmdCreateEntity = &cobra.Command{
		Use:   "create-entity",
		Short: "Create an entity",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating entity with args:", args)
		},
	}

	rootCmd.AddCommand(cmdCollision, cmdEntities)
	cmdCollision.AddCommand(cmdGenerateFromImage, cmdPreviewFromString)
	cmdEntities.AddCommand(cmdCreateEntity)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

