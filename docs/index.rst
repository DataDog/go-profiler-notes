go-profiler-notes
=================

This project was started by `Felix Geisendörfer <https://felixge.de/>`_ after joining Datadog's `Continuous Profiler <https://docs.datadoghq.com/profiler/>`_ team in early 2021. Initially it was just a loose collection of markdown files with notes on Go profiling. But over time and thanks to a few `contributors <https://github.com/DataDog/go-profiler-notes/graphs/contributors>`_, it has developed into one of the most in-depth resources on Go profiling.

In the future, additional topics such as runtime tracing, metrics, heap debugging as well as userland metrics and distributed tracing will be covered as well.

Support this project by giving it a |:star:| on GitHub |ico1|

.. |ico1| image:: https://img.shields.io/github/stars/DataDog/go-profiler-notes?style=social
   :alt: Github Stars
   :target: https://github.com/DataDog/go-profiler-notes

.. toctree::
   :hidden:

   profiling/index
   mental-model-for-go/index

.. toctree::
   :maxdepth: 1
   :caption: Mental Model for Go

   mental-model-for-go/goroutine-scheduler
   mental-model-for-go/garbage-collector

.. toctree::
   :maxdepth: 1
   :caption: Profilers

   profiling/cpu-profiler
   profiling/memory-profiler

.. toctree::
   :maxdepth: 1
   :caption: Misc

   misc/stack-traces
   misc/pprof
