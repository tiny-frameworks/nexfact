# Copyright 2026 Georg Hagn
# SPDX-License-Identifier: Apache-2.0

import os
import time
import subprocess
import logging
import json
import shutil
from datetime import datetime, timezone

# --- FESTE CONFIG & PATHS ---
LO_ENGINE_PATH = "/app/loContainerEngine.odt"
WATCH_DIR = os.getenv("WATCH_DIR", "/data/in")
ARTEFACTS_DIR = "/data/artefacts"
MODE = os.getenv("NEXFACT_MODE", "watch")
LOG_FILE = "/data/logs/nexfact.log"
TEMPLATE_PROFILE = "/root/.config/libreoffice/4"
HEARTBEAT_FILE = "/data/logs/nexfact.heartbeat"
ENV_FILE = "/data/logs/nexfact.env"

STATUS_ERROR = "error"
STATUS_SUCCESS = "success"

# --- DYNAMISCHE CONFIG (Zentrales Dictionary statt globaler Variablen) ---
cfg = {
    "HEARTBEAT_FREQUENCY": 30,
    "WATCH_SLEEP": 2,
    "CMD_TIMEOUT": 10,
    "CONFIG_RELOAD": 240
}

# Logging Setup
os.makedirs(os.path.dirname(LOG_FILE), exist_ok=True)
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[logging.FileHandler(LOG_FILE), logging.StreamHandler()]
)

def reload_config_from_env():
    """Liest die .env Datei nur, wenn sie existiert, und aktualisiert das cfg-Dict."""
    if os.path.exists(ENV_FILE):
        with open(ENV_FILE, 'r') as f:
            for line in f:
                # Kommentare und leere Zeilen ignorieren
                line = line.strip()
                if not line or line.startswith('#'):
                    continue
                
                # Zeile bei '=' splitten
                if '=' in line:
                    key, value = line.split('=', 1)
                    # Umgebungsvariable setzen (für den Fall, dass ich sie woanders brauche)
                    os.environ[key.strip()] = value.strip()
        
        # Jetzt das Konfig-Dict aktualisieren
        cfg["HEARTBEAT_FREQUENCY"] = int(os.getenv("HEARTBEAT_FREQUENCY", "30"))
        cfg["WATCH_SLEEP"] = int(os.getenv("WATCH_SLEEP", "2"))
        cfg["CMD_TIMEOUT"] = int(os.getenv("CMD_TIMEOUT", "10"))
        cfg["CONFIG_RELOAD"] = int(os.getenv("CONFIG_RELOAD", "240"))
        logging.info(f"Config geladen: Timeout={cfg['CMD_TIMEOUT']}s, Sleep={cfg['WATCH_SLEEP']}s, Heartbeat={cfg['HEARTBEAT_FREQUENCY']}s, Config={cfg['CONFIG_RELOAD']}s")
        
def write_heartbeat():
    with open(HEARTBEAT_FILE, "w") as f:
        f.write(datetime.now(timezone.utc).isoformat(timespec='seconds') + "\n")

# -------------------------
# Makro Ausführung
# -------------------------
def run_lo_macro_with_timeout(target_file_par, job_id):
    temp_profile = f"/tmp/lo_profile_{job_id}"
    macro_call = f'macro://./Standard.FileEngine.GeneratePdf("{target_file_par}")'
    logging.info(f"par-File: {target_file_par}")
    try:
        shutil.copytree(TEMPLATE_PROFILE, temp_profile, dirs_exist_ok=True)
    except Exception as e:
        logging.error(f"Fehler beim Kopieren des Template-Profils: {e}")
        return False

    cmd = [
        "soffice", "--headless", "--nologo", "--nodefault",
        "--norestore", "--nolockcheck",
        f"-env:UserInstallation=file://{temp_profile}",
        LO_ENGINE_PATH, macro_call
    ]
    # logging.info(f"CMD: {cmd}")
    try:
        # Hier nutzen wir den aktuellen Timeout aus dem cfg-Dict
        result = subprocess.run(cmd, timeout=cfg["CMD_TIMEOUT"], capture_output=True)
        if result.returncode != 0:
            logging.error(f"LO Error Output: {result.stderr.decode('utf-8')}")
        success = (result.returncode == 0)
        
    except subprocess.TimeoutExpired:
        logging.error(f"TIMEOUT! LibreOffice bei Job {job_id} hart beendet.")
        success = False
        os.system("pkill -9 soffice.bin") 
        
    finally:
        if os.path.exists(temp_profile):
            shutil.rmtree(temp_profile, ignore_errors=True)

    return success

