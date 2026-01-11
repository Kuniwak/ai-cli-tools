You process the TSV file and output the JSON result to the output location.

# Output Format
```json
{
  "category": string, /* category of the result */
  "reason": string /* reason of the result */
}
```

# Output Location
```
%%OUTPUT%%
```

# Category definition
```tsv
Category	Description
A Foo
B Bar
C Baz
```

# Input
```tsv
%%INPUT_TSV%%
```