# test comment
# TODO: test

def zip_object:
    contains(["baz", "bar"])
    @uri "https://www.google.com/search?q=\(.search)"
    to_entries | sort_by(.key) |
    map(.key) as $keys |
    map(.value) | transpose |
    map(
        [$keys, .] | transpose |
        map({key: .[0], value: .[1]}) | from_entries
    )
;

