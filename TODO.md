# Parsing Commands (4)
- [] Have a parsing tree of sub commands.
- [] Have pipelinable operations on push tables (find, filter, format ...)
    - [] These would not be separate cli commands but will be operations appliable on push tables.
    They can be interpreted in similar cli syntax (can have flags and arguments)
- [] List of all commands:
    - [] `np init`
    - [] `np help`
    - [] `np today`
    - [] `np view`
        - [] `np view -i index`
        - [] `np view date`
        - [] `np view today`
    - [] `np pt` ^ (* - subcommands, ^ - pipelinable) 
        - [] `np pt add -t title -d data -s script -a annotation` *
        - [] `np pt del -h hash` *
        - [] `np pt mod -h hash -t title -d data -s script -a annotation` *
        - [] `np pt -i index` ^
        - [] `np pt date` ^
        - [] `np pt today` ^
    - [] `np push`
        - [] `np push n`
    - [] `np diff` 
        - [] `np diff today -i` ^
        - [] `np diff today date` ^ 
        - [] `np diff date1 date2` ^
        - [] `np diff -i -j` ^
        - [] `np diff -i` ^
        - [] `np diff date` ^
        - [] `np diff -w|--week` ^
        - [] `np diff -m|--month` ^
        - [] `np diff tom` ^
    - [] `np status`
    - [] `np gc`
        - [] `np gc --from date`

# Push Table Design (1)
- Each row in the push table holds:
    - Hash of entity - hash(entity + title + entity type)
    - Hash of title
    - Pusher state
    - Annotations

# File System Design (2)
- Folders with name of hash title (store folder)
    - index file
        - holds title name
    - hash file
        - holds entity content
- Push table folder: (pts folder)
    - folder: date of pt
        - folders -> folders of title folder.
            - files -> hash of entity which, entity pusher state
                - also contain order, defined only once, when added, items added after are added
                in best known location, (after the one added)
            - index -> holds title pusher state
    - can pack more folders (n) to a pack and hold diffs.
- VIEW file, holds:
    - date of view, date of pt used.
- TODAY file, holds:
    - date of today, date of pt used.
```
.np/
    TODAY
        date of current today
        date of pt used
    VIEW
        date of current view
        date of pt used
    ./store
        ./[title-hash-1]
            index
                Name of the title
            [entity1-hash]
                Content of the entity
                type
            [entity2-hash]
                Content of the entity
                type
        ./[title-hash-2]
            index
                Name of the title
            [entity1-hash]
                Content
                type
    ./pt
        ./29072006 (loose)
            index.l
            ./[title-hash-1]
                index
                    title pusher state
                    annotations
                    order of children?
                [entity1-hash]
                    annotations
                    pusher state
                [entity2-hash]
                    annotations
                    pusher state
            ./[title-hash-2]
            ./[title-hash-3]
        ./31072006 (tight)
            index.t
                date of prevoius pt
                removed entitys hash
            ./[title-hash-1]
                index
                    order modification
                    annotation modification
                    pusher modification
                [entity2-hash]
                    pusher modification
                    annotation modification
                [entity4-hash]
                    new entity
                    psuher state
                    annotation modification
        ./02082006 (tight)
            index.t
                date of previous pt
                removed entitys hash
            ./[title-hash-4]
                index
                    order modification
                    annotation modification
                    pusher modification
                [entity-1-hash]
                    pusher state
                    annotation modification
```

# Language Syntax (3)
```
if
eif
!
$tag
^
```
- `^` marks the element as unmodifiable
- `$` marks the element as important
