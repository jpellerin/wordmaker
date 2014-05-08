# mkwords

A toy for conlangers or ... other weirdos. Use it to generate words from config files like those you can download from awkwords (http://bprhad.wz.cz/awkwords/).

To install, you'll need a working go build system. Then just:

```
go install github.com/jpellerin/wordmaker/mkwords
```

... and you should have a mkwords binary that you can run to ... make ... words.

The awkwords file format is fairly simple. Here is an example:

```
V:a/i/u/ei/ao/ia/ai
C:p/t/k/s/m/n/b/w/x/y/ts/l/sh/ch
T:p/t/k
F:s
N:m/n/rn/rl/nd/ng
r:CV(N/-CV/-CVN/)
```

Lines started with capital letters define word-part choices. Lines starting with 'r' define word construction patterns. That's about it.

Try `mkwords -h` to get some possibly helpful usage information.
