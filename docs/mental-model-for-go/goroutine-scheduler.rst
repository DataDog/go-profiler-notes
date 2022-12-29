Goroutine Scheduler
===================

Let’s talk about the scheduler first using the example below:

.. code:: go

   func main() {
       res, err := http.Get("https://example.org/")
       if err != nil {
           panic(err)
       }
       fmt.Printf("%d\n", res.StatusCode)
   }

Here we have a single goroutine, let’s call it ``G1``, that runs the
``main`` function. The picture below shows a simplified timeline of how
this goroutine might execute on a single CPU. Initially ``G1`` is
running on the CPU to prepare the http request. Then the CPU becomes
idle as the goroutine has to wait for the network. And finally it gets
scheduled onto the CPU again to print out the status code.

From the scheduler’s perspective, the program above executes like shown
below. At first ``G1`` is ``Executing`` on ``CPU 1``. Then the goroutine
is taken off the CPU while ``Waiting`` for the network. Once the
scheduler notices that the network has replied (using non-blocking I/O,
similar to Node.js), it marks the goroutine as ``Runnable``. And as soon
as a CPU core becomes available, the goroutine starts ``Executing``
again. In our case all cores are available, so ``G1`` can go back to
``Executing`` the ``fmt.Printf()`` function on one of the CPUs
immediately without spending any time in the ``Runnable`` state.

Most of the time, Go programs are running multiple goroutines, so you
will have a few goroutines ``Executing`` on some of the CPU cores, a
large number of goroutines ``Waiting`` for various reasons, and ideally
no goroutines ``Runnable`` unless your program exhibits very high CPU
load. An example of this can be seen below.

Of course the model above glosses over many details. In reality it’s
turtles all the way down, and the Go scheduler works on top of threads
managed by the operating system, and even CPUs themselves are capable of
hyper-threading which can be seen as a form of scheduling. So if you’re
interested, feel free to continue down this rabbit hole via Ardan labs
series on `Scheduling in
Go <https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html>`__
or similar material.

However, the model above should be sufficient to understand the
remainder of this guide. In particular it should become clear that the
time measured by the various Go profilers is essentially the time your
goroutines are spending in the ``Executing`` and ``Waiting`` states as
illustrated by the diagram below.
