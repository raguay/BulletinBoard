## build

```sh
rm -Rf build/bin/BulletinBoard.app/Contents/Resources/dialogs
wails build --platform "darwin/universal"
mkdir build/bin/BulletinBoard.app/Contents/Resources/dialogs
cp dialogs/* build/bin/BulletinBoard.app/Contents/Resources/dialogs
```
