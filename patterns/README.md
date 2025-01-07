# Concurrency Patterns
https://medium.com/@gopinathr143/go-concurrency-patterns-a-deep-dive-a2750f98a102

## Worker Pool Pattern
The worker pool pattern involves creating a group of worker goroutines to process tasks concurrently, limiting the number of simultaneous operations. This pattern is valuable when you have a large number of tasks to execute.

Example use cases:
* Handling incoming HTTP requests in a web server.
* Processing images concurrently.

### Notes
I have lots of failures playing with this.  Things I need to figure out
1. How to pick buffer size in "real world"
    - I can easily deadlock using this pattern if buffer doesn't fit all jobs
    - I assume in use cases I would not pre-know number of jobs
2. Upping number of jobs and decreasing number workers easiest way to "break"

### Next Steps
I should try and find a worker pool in an http server or the like, and see how they pick buffer sizes.

## Pipeline Pattern
The pipeline pattern structures a series of processing stages, with each stage executed concurrently. Data flows through these stages sequentially, allowing efficient data transformation and processing.

Example use cases:
* Data processing in ETL (Extract, Transform, Load) pipelines.
* Image processing pipelines in multimedia applications.

### Notes
This felt very natural for stepwise handling.

### Next Steps
I should explore how to handle errors.  In a real pipeline I think I would need to be able to reprocess on an error or the like.

## Fan-out/ Fan-in (Scatter gather)
The fan-out/fan-in pattern involves distributing tasks to multiple worker goroutines (fan-out) and then aggregating their results (fan-in). It’s useful for parallelizing tasks and combining their outcomes.

Example use cases:
* Web scraping multiple websites concurrently and merging the results.
* Aggregating data from multiple sensors in IoT applications.

### NOTES
I think I will need to look at other examples of this.  The example from the article seems very similar to worker pool... BUT this one works how I expect, whereas I ran into all kinds of trouble with deadlocks in worker pool.  Hmmm... I definitely need to better understand the difference.


## Others
### Mutex Pattern
 Protecting shared resources using mutexes (sync.Mutex) to ensure exclusive access.
### Semaphore Pattern
 Controlling access to resources by limiting the number of goroutines allowed at a time.
### Barrier Pattern
 Synchronizing multiple goroutines at specific points in their execution.
### WaitGroup Pattern
 Waiting for a collection of goroutines to finish before proceeding.