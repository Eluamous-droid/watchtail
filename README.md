Will continously monitor a directory and add tail -f to new files being created.

Keeps track of active tails and closes old one before starting new ones, provided the max amount of tails have been reached.
Perfect for the "Cloud native" application youre working with or your logging sidecar application when you just need logs thrown to stdout.

TODO: 

Include filename patterns

Never kill option

More tests

--version
