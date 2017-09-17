# Deduper: deduplication tool

**The problem**: one might make clones of important folders to back up on multiple devices.  Over
time, the multiple copies of files need to be deduplicated such as when one needs to refactor
and streamline the existing digital infrastructure be it your personal property or something more
elaborate.


*deduper* searches a list of folders for duplicate files and selectively deletes the duplicates.

```

Usage of ./deduper:
  -delete-prefix string
      delete dupes that start with this prefix
  -for-real
      minimal output; dry-run without this
  -ignore-suffixes string
      ignore all files with these comma separated suffixes
  -report
      print out report only.  This is on if 'delete-prefix' flag is omitted.  If on, nothing is deleted.
```

The following example searches all files in `/mnt/usb-disk` location and deletes any duplicates
with prefix supplied.  See below for more details on the algorithm.
```
deduper -for-real -delete-prefix 'some/where/duplicate/trip to asia/2011' /mnt/usb-disk
```

# Algorithm

1. A given set of folders is searched for all regular files ignoring any symlink and any files with
supplied suffixes.
2. It gathers a group of suspected duplicates.  Two files are suspected to be duplicates if they
have the same size and basename.
3. It goes through the duplicates and confirm that they are indeeded duplicates by comparing the
md5 sums.
4. If a prefix is supplied, it evaluates a set of paths to be deleted from each set of duplicates.
It does not delete anything from a set if the supplied prefix
matches all paths in a set.  Therefore, only if one or more copy is left undeleted
from a set of suspected duplicates, will any deletion be done.
5. It does not delete anything if `-report` flag is on.  This true even if `-for-real` flag is
supplied.

The tests included are supposed to cover all points above.
# Building
```
git clone github.com:kzw/deduper
cd deduper
make
```
