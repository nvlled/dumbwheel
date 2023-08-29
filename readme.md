# dumbwheel

A linux program to rebind mouse-thumb buttons
into mouse scroll wheel.

![sewer mouse](sewer-mouse.jpg)

## What for

I have at least three perfectly good mouse but
have a broken mouse wheel.

## Prerequisites

- root permission
- mouse must have thumb buttons

## Install

1. `$ go install github.com/nvlled/dumbwheel@latest`
2. Disable default mouse thumb button behaviour by
   adding these to your `~/.Xmodmap` file: (create if not existing)
   ```
   pointer = 1 2 3 4 5 6 7 91 92
   ```
3. Test run the program with root permission, example: `sudo $(which dumbwheel)`
4. Create service to run the program automatically (TODO)

## Usage

- press and hold thumb buttons up/down to scroll up/down
- double click thumb button and hold to scroll faster
- wiggle mouse while double pressed to accelerate more
