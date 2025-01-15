path
====

`path` is a CLI command to operate a path such as changing the order of a path with any index.

## Usage

```
input-of-directory-path | path args...
```

## Examples

```console
$ echo /path/to/a/b/c/d/e | path 3 4 5 6
a/b/c/d
$ echo /path/to/a/b/c/d/e | path 3..6
a/b/c/d
$ echo /path/to/a/b/c/d/e | path 3..-1
a/b/c/d
```

```console
$ echo /path/to/a/b/c/d/e | path 3 6 5 4
a/d/c/b
$ echo /path/to/a/b/c/d/e | path 3 3 3
a/a/a
```

```console
$ dirname /path/to/a/b/c/d/e
/path/to/a/b/c/d
$ echo /path/to/a/b/c/d/e | path 1..-1
/path/to/a/b/c/d
```

```console
$ echo /path/to/a/b/c/d/e | path 3
a
```

```console
>>>>>>> 9d24d55 (Rename app, dirgram to path)
$ echo /path/to/a/b/c/d/e
/path/to/a/b/c/d/e
$ echo /path/to/a/b/c/d/e | path 1..
/path/to/a/b/c/d/e
$ echo /path/to/a/b/c/d/e | path 3..
a/b/c/d/e
$ echo /path/to/a/b/c/d/e | path 3..-3
a/b
$ echo /path/to/a/b/c/d/e | path 3..-3 | path 1
a
$ echo /path/to/a/b/c/d/e | path 3..-4
a
$ echo /path/to/a/b/c/d/e | path 3
a
```

```console
$ echo ./local/file/1/2/3 | path 1..3
./local/file/1
```

```console
$ cat paths.txt
/home/babarot/car/toyota/corolla/levin
/home/babarot/car/toyota/
/home/babarot/car/toyota/supra/a80
/home/babarot/car/nissan
/home/babarot/car/nissan/skyline/bnr32
/home/babarot/car/nissan/skyline/er34
/home/babarot/car/honda/city
/home/babarot/car/honda/civic/eg6
/home/babarot/car/honda/integra
/home/babarot/car/subaru
$ cat paths.txt | path 4 | sort | uniq
honda
nissan
subaru
toyota
$ cat paths.txt | path 5 | sort | uniq -c
   1 city
   1 civic
   1 corolla
   1 integra
   2 skyline
   1 supra
```
