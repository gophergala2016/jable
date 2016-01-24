# Jable

A minimalist YouTube audio player for the console.

# Motivation

As a tabbro<sup>tm</sup>, my browser often overloads my RAM and CPU with things I don't really need. YouTube music tabs are one of the main perpetrators. And as a developer, I spend most of my time inside a terminal. You can see where this is going... Listening to YouTube music from your terminal, provided by Jable. 

Bonus points - it will also work for podcasts, talks, just about anything you don't have to actually look at.

# Installation

If you are on 64-bit Linux, [you can download the binary](https://github.com/gophergala2016/jable/releases):

Otherwise, you will have to build Jable yourself, instructions down below.

# Usage

Jable is technically a terminal inside your terminal. Once you start Jable you will be presented with `jable:` cli. Right now Jable has only basic functions:

- `play [QUERY]` searches YouTube for your [QUERY], fetches the first result and plays it. If Jable is already playing something, `play` command will add the song to the queue. You can queue up to 10 songs. 
- `pause` will pause the current playback.
- `resume` will continue the playback where it was previously paused.
- `skip` will skip the current song and start playing the next song in the queue.
- `help` presents you with a nice help dialog, listing all Jable features.
- `exit` gracefully quit Jable, but we all know Ctrl+C is more convinient.

# Build Jable from source

Grab the source:

`git clone https://github.com/gophergala2016/jable`
`cd jable`

First you need to install some dependencies:

`apt-get install libao-dev libmpg123-dev`
`go get bitbucket.org/weberc2/media/ao`
`go get bitbucket.org/weberc2/media/mpg123`

Then just build like any Go program:

`go build`

If something goes wrong, try building with cgo enabled, because imported `mpg123` and `ao` packages are linking before mentioned `libmpg123` and `libao`.

`CGO_ENABLED=1 go build`

If application is working, but there is no audio output, please make sure your libao coniguration points to correct audio driver in `/etc/libao.conf`.

# Gopher Gala

This is a submission for 2016 [Gopher Gala](http://gophergala.com/) by [kortemy](https://github.com/kortemy).


