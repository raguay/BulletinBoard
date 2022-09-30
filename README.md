![BulletinBoard](https://socialify.git.ci/raguay/BulletinBoard/image?description=1&forks=1&issues=1&language=1&name=1&owner=1&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Dark)

[![Richard's GitHub stats](https://github-readme-stats.vercel.app/api?username=raguay)](https://github.com/anuraghazra/github-readme-stats)


# BulletinBoard

![BulletinBoard](https://github.com/raguay/BulletinBoard/blob/main/images/BulletinBoard.png)

BulletinBoard is a program for showing information and dialogs to the user. Scripts can use it to get information from the user. It is used by the [EmailIt](https://GitHub.com/raguay/EmailIt) program to ask information from the user for filling in a template. Custom dialogs are easily created using the cli dialog builder. 

![Building a Dialog](https://github.com/raguay/BulletinBoard/blob/main/images/useryQuery-demo.gif)

This is a [Wails 2](https://wails.io/) version of my original [BulletinBoard-NWJS](https://github.com/raguay/BulletinBoard-NWJS) project. I use this program everyday in my workflow. I'm hopeful that you will find it useful as well. You can discuss about this program in the [discussions board](https://github.com/raguay/BulletinBoard/discussions).

## Table of Contents

- [How to Build](#how-to-build)
- [Using BulletinBoard](#using-bulletinboard)
- [Articles about BulletinBoard](#articles-about-bulletinboard)

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

The mask script file will also package the default dialogs for the user to use. If you use the wails build command to make the binary, you would need to copy the `./dialogs` directory contents to your `~/.config/bulletinboard/dialogs` directory to use them.

The executable file will be created in the `build/bin` directory. 

To run the development environment, you type:

```sh
wails dev
```

Alternatively, you can install and use the builds in the releases page. I've only used this program on a macOS system, but it should build and run okay on other OSes. I hopefully will make some builds for Windows and Linux soon.

## Using BulletinBoard

The first thing to be done is to create a command line alias to the main program in order to use the command line interface. To do that, add the following for the shell config file of your shell:

```sh
alias bb="/Applications/BulletinBoard.app/Contents/macOS/BulletinBoard"
```

Now you can use it in the command line with `bb`. Be careful not to use the command line without any arguments as that will run the gui Application.

## Articles about BulletinBoard

- [Building Bulletin Board](https://blog.customct.com/building-bulletin-board)
- [Adding a Bubbletea CLI Interface](https://blog.customct.com/adding-a-bubbletea-cli-interface)
