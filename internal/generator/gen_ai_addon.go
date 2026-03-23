package generator

import (
	"archive/zip"
	"fmt"
	"log"
)

// AIAddonGen implements AddonGenerator for the "ai" addon category.
// It writes internal/ai/client.go using the first selected provider.
type AIAddonGen struct{}

func (a *AIAddonGen) Generate(folderName string, addons []string, zw *zip.Writer) error {
	if len(addons) == 0 {
		return nil
	}
	// Use the first selected provider; only one client.go is generated.
	provider := addons[0]
	content := GenerateAIAddonContent(provider)
	if err := addToZip(zw, fmt.Sprintf("%s/internal/ai/client.go", folderName), content); err != nil {
		log.Printf("[ERROR] ai addon: %v", err)
		return fmt.Errorf("ai addon: %w", err)
	}
	return nil
}

// VectorStoreAddonGen implements AddonGenerator for the "vectorstore" addon category.
// It writes internal/vectorstore/store.go using the first selected store.
type VectorStoreAddonGen struct{}

func (v *VectorStoreAddonGen) Generate(folderName string, addons []string, zw *zip.Writer) error {
	if len(addons) == 0 {
		return nil
	}
	// Use the first selected store; only one store.go is generated.
	store := addons[0]
	content := GenerateVectorStoreContent(store)
	if err := addToZip(zw, fmt.Sprintf("%s/internal/vectorstore/store.go", folderName), content); err != nil {
		log.Printf("[ERROR] vectorstore addon: %v", err)
		return fmt.Errorf("vectorstore addon: %w", err)
	}
	return nil
}