# -------------------------
# Fehler-Logik für .err Datei
# -------------------------
def check_for_err_file(base_path):
    err_file = base_path + ".error"
    logging.info(f"err_file: {err_file}")
    if os.path.exists(err_file):
        try:
            with open(err_file, "r") as f:
                content = f.read().strip()
                logging.info(f"Content: {content}")
            os.remove(err_file)
            return content if content else "Unbekannter Fehler im Makro"
        except Exception as e:
            return f"Konnte .error Datei nicht lesen: {e}"
    return None

# -------------------------
# Stateless Job Processing
# -------------------------
def process_single_file(json_filename):
    full_json_path = os.path.join(WATCH_DIR, json_filename)
    logging.info(f"Starte Json-Job: {json_filename}")

    with open(full_json_path, 'r', encoding='utf-8') as f:
        job_data = json.load(f)

    # Die ID holen
    job_id = job_data.get("id")
    
    # Den absoluten Pfad zur .par Datei für LibreOffice bauen
    target_file_rel = f"{job_id}/{job_id}.par"
    target_file_abs = os.path.join(ARTEFACTS_DIR, target_file_rel)
    
    # Den Pfad zur .error Datei direkt im Log-Ordner generieren
    log_dir = os.path.dirname(LOG_FILE) # Ergibt "/data/logs"
    err_file_base = os.path.join(log_dir, job_id) # Ergibt "/data/logs/12345"
    
    logging.info(f"LO Aufruf mit: {target_file_abs}")
    logging.info(f"Erwarte potentielle Fehlerdatei unter: {err_file_base}.error")
    
    success = run_lo_macro_with_timeout(target_file_abs, job_id)
    err_message = check_for_err_file(err_file_base)

    if success and not err_message:
        logging.info(f"Erfolg: {json_filename}")
        job_data["status"] = STATUS_SUCCESS
    else:
        logging.error(f"Fehler bei: {json_filename}")
        job_data["status"] = STATUS_ERROR
        job_data["message"] = err_message if err_message else "Timeout oder Absturz bei Makro-Ausführung"

    with open(full_json_path, 'w', encoding='utf-8') as f:
        json.dump(job_data, f, ensure_ascii=False, indent=4)

    done_path = os.path.join(WATCH_DIR, json_filename.replace(".json", ".done"))
    os.rename(full_json_path, done_path)
    logging.info(f"Job beendet. Lieferschein erstellt: {os.path.basename(done_path)}")

# -------------------------
# Main Loop
# -------------------------
def main():
    logging.info("Starte LibreOffice Listener...")
    time.sleep(2)
    
    last_heartbeat = 0.0
    last_configload = 0.0
    
    if MODE == "watch":
        logging.info(f"Stateless Worker überwacht {WATCH_DIR}")

        while True:
            now = time.time()
            # Gekoppelt: Heartbeat schreiben UND Config neu einlesen passiert im selben Rutsch
            if now - last_heartbeat >= cfg["HEARTBEAT_FREQUENCY"]:
                write_heartbeat()
                last_heartbeat = now
 
            if now - last_configload >= cfg["CONFIG_RELOAD"]:
                reload_config_from_env()
                last_configload = now
                
            ready_files = [f for f in os.listdir(WATCH_DIR) if f.endswith(".json")]
            for f in ready_files:
                process_single_file(f)

            time.sleep(cfg["WATCH_SLEEP"])
    else:
        logging.info(f"Modus: RUN ONCE: Suche erste .json Datei in {WATCH_DIR}")
        for f in os.listdir(WATCH_DIR):
            if f.endswith(".json"):
                logging.info(f"RUN ONCE: found {f}")
                process_single_file(f)
                break

if __name__ == "__main__":
    main()
    
