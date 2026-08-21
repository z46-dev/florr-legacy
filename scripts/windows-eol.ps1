# If you're doing development on Windows and running the workflow locally fails,
# Run this script to convert all files from Windows EOL to Unix EOL. This is necessary
# because the workflow runs in a Linux environment and will fail if files have Windows EOL.

gofmt -w (Get-ChildItem -Recurse -Filter "*.go" -File | ForEach-Object FullName)