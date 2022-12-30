Profiler Comparison
===================

Here is an overview of the profilers built into the Go runtime.

.. list-table::
   :header-rows: 1

   * -
     - :doc:`CPU <cpu-profiler>`
     - :doc:`Memory <memory-profiler>`
     - :doc:`Block <block-profiler>`
     - :doc:`Mutex <mutex-profiler>`
     - :doc:`Goroutine <goroutine-profiler>`
   * - Production Safety
     - |:white_check_mark:|
     - |:white_check_mark:|
     - |:warning:| [#foot-block]_
     - |:white_check_mark:|
     - |:white_check_mark:| [#foot-goroutine]_
   * - Safe Rate
     - default
     - default
     - |:x:| [#foot-block]_
     - ``100``
     - ``1000`` |nbsp| goroutines
   * - Max Stack Depth
     - ``64``
     - ``32``
     - ``32``
     - ``32``
     - ``32`` |nbsp| - |nbsp| ``100`` |nbsp| [#foot-goroutine-api]_
   * - :ref:`profiling/cpu-profiler:profiler labels`
     - |:white_check_mark:|
     - |:x:|
     - |:x:|
     - |:x:|
     - |:white_check_mark:| [#foot-goroutine-api]_

The :doc:`thread-create-profiler` is not listed because it's broken.

.. [#foot-block] The block profiler can cause 5% or more CPU overhead, even when using a high rate value.
.. [#foot-goroutine] Before Go 1.19, this profile caused O(N) stop-the-world pauses where N is the number of goroutines. Expect ~1-10µsec pause per goroutine.
.. [#foot-goroutine-api] Depends on API.

.. |nbsp| unicode:: 0xA0 
   :trim:
