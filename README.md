# RimWorld Tag Fixer

Add the missing version tag to mods which has no issue, so that they won't constantly trigger warnings in third-party mod manager.

## Usage

First of all, subscribe the [No Version Warning][nvw], `tag-fixer` uses this as its database.

All you need to do is to specify a version to fix:

```sh
tag-fixer -v VERSION_TO_FIX
# For example
tag-fixer -v 1.6
```

You can also create your own database if you find No Version Warning database cannot satisfy you.

```sh
tag-fixer -v VERSION_TO_FIX -f FILE1 -f FILE2
# For example
tag-fixer -v 1.6 -f lists/1.6.xml
```

You can inspect files in `lists` folders as an example of database.

[nvw]: https://steamcommunity.com/sharedfiles/filedetails/?id=2599504692