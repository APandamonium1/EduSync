#!/bin/sh

# Decrypt the file
# mkdir $HOME/secrets
# --batch to prevent interactive command
# --yes to assume "yes" for questions
gpg --quiet --batch --yes --decrypt --passphrase="$FIREBASE_ADMIN_PASSPHRASE" \
--output edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json.gpg
# --output $HOME/secrets/edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json.gpg

# Optionally, you can output the decrypted file path for verification
echo "Decrypted file: edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json"