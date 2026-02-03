# Koneko
A modern cross platform implementation of the popular Win95 cursor chasing cat.
* currently only supports wayland / hyprland systems

### Compilation
If you are not editing the files kindly refrain from this and go to the execution steps.\
Make sure that you have go installed on your system using `go version` \
Clone the repository using `git clone https://github.com/SamrudhNelli/koneko.git`\
Inside the koneko directory run `CGO_CFLAGS="-w" go build -v -o koneko main.go` \
Note that the first compilation can take 10-15 minutes as it contains the huge [gotk4](https://github.com/diamondburned/gotk4) library.


### Execution
Download the executable `koneko` or clone the repo. \
Run the executable using `./koneko` in the parent directory. \
You could run it with opengl using `GSK_RENDERER=gl ./koneko` that does not produce any vulkan output and uses the opengl library.