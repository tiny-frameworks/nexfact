// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package eu.tfw.nexfact.zugferd;


import java.io.File;

public class RunOnceService extends JobFileProcessor {

    public RunOnceService(String inDir, String artefactsDir) {
        super(inDir, artefactsDir);
    }

    public void start() {
        System.out.println("Starting smart searching for first .jsoning in directory: " + inDir);

        while (true) {
            try {
                File dir = new File(inDir);
                // Nur nach Dateien suchen, die wirklich .json heißen (und nicht .processing)
                File[] jobs = dir.listFiles((d, name) -> name.endsWith(".json"));

                if (jobs != null) {
                    // Process only the 1. found .json
                    for (File jobFile : jobs) {
                        processJobFileSafe(jobFile);
                        return;
                    }
                }

                Thread.sleep(500);

            } catch (InterruptedException ie) {
                // Beendet die Schleife sauber, wenn der Container heruntergefahren wird (SIGTERM)
                System.out.println("Watcher interrupted. Shutting down.");
                Thread.currentThread().interrupt();
                break;
            } catch (Exception e) {
                // Verhindert, dass der Watcher bei einem globalen Fehler (z.B. Volume-Disconnect) komplett abstürzt
                System.err.println("Unexpected error in watch loop: " + e.getMessage());
            }
        }
    }
}
