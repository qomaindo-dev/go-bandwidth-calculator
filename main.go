package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

func main() {
	fmt.Println("=== Contoh Kalkulator Perhitungan Bandwidth dengan Membaca Kecepatan Internet Otomatis ===")

	// URL file untuk melakukan test kecepatan jaringan (gunakan file kecil untuk efisiensi)
	speedTestURL := "https://URL-FILE-DISINI.zip"

	// 1. Mengukur kecepatan download
	bandwidthMbps, err := measureDownloadSpeed(speedTestURL)
	if err != nil {
		fmt.Printf("Error mengukur kecepatan internet: %v\n", err)
		return
	}

	fmt.Printf("Kecepatan internet terukur: %.2f Mbps\n", bandwidthMbps)

	// 2. Mengambil ukuran file menggunakan URL yang diinginkan (bisa sama atau beda URL diatas)
	fileURL := speedTestURL // Custome URL berbeda disini jika di inginkan
	fileSize, unit, err := getFileSize(fileURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 3. Menghitung waktu download file
	fileSizeInMb := convertToMegabit(fileSize, unit)
	timeInSeconds := fileSizeInMb / bandwidthMbps
	timeFormatted := formatTime(timeInSeconds)

	// 4. Hasil dari perhitungan tersebut
	fmt.Printf("\n=== Hasil Perhitungan ===\n")
	fmt.Printf("URL File: %s\n", fileURL)
	fmt.Printf("Ukuran file: %.2f %s\n", fileSize, unit)
	fmt.Printf("Kecepatan internet: %.2f Mbps\n", bandwidthMbps)
	fmt.Printf("Total bandwidth digunakan: %.2f Megabit (Mb)\n", fileSizeInMb)
	fmt.Printf("Waktu estimasi unduh: %s\n", timeFormatted)
}

// Mengukur kecepatan download dengan mengunduh file kecil
func measureDownloadSpeed(url string) (float64, error) {
	startTime := time.Now()

	// Download file
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Membuang data (karena kita hanya perlu mengukur kecepatan, bukan menyimpan file)
	n, _ := io.Copy(io.Discard, resp.Body)

	// Menghitung kecepatan dalam Mbps
	duration := time.Since(startTime).Seconds()
	sizeInMb := float64(n) * 8 / (1024 * 1024) // Byte -> Megabit
	speedMbps := sizeInMb / duration

	return speedMbps, nil
}

// Mengambil ukuran file dari header HTTP (link sebelumnya diatas)
func getFileSize(url string) (float64, string, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("error: HTTP status %d", resp.StatusCode)
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, "", fmt.Errorf("error: Content-Length tidak ditemukan")
	}

	sizeBytes, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, "", err
	}

	// Konversi ke satuan MB/GB
	sizeMB := float64(sizeBytes) / (1024 * 1024)
	if sizeMB < 1024 {
		return sizeMB, "MB", nil
	} else {
		return sizeMB / 1024, "GB", nil
	}
}

// Mengkonversikan ukuran file ke Megabit (Mb)
func convertToMegabit(size float64, unit string) float64 {
	switch unit {
	case "GB":
		return size * 1024 * 8
	case "MB":
		return size * 8
	default:
		return size
	}
}

// Ubah format waktu
func formatTime(seconds float64) string {
	hours := math.Floor(seconds / 3600)
	remaining := seconds - hours*3600
	minutes := math.Floor(remaining / 60)
	secs := math.Floor(remaining - minutes*60)

	if hours > 0 {
		return fmt.Sprintf("%.0f jam %.0f menit %.0f detik", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%.0f menit %.0f detik", minutes, secs)
	} else {
		return fmt.Sprintf("%.0f detik", secs)
	}
}
