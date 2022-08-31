# BulletinBoard

BulletinBoard is a program for showing information and dialogs to the user. Scripts can use it to get information from the user. It is used by the [EmailIt](https://GitHub.com/raguay/EmailIt) to as information from the user for filling in a template. It currently only displays messages sent to it. This is a [Wails 2](https://wails.io/) version of my original [BulletinBoard-NWJS](https://github.com/raguay/BulletinBoard-NWJS) project. I use this program everyday in my workflow.

## How to build

You have to have [node.js](https://nodejs.org/en/), [go](https://go.dev/), and [Wails 2](https://wails.io) installed first. To build it, type:

```sh
wails build
```

To build as a macOS universal bundle, you type:

```sh
wails build --platform "darwin/universal"
```

Or, install the [mask](https://github.com/jacobdeichert/mask) script runner and use the `maskfile.md` in this directory by typing:

```sh
mask build
```

The executable file will be created in the `build/bin` directory. I've only used this program on a macOS system.

To run the development environment, you type:

```sh
wails dev
```

## Articles about BulletinBoard

- [Building Bulletin Board](https://blog.customct.com/building-bulletin-board)
