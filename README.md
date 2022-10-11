# Lambda 

![Lambda gopher picture](https://github.com/koss-null/lambda/blob/master/lambda_favicon.png?raw=true) 

Is a library to provide `map`, `reduce` and `filter` operations on slices in one pipeline.  
The slice can be set by generating function. Also the parallel execution is supported.  
It's strongly expected all function arguments to be **pure functions** (functions with no side effects, can be cached by
args).  

### Basic example:

```go
res := pipe.Slice(a).
	Map(func(x int) int { return x * x }).
	Map(func(x int) int { return x + 1 }).
	Filter(func(x int) bool { return x > 100 }).
	Filter(func(x int) bool { return x < 1000 }).
	Parallel(12).
	Do()
```

To play around with simple examples check out: `pkg/pipe/example/`.  
Also feel free to run it with `go run pkg/pipe/example/`.  
 
Here are some more examples (pretty mutch same as in the example package):  
First of all let's make a simple slice:  

```go
a := make([]int, 100)
for i := range a {
	a[i] = i
}
```

Let's wobble it a little bit:  
```go
pipe.Slice(a).
	Map(func(x int) int { return x * x }).
	Map(func(x int) int { return x + 1 }).
	Filter(func(x int) bool { return x > 100 }).
	Filter(func(x int) bool { return x < 1000 }).
	Do()
```
 
Here are some fun facts: 
* it's executed in *4 threads* by default. (I am still hesitating about it, so the value may still change [most likley to 1]) 
* if there is less than *5k* items in your slice, it will be executed in a single thread anyway. 
* you can set custom number of goroutines to execute on with `Parallel(int)` method. But it is limited from above to *256* (it's also may change if a good reason will be found).  
 
Do I have some more tricks? Shure!  
```go
pipe.Func(func(i int) (float32, bool) {
	rnd := rand.New(rand.NewSource(42))
	rnd.Seed(int64(i))
	return rnd.Float32(), true
}).
	Filter(func(x float32) bool { return x > 0.5 }).
```
This one generates an infinite random sequence greater than 0.5  
Use `Gen(n)` to generate exactly **n** items and filter them after,  
Or just `Take(n)` exactly **n** items evaluating them while you can (until the biggest `int` is sent to a function).  
 
Let's assemble it all in one:  
```go
pipe.Func(func(i int) (float32, bool) {
	rnd := rand.New(rand.NewSource(42))
	rnd.Seed(int64(i))
	return rnd.Float32(), true
}).
	Filter(func(x float32) bool { return x > 0.5 }).
	// Take 100500 values, don't stop generate before get all of them
	Take(100500).
	// There is no Sum() yet but no worries
	Reduce(func(x, y float32) float32 { return x + y})
```
Did you notice the `Reduce`?  
Here is what reduce does:  
Let's say there is **zero** element for type **T** where any `f(zero, T) = T`  
No we can add this **zero** at the beginning of the pipe and apply our function to the first two elements: `p[0]`[**zero**] and `p[1]`  
Result we will put into `p[0]` and move the whole pipe to the left by 1   
Due to what we've said before about **zero** `p[0] = p[1]` now   
Now we are doing the same operation for `p[0]` and `p[1]` while we have `p[1]`   
Eventually `p[0]` will be the result.  
 
Anything else?  
You can  sort faster than the wind using all power of your core2duo:  
```go
pipe.Func(func(i int) (float32, bool) {
	return float32(i) * 0.9, true
}).
	Map(func(x float32) float32 { return x * x }).
	Gen(100500).
	Sort(pipe.Less[float32]).
	// This is how you can sort in parallel (it's rly faster!)
	Parallel(12).
	Do()
```
Also if you don't whant to carry the result array on your shoulders and only worry about the amount, you better use
`Count()` instead of `Do()`.  
 
What's next?  
I hope to provide some roadmap of the project soon.   
Also I am going to craft some unit-tests and may be set up github pipelines eventually.   
Before any roadmap and contribution policies being written I will unlikely accept any merge requests, but I will consider it in some
time.   
Anyway feel free to fork, inspire and use! I will try to supply all version tags by some manual testing and quality
control at least.   
