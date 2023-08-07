# Go-Practitioner
This Repo will consist of some practitioner level courses for Go.
- [Concurrent Programmung in Go](#concurrent-programming-in-go)

## Concurrent Programming in Go
### Concurrency in Go
- Demo: Goroutines and Waitgroups
    - see Goroutines_and_waitgroups folder
- Demo: Channels
    - see Channels folder
### Goroutines
- What is a Goroutine?
    - ```
        [   Go program  ]
        [   goroutine   ]
        [   Scheduler (Go runtime)  ]
        [   threads     ]
        [   Operating system    ]
        ```
    - Go routines sit between the Go Program and the Scheduler.
    - Threads are inbetween the Scheduler and Operating system.
    - Threads are the construction that operating systems use to manage thier concurrency.
    - The scheduler will map goroutines onto operating system threads. Remeber threads are where the operating system actually executes those concurrent tasks.
    - So if a goroutine wants to accomplish anything it has to be scheduled onto a thread.
    - The scheduler can optimize the amount of time that the goroutines are spending on threads to improve the performance of our program.
    - Threads versus Goroutines
        - Thread
            - Have own execution stack
            - Fixed stack space (around 1 MB)
            - Managed by Operating System
            - Relatively expensive
        - Goroutine
            - Have own execution stack
            - Variable stack space (starts @ 2 KB) (this makes goroutines extremely light weight in terms of memory foot print)
            - Managed by Go runtime (scheduler)
            - Inexpensive
- Lifecycle of a Goroutine
    - Three states of a goroutine lifecycle
        - Create
        - Execute
            - Go back and forth from Blocking and Running
            - goroutines also have the responsibilty of execute deferred functions as well.
            - Why would a goroutine block
                - System calls
                - Sleeping
                - Network
                - Blocking
                - Queuing
                - etc. (too many go routines and not enough threads.)
        - Exit
- Advice Regarding Goroutines
    - Goroutines are cheap - use them!
    - Know how a goroutine will stop when you start it
    - Use channels to communicate between goroutines
    - Use sync.WaitGroup to synchronize completion of tasks
### Channels