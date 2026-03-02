package bmgt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type cardAPIResponse struct {
	Data []struct {
		CardImages []struct {
			ImageURL string `json:"image_url"`
		} `json:"card_images"`
	} `json:"data"`
}

func FetchAndSaveImage(name, dir string) error {
	// カード名から画像URLを取得する
	// カード名にスペースなどが含まれるためURLエンコードが必要
	apiURL := fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?name=%s", url.QueryEscape(name))
	
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("カードが見つかりません: %s (Status: %d)", name, resp.StatusCode)
	}

	var apiRes cardAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return fmt.Errorf("JSONデコード失敗: %w", err)
	}

	if len(apiRes.Data) == 0 || len(apiRes.Data[0].CardImages) == 0 {
		return fmt.Errorf("画像URLが見つかりませんでした: %s", name)
	}

	imageURL := apiRes.Data[0].CardImages[0].ImageURL
	imgResp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("画像ダウンロード失敗: %w", err)
	}
	defer imgResp.Body.Close()

	path := filepath.Join(dir, name+".jpg")
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ファイル作成失敗: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, imgResp.Body)
	return err
}