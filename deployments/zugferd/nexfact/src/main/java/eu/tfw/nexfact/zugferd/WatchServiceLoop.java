// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package eu.tfw.nexfact.zugferd;


import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.PrintStream;
import java.nio.file.Path;
import java.time.Duration;
import java.time.ZoneOffset;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Properties;
import org.tinylog.Logger;

public class WatchServiceLoop extends JobFileProcessor {

    private static final DateTimeFormatter ISO_FORMATTER = DateTimeFormatter.ofPattern("uuuu-MM-dd'T'HH:mm:ssxxx");
    private ZonedDateTime lastHeartbeat;
    private ZonedDateTime lastConfigload;
    private final String heartbeatFile;
    private int configReloadFrequency;
    private int watchSleep;
    private int heartbeatFrequency;
           
     
    public WatchServiceLoop(String inDir, String artefactsDir) {
        super(inDir, artefactsDir);
        Path logPath = artefactsRoot.getParent().resolve("logs").toAbsolutePath();
        heartbeatFile = logPath.resolve(Main.HEARTBEAT_FILE).toString();
    }
    
    public void start() {
        Logger.info("Starting smart polling in directory: " + inDir);
        ZonedDateTime now = ZonedDateTime.now(ZoneOffset.UTC);
        
        reloadConfig();
        lastConfigload = now;
  
        writeHeartbeat(now);
        lastHeartbeat = now;

        while (true) {
            now = ZonedDateTime.now(ZoneOffset.UTC);
             if (Duration.between(lastHeartbeat, now).getSeconds() > heartbeatFrequency) {
                writeHeartbeat(now);
                lastHeartbeat = now;
            }
            if (Duration.between(lastConfigload, now).getSeconds() > configReloadFrequency) {
                this.reloadConfig();
                lastConfigload = now;
            }

            try {
                File dir = new File(inDir);
                // Nur nach Dateien suchen, die wirklich .json heißen (und nicht .processing)
                File[] jobs = dir.listFiles((d, name) -> name.endsWith(".json"));

                if (jobs != null) {
                    for (File jobFile : jobs) {
                        processJobFileSafe(jobFile);
                    }
                }
                Thread.sleep(watchSleep);
            } catch (InterruptedException ie) {
                // Beendet die Schleife sauber, wenn der Container heruntergefahren wird (SIGTERM)
                Logger.info("Watcher interrupted. Shutting down.");
                Thread.currentThread().interrupt();
                break;
            } catch (Exception ex) {
                // Verhindert, dass der Watcher bei einem globalen Fehler (z.B. Volume-Disconnect) komplett abstürzt
                Logger.error("Unexpected error in watch loop: {}", ex.getMessage());
            }
        }
    }

    private void writeHeartbeat(ZonedDateTime zdt) {
        try (PrintStream out = new PrintStream(heartbeatFile)) {
            out.println(zdt.format(ISO_FORMATTER)); // → "2026-06-14T15:51:14+00:00"
        } catch (FileNotFoundException ex) {
            Logger.error("Error writing heartbeat: {}" + ex.getMessage());
        }
    }
    
    private void reloadConfig() {
        Path logPath = artefactsRoot.getParent().resolve("logs").toAbsolutePath();
        String configFile = logPath.resolve(Main.CONFIG_FILE).toString();

        Properties envVars = new Properties();
        // 1. read file
        try (FileInputStream fis = new FileInputStream(configFile)) {
            envVars.load(fis);
        } catch (Exception ex) {
            Logger.error("Failed to read properties file for {}: {}", configFile, ex.getMessage());
        }

        // 2. transfer with Lambda-Expression to System.setProperty
        envVars.forEach((key, value) -> System.setProperty((String) key, (String) value));

        watchSleep = Integer.parseInt(System.getProperty("WATCH_SLEEP", Main.DEFAULT_WATCH_SLEEP));
        heartbeatFrequency = Integer.parseInt(System.getProperty("HEARTBEAT_FREQUENCY", Main.DEFAULT_HEARTBEAT_FREQUENCY));
        configReloadFrequency = Integer.parseInt(System.getProperty("CONFIG_RELOAD", Main.DEFAULT_CONFIG_RELOAD));

        Logger.info("Configuration loaded from: {}", configFile);
        Logger.info("    Configuration: WatchSleep = {}, Heartbeat = {}, configReload = {}", watchSleep, heartbeatFrequency, configReloadFrequency);
    }
}
