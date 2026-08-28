# name of sever binary
binary_name="server"

# checking for ./$binary_name existing
if [ ! -e ./$binary_name ]; then
  echo "./$binary_name was not found, trying to compile"
  # checking if go is avalible
  if ! command -v go >/dev/null 2>&1 ; then
    echo "Go is not found, make sure that it is installed and in your \$PATH"
    exit 1
  fi
  # downloading dependencies? (я за го не шарю, напиши чё мы тут делаем если я не так понял)
  go mod download
  # compiling server binary as $binary_name
  CGO_ENABLED=0 go build -o $binary_name ./cmd/server/
  # making file executable
  chmod u+x $binary_name
  echo "file \"$binary_name\" was compiled succesfully, running:"
fi

# exporting environment variables
export BASE_SNAKE_LENGTH=5
export MAP_HEIGHT=50
export MAP_WIDTH=50
export TICK_DURATION=1s
export INTEREST_RADIUS=10
export SPAWN_CHUNK_PADDING_FROM_BORDER=1
# running binary
./$binary_name
