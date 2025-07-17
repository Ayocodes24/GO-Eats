package restaurant

import (
	"context"
	"fmt"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var ImageUpdateLock *sync.Mutex = &sync.Mutex{}

func (restSrv *RestaurantService) AddMenu(ctx context.Context, menu *restaurant.MenuItem) (*restaurant.MenuItem, error) {
	_, err := restSrv.db.Insert(ctx, menu)
	if err != nil {
		return &restaurant.MenuItem{}, err
	}
	return menu, nil
}

func (restSrv *RestaurantService) UpdateMenuPhoto(ctx context.Context, menu *restaurant.MenuItem) {
	if restSrv.env == "TEST" {
		return
	}
	client := &http.Client{}
	downloadClient := &unsplash.DefaultHTTPImageClient{}
	fs := &unsplash.DefaultFileSystem{}
	menuImageURL := unsplash.GetUnSplashImageURL(client, menu.Name)
	imageFileName := fmt.Sprintf("menu_item_%d.jpg", menu.MenuID)
	imageFileLocalPath := fmt.Sprintf("uploads/%s", imageFileName)
	imageFilePath := filepath.Join(os.Getenv("LOCAL_STORAGE_PATH"), imageFileName)
	err := unsplash.DownloadImageToDisk(downloadClient, fs, menuImageURL, imageFilePath)
	if err != nil {
		slog.Info("UnSplash Failed to Download Image", "error", err)
	}

	go func() {
		ImageUpdateLock.Lock()
		defer ImageUpdateLock.Unlock()
		setFilter := database.Filter{"photo": imageFileLocalPath}
		whereFilter := database.Filter{"menu_id": menu.MenuID}
		select {
		case <-ctx.Done():
			slog.Error("UnSplash Worker::", "error", ctx.Err().Error())
			return
		default:
			_, err := restSrv.db.Update(context.Background(), "menu_item", setFilter, whereFilter)
			if err != nil {
				slog.Info("UnSplash DB Image Update", "error", err)
			}
		}
	}()
}
