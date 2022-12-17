# go-run-and

Can't remember the command to show the size of a directory? Just ask for it!

go-run-and (gra) uses OpenAI Codex to create bash scripts from natural language.

## Installation

```bash
go install github.com/lucabeetz/gra@latest
```

## Usage

```bash
gra "show me the largest file in my downloads folder"
```

## Examples

Count file types in downloads

```bash
> gra "count number of file types in downloads"
Suggested:
# count number of file types in downloads (on macos)

find ~/Downloads -type f | sed -E 's/.*\.([^.]+)$/\1/' | sort | uniq -c | sort -nr
Run? [y/N] | [e] for explanation
> e
Explanation:
find ~/Downloads -type f
find all files in ~/Downloads

sed -E 's/.*\.([^.]+)$/\1/'
extract the file extension

sort
sort the file extensions

uniq -c
count the number of occurrences of each file extension

sort -nr
sort the file extensions by number of occurrences
Run? [y/N]
> y
   2 txt
   2 png
   2 json
   2 ics
   2 DS_Store
   1 mp4
   1 mov
   1 localized
   1 jpg
   1 docx
```