# Push ke new tag
git add .
git commit -m "feat: fitur baru"
git tag -a v2.0 -m "Rilis versi 2.0"
git push origin v2.0

# Merge tag ke branch
git checkout main
git merge v2.1
git push origin main