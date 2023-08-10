# Go-Practitioner
This Repo will consist of some practitioner level courses for Go.
- [Concurrent Programmung in Go](#concurrent-programming-in-go)
- [Testing in Go](#testing-in-go)

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
### Common Concurrency Patterns
- Non-blocking Error Channels
    - ```
        var (
            in = make(chan string)
            out = make(chan int)
            errCh = make(chan error, 1)
        )
        func worker(in <-chan string, out chan <- int, err chan<- error) {
            for msg := range in {
                // converts strings into integers with the Atoi function from the string conversion package.
                i, err := strconv.Atoi(msg)
                if err != nil {
                    errCh <- err
                    return
                }
                out <- i
            }
        }
        ```
    - So by adding the buffer to the errCh it allow us to make sure the channel always has a place to put the error at the end of the function.
- Encapsulating Goroutines
    - ```
        var (
            in = make(chan string)
            out = make(chan int)
            errCh = make(chan error)
        )
        func worker(in <-chan string, out chan <- int, err chan<- error) {
            for msg := range in {
                // converts strings into integers with the Atoi function from the string conversion package.
                i, err := strconv.Atoi(msg)
                if err != nil {
                    errCh <- err
                    return
                }
                out <- i
            }
        }
        ```
    - When trying to send a nil value in a channel the application panics
    - ```
        var (
            in = make(chan string)
        )
        func worker(in <-chan string) (chan int, chan error) {
            // intializing the channels inside the function guarantees that nil will not be a value passed into it for out.
            out := make(chan int)
            errCh := make(chan error, 1)
            go func() {
                for msg := range in {
                    i, err := strconv.Atoi(msg)
                    if err != nil {
                        errCh <-
                        return
                    }
                    out <- i
                }
            }()
            return out, errCh
        }
        ```
- Demo: Non-blocking Error Channels and Encapsulating Goroutines
    - see N_e_c_e_g folder
- Messaging Patterns
    - Foundational Concurrency Patterns
        - Single prducer single consumer
            - Producers and consumers are the different sides of the channel
            - single producer means that there is one goroutine sending messages into a channel.
            - single consumer means that there is one goroutine receiving messages from a channel.
        - Single producer multiple consumer
            - multiple goroutines that are pulling messages from the same channel
        - Multiple producer single consumer
            - multiple channels, one goroutine
        - multiple producers multiple consumers
- Demo: Single Producer, Single Consumer
    - see Single_producer_single_consumer folder
- Demo: Single Producer, Multiple Consumer
    - see Single_producer_multiple_consumer folder
- Demo: Multiple Producer, Single Consumer
    - see Multiple_producer_single_consumer folder
    - with multiple producers you need some sort of supervisory function like with the addition goroutine in the reserverInventory function.
    - it keeps an eye on things to make sure it knows when the last message was sent. Because none of the workers knows if it's the last one to send messages.
- Demo: Multiple Producer, Multiple Consumer
    - see Multiple_producer_multiple_consumer folder
### Additional Tools to Support Concurrent Programming
- The Sync Package - Mutexes
    - Mutexes
    - ```
        var m sync.Mutex
        func main() {
            go func() {
                m.Lock()
                defer m.Unlock()
                // do thing
            }()
            go func() {
                m.Lock()
                defer m.Unlock()
                // do other thing
            }()
        }
        ```
    - When the m.Lock() tries to aquire the lock the Go scheduler is going to block all of the goroutines that are trying to acquire that lock at the lock method call.
    - It will then allow one through until unlock is called.
    - As soon as the goroutine calls unlock, then another goroutine is allowed to do what's called acquiring the lock, and then it can proceed.
    - What this allows us to do is this introduces a mechanism to synchronize or make goroutines execute sequentially.
- Demo: Mutexes
    - see Mutexes folder
- Demo : sync.Once
    - see Once folder
    - ```go get github.com/mattn/go-sqlite3``` for database
    - One of the things we do with SQL databases in our programs is we have to create what's called a db object.
    - A db object, that's in the SQL package, and its signature often looks like this. ```var db *sql.DB```
- Demo: The Race Detector
    - see Race_detector folder
- Contexts
    - The Purpose of Contexts
        - The context package was developed as a way to communicate cancelation between goroutines.
        - We can communicate using the context that the operation is no longer necessary, and that context will communicate that canelation over to the goroutines.
        - There are basic forms of the context package.
            - Cancel
            - Timeout
- Demo: Context with Cancel
    - see Context_with_cancel folder
- Demo: Context with Timeout
    - see Context_with_timeout folder
## Testing in Go
### Introduction
- Introduction
    - Why Write Tests?
        - Initial Build: Ensure correctness of application
        - Production: Identify errors before users are impacted
        - Enhancements: Prevent new development from breaking older feature
    - Testing Support in Go
        - Test Runner: some program that is able to execute our tests for us.
        - Testing API: some way to interact with the test and the test runner itself, to let the test runner know what's going on with our test.
        - Assertion/Expectation API: to assert that we expect our test to prove something or state is present in our application.
        - There actually is no assertion or expectation API built into Go.
- Demo: The First Test
    - see The_first_test folder
### Testing Business Logic
- The Testing Pyramid
    - Unit at the bottom: Prove that individual units of logic are correct
    - then Component: Confirm that different application concerns (i.e. packages) perform correctly
    - then Integration: Validate that entire program works as expected
    - then at the top is End to end: Demonstrate that entire system works correctly together
- Testing API
    - The T object is the way that we communicate with the test runner
    - Communicate with test runner
        - Report failures
        - Logging
        - Configure how to test is executed
        - Iteract with enviornment
    - No assertions!
- Reporting Test Failures
    - Non-immediate failures vs Immediate failures
        - Non-immediate
            - t.Fail(): this marks the test as failed and moves on.
            - t.Error(...interface{}): if we want a little bit more information. It's the same as the fail method but it's followed by a call to the log method on that T object.
            - t.Errorf(string,...interface{}): uses a formatting string instead of just sending data to the test log.
        - Immediate
            - t.FailNow()
            - t.Fatal(...interface{})
            - t.Fatalf(string,...interface{})
- Why are Assertions Missing?
    - ```
        // reduce learning curve
        func TestThing(t *testing.T){
            l, r := 2, 4
            expect := 6

            got := l + r

            if got != expect {
                // report failure
            }
        }
        ```
    - Checking for test failures uses same construct as production code.
    - If you can write production code, you can write tests.
- ```
    // focus on common concerns
    func TestThing(t *testing.T){
            l, r := 2, 4
            expect := 6

            got := l + r

            assert.Equal(got, expect)       // assertion style
            Expect(expect).To(Equal(got))   // expect style
    }
    ```
    - Both styles are equally valid and subject to team preference
- Demo: Writing a Unit Test
    - see Writing_a_unit_test folder