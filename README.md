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
- Buffered and Unbuffered Channels
    - Unbuffered Channels
        - ```
            var ch = make(chan string)
            func sender() {
                ch <- "message"
            }
            func receive() {
                msg := <- ch
                fmt.Println(msg)
            }
            func main() {
                // synchronization elided for clarity
                go sender()
                go receiver()
            }
            ```
        - The senders and the receivers must be synchronized, and it's one of the resposibilities of the scheduler to make sure that that synchronization happens.
    - Buffered Channels
        - ```
            var ch = make(chan string, 1)
            func sender() {
                ch <- "message"
            }
            func receive() {
                msg := <- ch
                fmt.Println(msg)
            }
            func main() {
                // synchtonization elided for clarity
                go sender()
                go receiver()
            }
            ```
        - The second parameter of one tells the make function to construct a channel that has an internal buffer or an internal storage capacity.
        - In this case that capacity is 1, and it allows the channel to store a single message within itself.
        - So now the channel itself can receive that message itself into its buffer. Meaning it doesn't have to be blocked waiting for a receiver.
        - And because the message is sitting in the buffer the recieve function doesn't have to block waiting for a sender because it's already waiting in the buffer.
        - This decouples the sending side of our channel from the receiving side of our channel. Which is one of the primary use cases of buffered channels.
- Demo: Buffered and Unbuffered Channels
    - see Buffered_and_unbuffered_channels folder
    - be mindful of using buffers because you application will have to allocate memory for that buffer.
- Directional Channels
    - ```
        // synchronization code elided for clarity
        func main() {
            ch := make(chan string)         // bidirectional channel
            go func(ch chan string) {
                ch <- "message"
            }(ch)
            go func(ch chan string) {
                fmt.Println(<-ch)
            }(ch)
        }
        ```
    - The reason we pass the channels into the goroutines is that it allows us to do something special and that is creating a directional channel.
    - ```
        // synchronization code elided for clarity
        func main() {
            ch := make(chan string)         // bidirectional channel
            go func(ch chan <- string) {    // send-only channel
                ch <- "message"
            }(ch)
            go func(ch <- chan string) {    // receive-only channel
                fmt.Println(<-ch)
            }(ch)
        }
        ```
- Demo: Directional Channels
    - see Directional_channels folder
- Control Flow with Channels
    - select statements
    - for loops
- Select Statement
    - ```
        func main() {
            ch := make(chan int, 1)         // buffered channels prevents need for
            ch := make(chan string, 1)      // synchronization

            ch1 <- 999
            // ch2 <- "message"

            select {
                case msg := <-ch1:
                    fmt.Println(msg)
                case msp := <-ch2:
                    fmt.Println(msg)
                default:
                    fmt.Println("default")
            }
        }
        ```
    - Select statements allow us to operate between multiple cases based on channel operations that can proceed or not.
    - If both ch1 and ch2 can be sent, then the select statements output will be undefined.
    - The reason for this is because select statements work with channels, and channels imply concurrent behaviors.
    - When multiple cases are valid, one case is selected pseudo-randomly
    - if there are no cases to select we will get a deadlock, our program will crash
    - Blocking select statements wait till at least one case is actionable.
    - You can add a default if you don't want a blocking select.
- Demo: Select Statements
    - see Select_statements folder
- For Loops
    - ```
        func main() {
            ch := make(chan string, 3)

            for _, v := range [...]arr{"foo", "bar", "baz"} {
                ch <- v
            }
            close(ch)       //closes the inlet side of the channel.
            for msg := range ch {
                fmt.Println(msg)        // foo bar baz deadlock
            }
        }
        ```
    - when looping over a channel we're looping over an open collection that's why we get the deadlock at the end.
    - after adding close(ch) we no longer run into the deadlock problem.
- Demo: For Loops
    - see For_loops folder