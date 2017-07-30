# Deduper: deduplication tool

*deduper* searches a list of folders for duplicate files and selectively deletes the duplicates.

```
Usage of ./deduper:
  -delete-prefix string
      delete dupes with paths that start with this prefix
  -quiet
      minimal output
  -report
      print out report only.  This is on if 'delete-prefix' flag is omitted.  If on, nothing is deleted.
```

The following example searches all files in `/mnt/usb-disk` location and deletes any duplicates
with prefix supplied.  See below for more details on the algorithm.
```
deduper -delete-prefix 'some/where/duplicate/trip to asia/2011' /mnt/usb-disk
```

The problem: one might make clones of important folders to back up on multiple devices.  Over
time, the multiple copies of files need to be deduplicated such as when one needs to refactor
and streamline the existing digital infrastructure be it your personal property or something more
elaborate.

# Algorithm

1. A given set of folders is searched for all regular files ignoring any symlink. 
2. It gathers a group of suspected duplicates.  Two files are suspected to be duplicates if they
have the same size and basename.
3. It goes through the duplicates and confirm that they are indeeded duplicates by comparing the
md5 sums.
4. If a prefix is supplied, it determines the set of paths
to be deleted from each set of duplicates.  It does not delete any if the supplied prefix
matches all paths of a set of duplicates.
5. It does not delete anything if `-report` flag is on.

The tests included are supposed to cover all points above.
# Building
```
git clone github.com:kzw/deduper
cd deduper
make
```
