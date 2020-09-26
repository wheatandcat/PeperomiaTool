#!/bin/sh

gcloud app versions list --service=peperomia-backend | sort -rk 4,4 | awk '{print $2}' | tail -n +5 | xargs -I {} gcloud app versions delete {}
gcloud app versions list --service=peperomia-web | sort -rk 4,4 | awk '{print $2}' | tail -n +5 | xargs -I {} gcloud app versions delete {}