FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# 1. Schlanke Installation
RUN apt-get update && apt-get install -y --no-install-recommends \
    libreoffice-writer \
    python3-uno \
    fonts-dejavu \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# 2. Das "Golden Profile" injizieren
# Wir kopieren dein vorkonfiguriertes Profil direkt an den Standard-Ort
COPY resources/lo_user_profile/ /root/.config/libreoffice/4/

# Pfade im Profil von lokalem Pfad auf Container-Pfad umbiegen
RUN sed -i 's|/home/[^/]*|/root|g' /root/.config/libreoffice/4/user/registrymodifications.xcu

COPY resources/nexfact_manager.py /app/nexfact_manager.py
COPY resources/loContainerEngine.odt /app/loContainerEngine.odt

WORKDIR /app

# Standardmäßig watch, kann via -e NEXFACT_MODE=once überschrieben werden
ENV NEXFACT_MODE=watch
ENV WATCH_DIR=/data/in

# Wir starten das Python-Script als Hauptprozess
CMD ["python3", "nexfact_manager.py"]

# OCI Metadata Labels
LABEL org.opencontainers.image.title="nexfact-pdfmanager"
LABEL org.opencontainers.image.description="Headless LibreOffice container worker for nexfact"
LABEL org.opencontainers.image.authors="Georg Hagn <tiny-frameworks>"
LABEL org.opencontainers.image.url="https://codeberg.org/tiny-frameworks/nexfact"
LABEL org.opencontainers.image.source="https://codeberg.org/tiny-frameworks/nexfact"
LABEL org.opencontainers.image.documentation="https://codeberg.org/tiny-frameworks/nexfact/src/branch/main/README.md"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.vendor="tiny-frameworks"


