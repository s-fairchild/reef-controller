#!/bin/bash

set -o errexit

main() {
    realHome="$(homeCheck)"
    userInput="$1"
    if [[ -n $board ]]; then
        getBoard $board
    fi
    binaryFile="${userInput:-"build/release"}"

    if launchOpenocd "$binaryFile" "$board"; then
        launchMinicom "$binaryFile"
    fi
}

getBoard() {
    case $1 in
    arduino-uno)
        b="at32ap7000.cfg";;
    pico)
        b="rp2040.cfg";;
    esac
    if [[ -n $b ]]; then
        board="$b"
    fi
}

# homeCheck Finds the openocd repo path if this script is ran with sudo
homeCheck() {

    if [[ $UID -eq 0 ]]; then
        echo "/home/$SUDO_USER"
    else
        echo "$HOME"
    fi
}

launchOpenocd() {
    defaultBoard="rp2040.cfg"
    board="${board:-$defaultBoard}"
    openocd -s "${realHome}/src/pico/openocd/tcl/" -f interface/picoprobe.cfg -f target/"$board" -c "program ${1} verify reset exit"
}

launchMinicom() {
    local termDev
    termDev="/dev/ttyACM0"
	if [[ -a $termDev ]]; then
        minicom -D $termDev -b 115200
    else
        echo "No device $termDev found"
    fi
}

main "$1"
