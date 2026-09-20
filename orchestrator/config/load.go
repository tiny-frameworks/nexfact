// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

func checkContainer() error {
	heartbeatFrequency := SystemParams.Engine.HeartbeatFrequency * 2
	providerRoot := filepath.Join(SellerParams.SellerRoot, "data")

	if SystemParams.Engine.PdfProvider == ProviderContainer {
		hbFile := filepath.Join(providerRoot, "pdfs", "logs", "nexgate.heartbeat")
		if _, err := writer.CheckHeartbeat(hbFile, heartbeatFrequency); err != nil {
			return err
		}
	}
	if SystemParams.Engine.ZugferdProvider == ProviderContainer {
		hbFile := filepath.Join(providerRoot, "zugferds", "logs", "nexgate.heartbeat")
		if _, err := writer.CheckHeartbeat(hbFile, heartbeatFrequency); err != nil {
			return err
		}
	}
	return nil
}

func checkTimeout(timeout, lower, upper time.Duration) error {
	if timeout < lower {
		return errors.New(
			errors.InvalidValue,
			fmt.Sprintf("timeout too small: %v", timeout),
			"orchestrator.config.load.timeout",
		)
	}
	if timeout > upper {
		return errors.New(
			errors.InvalidValue,
			fmt.Sprintf("timeout too big: %v", timeout),
			"orchestrator.config.load.timeout",
		)
	}
	return nil
}

func writeEnv() error {
	sellerEnv := map[string]time.Duration{
		"CMD_TIMEOUT":         SystemParams.Engine.CmdTimeout,
		"HEARTBEAT_FREQUENCY": SystemParams.Engine.HeartbeatFrequency,
		"WATCH_SLEEP":         SystemParams.Engine.WatchSleep,
		"CONFIG_RELAOAD":      SystemParams.Engine.ConfigReloadFrequency,
	}

	sellerDataPath := filepath.Join(SellerParams.SellerRoot, "data")
	sellerPdfEnvPath := filepath.Join(sellerDataPath, "pdfs", "logs", "nexgate.env")
	sellerZugferdEnvPath := filepath.Join(sellerDataPath, "zugferds", "logs", "nexgate.env")

	pFile, pErr := os.Create(sellerPdfEnvPath)
	if pErr != nil {
		return errors.Wrap(errors.WriteError, fmt.Sprintf("Error creating file: %s", sellerPdfEnvPath), "orchestrator.config.load.writeEnv", pErr)
	}
	defer pFile.Close()

	zFile, zErr := os.Create(sellerZugferdEnvPath)
	if zErr != nil {
		return errors.Wrap(errors.WriteError, fmt.Sprintf("Error creating file: %s", sellerZugferdEnvPath), "orchestrator.config.load.writeEnv", zErr)
	}
	defer zFile.Close()

	// An elegant Go approach: Wrap both files in a MultiWriter. A single write operation writes to both.
	multiWriter := io.MultiWriter(pFile, zFile)

	for key, val := range sellerEnv {
		secondString := strconv.FormatFloat(val.Seconds(), 'f', -1, 64)
		envString := fmt.Sprintf("%s=%s\n", key, secondString)
		if _, err := multiWriter.Write([]byte(envString)); err != nil {
			return errors.Wrap(errors.WriteError, "failed to write env configuration", "orchestrator.config.load.writeEnv", err)
		}
	}

	return nil
}
