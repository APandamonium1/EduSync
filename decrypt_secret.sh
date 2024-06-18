#!/bin/sh

# Decrypt the file
mkdir $HOME/secrets
# --batch to prevent interactive command
# --yes to assume "yes" for questions
gpg --quiet --batch --yes --decrypt --passphrase="$FIREBASE_ADMIN_PASSPHRASE" \
# --output $HOME/secrets/edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json.gpg
--output edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json.gpg