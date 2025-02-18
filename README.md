# ❤️CMS

LoveCMS is a CMS for Hugo. I'm writing it because:

1. My partner wants to move off Wordpress
2. I want a little project that I can use to learn Go

# Design Goals
* -Build with an eye toward extensibility. Certain features are going to be dependent on the underlying generator or hosting platform. These features should be written in such a way that they can easily be swapped out to other platforms.
* Hide as much of the technical minutia as possible. At least up front, there's not going to be a way to work around the fact that we'll need a working git repo for a static site, but after that the goal is the user should never have to touch git or the command line.

# Requirements

I'm tracking requirements and progress through the [❤️CMS Roadmap](https://github.com/users/jaredhaight/projects/1/views/1)