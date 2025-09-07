package app

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/rafaellevissa/rox-partner/internal/db"
	"github.com/rafaellevissa/rox-partner/internal/parser"
	"github.com/rafaellevissa/rox-partner/internal/unzip"
	"github.com/rafaellevissa/rox-partner/pkg/logger"
)

func IngestTrades(dsn string, path string, tmpDir string, batchSize int) error {
	l := logger.New()
	l.Infof("starting ingest for %s", path)

	dbConn, err := db.Connect(dsn)
	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}
	defer dbConn.Close()

	if err := db.EnsureSchema(dbConn); err != nil {
		return fmt.Errorf("ensure schema: %w", err)
	}

	var zips []string
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat path: %w", err)
	}

	if info.IsDir() {
		err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && filepath.Ext(p) == ".zip" {
				zips = append(zips, p)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("walk dir: %w", err)
		}
	} else if filepath.Ext(path) == ".zip" {
		zips = append(zips, path)
	} else {
		return fmt.Errorf("file %s is not a .zip", path)
	}

	l.Infof("found %d zip(s) to process", len(zips))

	numWorkers := runtime.NumCPU()
	runtime.GOMAXPROCS(numWorkers)
	l.Infof("using %d unzip workers and %d CSV workers", numWorkers, numWorkers)

	csvCh := make(chan string, numWorkers*2)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range csvCh {
				ext := filepath.Ext(f)
				if ext != ".csv" && ext != ".txt" {
					continue
				}
				processCSV(f, dbConn, batchSize, l)
			}
		}()
	}

	var unzipWg sync.WaitGroup
	for _, z := range zips {
		unzipWg.Add(1)
		go func(zipPath string) {
			defer unzipWg.Done()
			l.Infof("extracting zip: %s", zipPath)
			files, err := unzip.Extract(zipPath, tmpDir)
			if err != nil {
				l.Errorf("extract %s: %v", zipPath, err)
				return
			}
			for _, f := range files {
				csvCh <- f
			}
		}(z)
	}

	unzipWg.Wait()
	close(csvCh)
	wg.Wait()

	l.Info("ingest finished")
	return nil
}

func processCSV(f string, dbConn *sql.DB, batchSize int, l *logger.Logger) {
	l.Infof("parsing %s", f)
	trades, err := parser.ParseCSV(f)
	if err != nil {
		l.Errorf("parse %s: %v", f, err)
		return
	}

	for i := 0; i < len(trades); i += batchSize {
		end := i + batchSize
		if end > len(trades) {
			end = len(trades)
		}
		batch := trades[i:end]

		if err := db.InsertTradesBulk(dbConn, batch); err != nil {
			l.Errorf("insert batch: %v", err)
		} else {
			l.Infof("inserted batch %d..%d", i, end)
		}
	}

	_ = os.Remove(f)
}
