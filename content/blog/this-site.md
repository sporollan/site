---
title: "This Site"
date: "2026-01-02"
tags: ["site", "projects"]
description: "A quick thought about the existance of this site"
---

## Why did I build a site?

After months of building all sorts of projects, I wanted a single 
place to pull everything together. A home for all the code and notes I write. 
It’s mainly for me, a way to see my own progress 
and stay focused on diving deep into things I enjoy.

This is what this site is. A tool for visualizing and organizing my journey.
Over time it should contain whatever templates and real projects I made that 
are worth highlighting, and hopefully a quick description on the technical
challenges.

## Why Go?

Nothing special. Just decided it would be nice to immerse myself in Go for a few days. After learning basic syntax and common rules about how things are done, jumped into solving all the concurrency problems we've seen at uni, where we used to do java.

No more trying to manage locks to avoid deadlocks, and no more worrying about how to protect data shared among threads. Go lets you pipe data between independent processes (goroutines) and structure your whole program around that principle. But I still had to figure out what this had to do with a static site generator.

After some research and taking great inspiration on [Hugo](https://github.com/gohugoio/hugo) I went on to attempt my own version. Its rushed but it suits my needs. Just the minimal internals to handle many templates and markdown content. And the pipelines to build, test, and deploy to my s3 bucket shared only to cloudflare.

It's a start. A simple tool built to learn, and a place to document what comes next.