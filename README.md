# pipe

[![Actions Status](https://github.com/variantdev/pipe/workflows/Go/badge.svg)](https://github.com/variantdev/pipe/actions?query=workflow%3AGo)

`pipe` is a fork of the awesome [go-pipe](https://github.com/go-pipe/pipe) with a few enhancements:

- Support for Go modules
- Support for Go `context.Context` for timeout, cancellation and deadline
- Tested compatibility with macOS(ore more concretely unit tests to ensure it's working on macOS)
- Exec logger for recording executed commands and context propagation for e.g. tracing and metering (See the OpenTelemetry example)

## Usage

```go
// `p` is basically the following shell script rewritten in Go:
//
// #!/usr/bin/env sh
// PIPE_NEW_VAR=new
// echo $PIPE_OLD_VAR $PIPE_NEW_VAR
// PIPE_NEW_VAR=after
// go run a-go-app-prints-PIPE_NEW_VAR
// echo hello | sed s/l/k/g
//
p := pipe.Script(
    pipe.SetEnvVar("PIPE_NEW_VAR", "new"),
    pipe.System("echo $PIPE_OLD_VAR $PIPE_NEW_VAR"),
    pipe.SetEnvVar("PIPE_NEW_VAR", "after"),
    func(s *pipe.State) error {
        count := 0
        prefix := "PIPE_NEW_VAR="
        for _, kv := range s.Env {
            if strings.HasPrefix(kv, prefix) {
                count++
            }
        }
        if count != 1 {
            return fmt.Errorf("found %d environment variables", count)
        }
        return nil
    },
     pipe.Line(
        pipe.Print("hello"),
        pipe.Exec("sed", "s/l/k/g")
    )
)

// output contains everything written to stdout by running the script
output, err := pipe.Output(p)

// combimed contains everything written to stdout and stderr by running the script
combined, err := pipe.CombinedOutput(p)

// errs contains everything written to stderr by running the script
output, errs, err := pipe.DividedOutput(p))

// Cancellation with low-level API
started := time.Now()
s := pipe.NewState(nil, nil)

if err := s(p); err != nil {
  // when the script failed to initialize
}

ch := make(chan error)
go func() {
    ch <- s.RunTasks()
}()
time.Sleep(100 * time.Millisecond)
s.Kill()

err := <-ch

if err.Error() == "explicitly killed" {
  // When the script 
} else {
  // When the script succeeded
}

// context.Context support
started := time.Now()
// ctx, cancel := context.WithDeadline(context.Background(), ...)
// ctx, cancel := context.WithCancel(context.Background())
ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Millisecond)
s := pipe.NewState(nil, nil)

if err := s(p); err != nil {
  // when the script failed to initialize
}

ch := make(chan error)
go func() {
    ch <- s.Run(ctx)
}()

// `cancel()` to cancel

// If not `cancel()`ed, this may be an `explicitly killed` error due to the timeout
err := <-ch

if err.Error() == "context canceled" {
  // The script is cancelled
} else if err.Error() == "context deadline exceeded" 
  // The script is timed out
} else {
  // The script succeeded
}
```

See [gopkg.in/pipe.v2](https://gopkg.in/pipe.v2) for the upstream documentation and usage details.

## Why you'd want to choose this over go-pipe?

`go-pipe` seemed the best option above all the options existed at the time of writing this, but it lacked the following things I wanted:

- Support for Go modules
- Support for Go `context.Context` for timeout, cancellation and deadline
- Tested compatibility with macOS(ore more concretely unit tests to ensure it's working on macOS)
- Exec logger for recording executed commands and context propagation for e.g. tracing and metering (See the OpenTelemetry example)

This project just adds them without (mostly) keeping the same API as `go-pipe`. I'm very eager to submit pull requests to `go-pipe/pipe` if it makes sense but until then I'll use this as a better alternative to `go-pipe` for my specific use-cases.

## Alternatives

For shell scripting in Go, and more concretely "replacing bash scripts with Go apps", I've considered following alternatives.

- https://github.com/ebuchman/go-shell-pipes, 2015, supports pipelines built from only slices-of-strings and only `exec.Cmd`s
- https://github.com/codeskyblue/go-sh, 2013-2019, `exec.Cmd`-like API with command-chaining and sessions
- https://github.com/go-pipe/pipe, 2014, [blog(2013)](https://blog.labix.org/2013/04/15/unix-like-pipelines-for-go), Go-native DSL for shell scripting and pipelining, supports `State.Kill` for cancellation, has ability to add custom DSL functions
- https://github.com/mattn/go-pipeline, 2017-, [blog(2015)](https://mattn.kaoriya.net/software/lang/go/20151030131242.htm), supports pipelines built from only slices-of-strings, best for capturing output/combimed output given commands(slices-of-strings)
- https://github.com/b4b4r07/go-pipe, 2019-, runs multiple `exec.Cmd`s in a pipeline and capture the output as bytes.Buffer 
- https://github.com/urjitbhatia/gopipe, 2016-, stream-filter like API for Go, not exactly for shell scripting
- https://github.com/spatialcurrent/go-pipe, 2019-, RX-like API for Go
- https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html, vanilla Go os/exec, versatile but lots of boilterplate
- https://github.com/progrium/go-basher, GO API for embedded bash
