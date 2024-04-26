
```
# Heuristics
(Mixing Go and Haskell notations...)

## Distances:
0.
	a -> b <~~> a -> b // identical (ID)
1.
	(x1,x2) -> (y1,y2) <~~> (x2,x1) -> (y1,y2) // permute params (PP)
	(x1,x2) -> (y1,y2) <~~> (x1,x2) -> (y2,y1) // permute results (PR)
2. 
	(x1,x2) -> (y1,y2) ~~> x2 -> (y1,y2) // weaken params (WP)
	(x1,x2) -> (y1,y2) ~~> (x1,x2) -> y2 // weaken results (WR)
```

```
- subtyping

- currying? convenient or trouble?

- a -> b ~~ a -> c // possible first step
- a -> b ~~ c -> b // possible last step
```

```
- Fuzzy context: “Type name, type only name” ~~> Guess what user means (module, package, version)
    - type (v.)
    - according to popularity?

- Beyond expectation:
    - `func([]string) []int` ~~> `func[a,b any]([]a) []b`
    - `func([]string) []int` ~~> `func[a,b any]([]a, func(a) b) []b` + `func(string) int`

...
```