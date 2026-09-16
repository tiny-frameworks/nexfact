#!/bin/bash

# Starts containers for pdf and zugferd
# You MUST have following directories : f.e.: see env/sellers/example/data
# -------------------------------------------------------------------

# Start Terminal im pdf-Container
podman run --rm --name nexfact-pdfservice --entrypoint bash -it \
  -v ./data/pdfs/in:/data/in:Z \
  -v ./data/artefacts:/data/artefacts:Z,shared \
  -v ./data/pdfs/logs:/data/logs:Z \
nexfact_pdfmanager

# Starte den pdf-Container mit watch im Hintergrund (Detached)
podman run -d --name nexfact-pdfwatcher \
  -e NEXGATE_MODE=watch \
  -e WATCH_DIR=/data/in \
  -v ./data/pdfs/in:/data/in:Z \
  -v ./data/artefacts:/data/artefacts:Z,shared \
  -v ./data/pdfs/logs:/data/logs:Z \
  nexfact-pdfmanager


# Start Terminal im zugferd-Container
podman run --rm --name nexfact-zfservice --entrypoint bash -it \
  -v ./data/pdfs/in:/data/in:Z \
  -v ./data/artefacts:/data/artefacts:Z,shared \
  -v ./data/pdfs/logs:/data/logs:Z \
nexfact-zugferd

# Starte den zugferd-Container mit watch im Hintergrund (Detached)
podman run -d --name nexfact-zfwatcher \
  -e NEXGATE_MODE=watch \
  -e WATCH_DIR=/data/in \
  -v ./data/zugferds/in:/data/in:Z \
  -v ./data/artefacts:/data/artefacts:Z,shared \
  -v ./data/zugferds/logs:/data/logs:Z \
  nexfact-zugferd

