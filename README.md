Will continously monitor a directory and add tail -f to new files being created.
Keeps track of active tails and closes old one before starting new ones, provided the max amount of tails have been reached.

TODO: 
Exclude filename patterns
Include filename patterns
Dont tail a new directory if created
