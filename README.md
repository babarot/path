dirgram
=======

dirgram = directory + anagram

## Usage

```
input-of-directory-path | dirgram args...
```

## Examples

```console
$ echo /path/to/a/b/c/d/e | dirgram 3 4 5 6
a/b/c/d
$ echo /path/to/a/b/c/d/e | dirgram 3..6
a/b/c/d
$ echo /path/to/a/b/c/d/e | dirgram 3..-1
a/b/c/d
```

```console
$ echo /path/to/a/b/c/d/e | dirgram 3 6 5 4
a/d/c/b
$ echo /path/to/a/b/c/d/e | dirgram 3 3 3
a/a/a
```

```console
$ echo /path/to/a/b/c/d/e | dirgram 3 6 5 4
$ dirname /path/to/a/b/c/d/e
/path/to/a/b/c/d
$ echo /path/to/a/b/c/d/e | dirgram 1..-1
/path/to/a/b/c/d
```

```console
$ echo /path/to/a/b/c/d/e | dirgram 3
a
```

```console
$ echo /path/to/a/b/c/d/e
/path/to/a/b/c/d/e
$ echo /path/to/a/b/c/d/e | dirgram 1..
/path/to/a/b/c/d/e
$ echo /path/to/a/b/c/d/e | dirgram 3..
a/b/c/d/e
$ echo /path/to/a/b/c/d/e | dirgram 3..-3
a/b
$ echo /path/to/a/b/c/d/e | dirgram 3..-4
a
$ echo /path/to/a/b/c/d/e | dirgram 3
a
```
