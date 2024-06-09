# RimWorld Tag Fixer

Add the missing version tag to mods which has no issue, so that they won't constantly trigger warnings in third-party mod manager.

## Usage

First of all, subscribe the [No Version Warning][nvw], `tag-fixer` uses this as its database.

```sh
tag-fixer PATH_TO_WORKSHOP -v VERSION_TO_FIX
```

For example:

```sh
tag-fixer "G:\Steam\steamapps\workshop\content\294100" -v 1.5
```

You can also add your own database by using `-f` flag

```
tag-fixer PATH_TO_WORKSHOP -v VERSION_TO_FIX -f FILE1 -f FILE2
```

For example:

```sh
tag-fixer "G:\Steam\steamapps\workshop\content\294100" -v 1.5 -f ModIdsToFix.xml
```

[nvw]: https://steamcommunity.com/sharedfiles/filedetails/?id=2599504692