// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package eu.tfw.nexfact.zugferd;

import org.tinylog.Logger;

public class Main {
    public static String VERSION = "1.0.0";
    public static String MODE;
    public static String WATCH_DIR;
    public static String ARTEFACTS_ROOT;
    public static String DEFAULT_WATCH_DIR;
    public static String DEFAULT_ARTEFACTS_ROOT;
    public static String HEARTBEAT_FILE = "nexfact.heartbeat";
    public static String CONFIG_FILE = "nexfact.env";
    public static String DEFAULT_HEARTBEAT_FREQUENCY = "25"; // seconds
    public static String DEFAULT_CONFIG_RELOAD = "240"; // seconds
    public static String DEFAULT_WATCH_SLEEP = "2"; // seconds
    
    public static void main(String[] args) throws Exception {

        Logger.info("Starte ZUGFeRD System with Version: {}", VERSION);
        // define defaults
        DEFAULT_WATCH_DIR = System.getenv("/data/in");
        DEFAULT_ARTEFACTS_ROOT = "/data/artefacts";
        
        // read env vars
        MODE = System.getenv("NEXFACT_MODE");
        WATCH_DIR = System.getenv("WATCH_DIR");
        ARTEFACTS_ROOT = System.getenv("ARTEFACTS_ROOT");
        
        // set vars (env or default)
        if (ARTEFACTS_ROOT == null || ARTEFACTS_ROOT.isBlank()) {
            ARTEFACTS_ROOT = DEFAULT_ARTEFACTS_ROOT;
        }
        
        if ((WATCH_DIR.isEmpty())) {
            WATCH_DIR = DEFAULT_WATCH_DIR;
        }

        // switch to runOncce od watch 
        if ("watch".equals(MODE)) {
            runWatch();
        } else {
            runOnce();
        }
    }

    private static void runOnce() throws Exception {
        RunOnceService once;
        once = new RunOnceService(WATCH_DIR, ARTEFACTS_ROOT);
        once.start();
     }

    private static void runWatch() throws Exception {
        WatchServiceLoop watchLoop;
        watchLoop = new WatchServiceLoop(WATCH_DIR, ARTEFACTS_ROOT);
        watchLoop.start();
    }
  }
