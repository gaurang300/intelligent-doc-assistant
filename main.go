package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"intelligent-doc-assistant/internal/parser"
	"intelligent-doc-assistant/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path-to-codebase>", os.Args[0])
	}
	codebasePath := os.Args[1]

	store := storage.NewStore()
	if store == nil {
		log.Fatal("Failed to create store")
	}

	ctx := context.Background()

	processFile := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// Parse file for functions & comments
			chunks, err := parser.ParseGoFile(path)
			if err != nil {
				log.Printf("Failed to parse %s: %v", path, err)
				return nil
			}

			// Store all chunks from the file
			if err := store.StoreChunks(ctx, chunks); err != nil {
				log.Printf("Failed to store chunks for %s: %v", path, err)
			}
		}
		return nil
	}

	if err := filepath.Walk(codebasePath, processFile); err != nil {
		log.Fatal("Error walking codebase:", err)
	}

	fmt.Println("✅ Codebase ingestion completed successfully")
}
